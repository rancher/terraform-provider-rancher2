package rancher2

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	ghodssyaml "github.com/ghodss/yaml"
	gover "github.com/hashicorp/go-version"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api/v1"
)

const (
	clusterProjectIDSeparator = ":"
	passDigits                = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	passDefaultLen            = 20
	maxHTTPRedirect           = 5
)

func getMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func AreEqual(o, n interface{}) bool {
	return reflect.DeepEqual(o, n)
}

func Base64Encode(s string) string {
	if len(s) == 0 {
		return ""
	}
	data := []byte(s)

	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}
	data, err := base64.StdEncoding.DecodeString(s)

	return string(data), err
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func getKubeConfigFromObj(kubeconfig *kubeconfig.Config) (string, error) {
	if kubeconfig == nil {
		return "", nil
	}
	config, err := interfaceToMap(kubeconfig)
	if err != nil {
		return "", err
	}

	return mapInterfaceToYAML(config)
}

func getObjFromKubeConfig(config string) (*kubeconfig.Config, error) {
	kubeconfig := &kubeconfig.Config{}
	if len(config) == 0 {
		return kubeconfig, nil
	}
	kubeconfigMap, err := ghodssyamlToMapInterface(config)
	if err != nil {
		return nil, fmt.Errorf("Yaml unmarshall kube_config %v", err)
	}
	kubeconfigJSON, err := mapInterfaceToJSON(kubeconfigMap)
	if err != nil {
		return nil, fmt.Errorf("Json marshall kube_config: %v", err)
	}
	err = jsonToInterface(kubeconfigJSON, kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Json unmarshall kube_config: %v", err)
	}

	return kubeconfig, nil
}

func getTokenFromKubeConfig(config string) (string, error) {
	if len(config) == 0 {
		return "", nil
	}
	kubeconfig, err := getObjFromKubeConfig(config)
	if err != nil {
		return "", err
	}
	if kubeconfig == nil || kubeconfig.AuthInfos == nil || len(kubeconfig.AuthInfos) == 0 {
		return "", nil
	}

	return kubeconfig.AuthInfos[0].AuthInfo.Token, nil
}

func getTokenIDFromKubeConfig(config string) (string, error) {
	token, err := getTokenFromKubeConfig(config)
	if err != nil {
		return "", err
	}
	return splitTokenID(token), nil

}

func updateKubeConfigToken(config, token string) (string, error) {
	if len(token) == 0 {
		return config, nil
	}
	if len(config) == 0 {
		return "", nil
	}
	kubeconfig := &kubeconfig.Config{}
	err := jsonToInterface(config, kubeconfig)
	if err != nil {
		return "", err
	}
	if kubeconfig == nil || kubeconfig.AuthInfos == nil || len(kubeconfig.AuthInfos) == 0 {
		return "", nil
	}

	return kubeconfig.AuthInfos[0].AuthInfo.Token, nil
}

func TrimSpace(val interface{}) string {
	return strings.TrimSpace(val.(string))
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func GetRandomPass(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]byte, n)
	for i := range b {
		b[i] = passDigits[rand.Int63()%int64(len(passDigits))]
	}
	return string(b)
}

func HashPasswordString(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Problem encrypting password: %v", err)
	}
	return string(hash), nil
}

func NewListOpts(filters map[string]interface{}) *types.ListOpts {
	listOpts := clientbase.NewListOpts()
	if filters != nil {
		listOpts.Filters = filters
	}

	return listOpts
}

func DoUserLogin(url, user, pass, ttl, desc, cacert string, insecure bool) (string, string, error) {
	loginURL := url + "/v3-public/localProviders/local?action=login"
	loginData := `{"username": "` + user + `", "password": "` + pass + `", "ttl": ` + ttl + `, "description": "` + desc + `"}`
	loginHead := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	// Login with user and pass
	loginResp, err := DoPost(loginURL, loginData, cacert, insecure, loginHead)
	if err != nil {
		return "", "", err
	}

	if loginResp["type"].(string) != "token" || loginResp["token"] == nil {
		return "", "", fmt.Errorf("Doing  user logging: %s %s", loginResp["type"].(string), loginResp["code"].(string))
	}

	return loginResp["id"].(string), loginResp["token"].(string), nil
}

