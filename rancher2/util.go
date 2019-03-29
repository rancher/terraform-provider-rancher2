package rancher2

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	clusterProjectIDSeparator = ":"
	passDigits                = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	passDefaultLen            = 20
)

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
		return "", fmt.Errorf("[ERROR] problem encrypting password: %v", err)
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
	loginURL := url + "-public/localProviders/local?action=login"
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
		return response, fmt.Errorf("[ERROR] Doing post: URL is nil")
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

func NormalizeURL(url string) string {
	if url == "" {
		return ""
	}

	url = strings.TrimSuffix(url, "/")

	if !strings.HasSuffix(url, "/v3") {
		url = url + "/v3"
	}

	return url
}

func IsNotFound(err error) bool {
	return clientbase.IsNotFound(err)
}

func splitTokenID(token string) string {
	separator := ":"

	if strings.Contains(token, separator) {
		return token[0:strings.Index(token, separator)]
	}

	return token
}

func splitID(id string) (clusterID, resourceID string) {
	separator := ":"
	if strings.Contains(id, separator) {
		return id[0:strings.Index(id, separator)], id[strings.Index(id, separator)+1:]
	}
	return "", id
}

func clusterIDFromProjectID(projectID string) (string, error) {
	if projectID == "" || !strings.Contains(projectID, clusterProjectIDSeparator) {
		return "", fmt.Errorf("[ERROR] Getting clusted ID from project ID: Bad project id format %s", projectID)
	}

	return projectID[0:strings.Index(projectID, clusterProjectIDSeparator)], nil
}

func splitProjectID(id string) (clusterID, projectID string) {
	id = strings.TrimSuffix(id, clusterProjectIDSeparator)

	if strings.Contains(id, clusterProjectIDSeparator) {
		return id[0:strings.Index(id, clusterProjectIDSeparator)], id
	}

	return id, ""
}

func toArrayString(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = v.(string)
	}
	return out
}

func toArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func toMapString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for i, v := range in {
		out[i] = v.(string)
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