func DoPost(url, data, cacert string, insecure bool, headers map[string]string) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	if url == "" {
		return response, fmt.Errorf("Doing post: URL is nil")
	}

	jsonBytes := []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return response, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	if cacert != "" {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(cacert)); !ok {
			log.Println("No certs appended, using system certs only")
		}
		transport.TLSClientConfig.RootCAs = rootCAs
	}

	client.Transport = transport

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func DoGet(url, username, password, token, cacert string, insecure bool) ([]byte, error) {
	start := time.Now()

	if url == "" {
		return nil, fmt.Errorf("Doing get: URL is nil")
	}
	log.Println("Getting from ", url)

	client := &http.Client{
		Timeout: time.Duration(60 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxHTTPRedirect {
				return fmt.Errorf("Stopped after %d redirects", maxHTTPRedirect)
			}
			if len(token) > 0 {
				req.Header.Add("Authorization", "Bearer "+token)
			} else if len(username) > 0 && len(password) > 0 {
				s := username + ":" + password
				req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
			}
			return nil
		},
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	if cacert != "" {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(cacert)); !ok {
			log.Println("No certs appended, using system certs only")
		}
		transport.TLSClientConfig.RootCAs = rootCAs
	}
	client.Transport = transport

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Doing get: %v", err)
	}
	if len(token) > 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	} else if len(username) > 0 && len(password) > 0 {
		s := username + ":" + password
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Doing get: %v", err)
	}
	defer resp.Body.Close()

	// Timings recorded as part of internal metrics
	log.Println("Time to get req: ", float64((time.Since(start))/time.Millisecond), " ms")

	return ioutil.ReadAll(resp.Body)

}

func NormalizeURL(input string) (string, error) {
	if input == providerDefaultEmptyString {
		return "", fmt.Errorf("Normalizing url: no api_url provided")
	}
	if input == "" {
		return "", nil
	}
	u, err := url.Parse(input)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return "", fmt.Errorf("Normalizing url %s: %v", input, err)
	}
	// Setting empty url path
	u.Path = ""
	return u.String(), nil
}

func IsUnknownSchemaType(err error) bool {
	return strings.Contains(err.Error(), "Unknown schema type")
}

func IsNotAccessibleByID(err error) bool {
	return strings.Contains(err.Error(), "can not be looked up by ID")
}

func IsNotFound(err error) bool {
	return clientbase.IsNotFound(err)
}

// IsForbidden checks if the given APIError is a Forbidden HTTP statuscode
func IsForbidden(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusForbidden
}

// IsNotAllowed checks if the given APIError is a Method Not Allowed HTTP statuscode
func IsNotAllowed(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusMethodNotAllowed
}

// IsConflict checks if the given APIError is a Conflict HTTP statuscode
func IsConflict(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusConflict
}

// IsServerError checks if the given APIError is a Internal Server Error HTTP statuscode
func IsServerError(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusInternalServerError
}

// IsBadGatewayError checks if the given APIError is a Bad Gateway Server Error HTTP statuscode
func IsBadGatewayError(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusBadGateway
}

// IsServiceUnavailableError checks if the given APIError is a Service Unavailable Server Error HTTP statuscode
func IsServiceUnavailableError(err error) bool {
	apiError, ok := err.(*clientbase.APIError)
	if !ok {
		return false
	}

	return apiError.StatusCode == http.StatusServiceUnavailable
}

func splitTokenID(token string) string {
	separator := ":"

	if strings.Contains(token, separator) {
		return token[0:strings.Index(token, separator)]
	}

	return token
}

func splitBySep(data, sep string) []string {
	if len(sep) == 0 {
		return nil
	}
	return strings.Split(data, sep)
}

func splitID(id string) (clusterID, resourceID string) {
	separator := "."

	if strings.Contains(id, separator) {
		return id[0:strings.Index(id, separator)], id[strings.Index(id, separator)+1:]
	}
	return "", id
}

func splitRegistryID(id string) (namespaceID, projectID, resourceID string) {
	separator := "."

	result := strings.Split(id, separator)

	switch count := len(result); count {
	case 2:
		return "", result[0], result[1]
	case 3:
		return result[0], result[1], result[2]
	}

	return "", "", id
}

func clusterIDFromProjectID(projectID string) (string, error) {
	if projectID == "" || !strings.Contains(projectID, clusterProjectIDSeparator) {
		return "", fmt.Errorf("[ERROR] Getting clusted ID from project ID: Bad project id format %s", projectID)
	}

	return projectID[0:strings.Index(projectID, clusterProjectIDSeparator)], nil
}

func splitProjectIDPart(id string) (projectID string) {
	id = strings.TrimSuffix(id, clusterProjectIDSeparator)

	if strings.Contains(id, clusterProjectIDSeparator) {
		return id[strings.Index(id, clusterProjectIDSeparator)+1:]
	}

	return ""
}

func splitProjectID(id string) (clusterID, projectID string) {
	id = strings.TrimSuffix(id, clusterProjectIDSeparator)

	if strings.Contains(id, clusterProjectIDSeparator) {
		return id[0:strings.Index(id, clusterProjectIDSeparator)], id
	}

	return id, ""
}

func splitAppID(id string) (projectID, appID string, err error) {
	separator := clusterProjectIDSeparator

	fields := strings.Split(id, separator)

	if len(fields) != 3 {
		return "", "", fmt.Errorf("[ERROR] Getting App ID: Bad project id format %s", id)
	}

	return fields[0] + separator + fields[1], fields[1] + separator + fields[2], nil
}

func updateVersionExternalID(externalID, version string) string {
	//Global catalog url: catalog://?catalog=demo&template=test&version=1.23.0
	//Cluster catalog url: catalog://?catalog=c-XXXXX/test&type=clusterCatalog&template=test&version=1.23.0
	//Project catalog url: catalog://?catalog=p-XXXXX/test&type=projectCatalog&template=test&version=1.23.0

	str := strings.TrimPrefix(externalID, AppTemplateExternalIDPrefix)
	values := strings.Split(str, "&")
	out := AppTemplateExternalIDPrefix
	for _, v := range values {
		if strings.Contains(v, "version=") {
			pair := strings.Split(v, "=")
			if pair[0] == "version" {
				pair[1] = version
			}
			v = pair[0] + "=" + pair[1]
		}
		out = out + "&" + v
	}
	return out
}

func toArrayString(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toArrayStringSorted(in []interface{}) []string {
	if in == nil {
		return nil
	}
	out := toArrayString(in)
	sort.Strings(out)
	return out
}

func toArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func toArrayInterfaceSorted(in []string) []interface{} {
	if in == nil {
		return nil
	}
	sort.Strings(in)
	out := toArrayInterface(in)
	return out
}

func toMapString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toMapByte(in map[string]interface{}) map[string][]byte {
	out := make(map[string][]byte)
	for i, v := range in {
		if v == nil {
			out[i] = []byte{}
			continue
		}
		value := v.(string)
		out[i] = []byte(value)
	}
	return out
}

func toMapInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for i, v := range in {
		out[i] = v
	}
	return out
}

func jsonToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := jsonToInterface(in, &out)
	return out, err
}

func jsonToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := json.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func interfaceToMap(in interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{})

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ghodssyamlToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := ghodssyaml.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, err
}

func ghodssyamlToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := ghodssyaml.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func yamlToMapInterface(in string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(in), &out)
	if err != nil {
		return nil, err
	}
	return out, err
}

func yamlToInterface(in string, out interface{}) error {
	if out == nil {
		return nil
	}
	err := yaml.Unmarshal([]byte(in), out)
	if err != nil {
		return err
	}
	return err
}

func mapInterfaceToJSON(in map[string]interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	return interfaceToJSON(in)
}

func mapInterfaceToYAML(in map[string]interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	return interfaceToYAML(in)
}

func interfaceToJSON(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func interfaceToYAML(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := yaml.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func interfaceToGhodssyaml(in interface{}) (string, error) {
	if in == nil {
		return "", nil
	}
	out, err := ghodssyaml.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), err
}

func YAMLToJSON(in string) (string, error) {
	if len(in) == 0 {
		return "", nil
	}
	out, err := ghodssyaml.YAMLToJSON([]byte(in))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func JSONToYAML(in string) (string, error) {
	if len(in) == 0 {
		return "", nil
	}
	out, err := ghodssyaml.JSONToYAML([]byte(in))
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func FileExist(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func newString(value string) *string {
	return &value
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

func newEmptyString() *string {
	b := ""
	return &b
}

func IsVersionLessThan(ver1, ver2 string) (bool, error) {
	v1, err := gover.NewVersion(ver1)
	if err != nil {
		return false, err
	}
	v2, err := gover.NewVersion(ver2)
	if err != nil {
		return false, err
	}
	return v1.LessThan(v2), nil
}

func IsVersionGreaterThanOrEqual(ver1, ver2 string) (bool, error) {
	v1, err := gover.NewVersion(ver1)
	if err != nil {
		return false, err
	}
	v2, err := gover.NewVersion(ver2)
	if err != nil {
		return false, err
	}
	return v1.GreaterThanOrEqual(v2), nil
}

func sortVersions(list map[string]string) ([]*gover.Version, error) {
	var versions []*gover.Version
	for key := range list {
		v, err := gover.NewVersion(key)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	sort.Sort(gover.Collection(versions))
	return versions, nil
}

func getLatestVersion(list map[string]string) (string, error) {
	sorted, err := sortVersions(list)
	if err != nil {
		return "", err
	}

	return sorted[len(sorted)-1].Original(), nil
}

func structToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	relType := reflect.TypeOf(item)
	relValue := reflect.ValueOf(item)
	switch relType.Kind() {
	case reflect.Ptr:
		relValue = reflect.ValueOf(item).Elem()
		if !relValue.IsValid() {
			return res
		}
		relType = reflect.ValueOf(item).Elem().Type()
	}
	for i := 0; i < relType.NumField(); i++ {
		tags := strings.Split(relType.Field(i).Tag.Get("json"), ",")
		tag := tags[0]
		if tag != "" && tag != "-" {
			switch relType.Field(i).Type.Kind() {
			case reflect.Slice, reflect.Array:
				subtype := relValue.Field(i).Type().Elem()
				if subtype.Kind() == reflect.Struct {
					subvalue := reflect.ValueOf(relValue.Field(i).Interface())
					field := make([]interface{}, subvalue.Len())
					for i := 0; i < subvalue.Len(); i++ {
						field[i] = structToMap(subvalue.Index(i).Interface())
					}
					res[tag] = field
				} else {
					res[tag] = relValue.Field(i).Interface()
				}
			case reflect.Ptr:
				subvalue := reflect.ValueOf(relValue.Field(i).Interface())
				if subvalue.IsNil() {
					res[tag] = nil
					break
				}
				subtype := relValue.Field(i).Type().Elem()
				if subtype.Kind() == reflect.Struct {
					res[tag] = structToMap(relValue.Field(i).Interface())
				} else {
					res[tag] = relValue.Field(i).Interface()
				}
			case reflect.Struct:
				res[tag] = structToMap(relValue.Field(i).Interface())
			default:
				res[tag] = relValue.Field(i).Interface()
			}
		} else {
			if relType.Field(i).Type.Kind() == reflect.Struct {
				data := structToMap(relValue.Field(i).Interface())
				for i := range data {
					res[i] = data[i]
				}
			}
		}
	}

	return res
}
