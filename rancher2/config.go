package rancher2

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
	"golang.org/x/crypto/bcrypt"
)

const (
	rancher2ClientAPIVersion          = "/v3"
	rancher2CatalogAPIVersion         = "/v1"
	rancher2CatalogTypePrefix         = "catalog.cattle.io"
	rancher2ManagementV2TypePrefix    = "management.cattle.io"
	rancher2ReadyAnswer               = "pong"
	rancher2RetriesWait               = 5
	rancher2WaitFalseCond             = 120
	rancher2RKEK8sSystemImageVersion  = "2.3.0"
	rancher2NodeTemplateChangeVersion = "2.3.3" // Change node template id format
	rancher2TokeTTLMinutesVersion     = "2.4.6" // ttl token is readed in minutes
	rancher2TokeTTLMilisVersion       = "2.4.7" // ttl token is readed in miliseconds
	rancher2UILandingVersion          = "2.5.0" // ui landing option
	rancher2NodeTemplateNewPrefix     = "cattle-global-nt:nt-"
	rancher2DefaultTimeout            = "120s"
	rancher2DefaultLocalClusterID     = "local"
)

// Client are the client kind for a Rancher v3 API
type Client struct {
	Management *managementClient.Client
	CatalogV2  map[string]*clientbase.APIBaseClient
	Cluster    map[string]*clusterClient.Client
	Project    map[string]*projectClient.Client
}

// Config is the configuration parameters for a Rancher v3 API
type Config struct {
	TokenKey             string `json:"tokenKey"`
	URL                  string `json:"url"`
	CACerts              string `json:"cacert"`
	Insecure             bool   `json:"insecure"`
	Bootstrap            bool   `json:"bootstrap"`
	ClusterID            string `json:"clusterId"`
	ProjectID            string `json:"projectId"`
	Timeout              time.Duration
	RancherVersion       string
	K8SDefaultVersion    string
	K8SSupportedVersions []string
	Sync                 sync.Mutex
	Client               Client
}

// GetRancherVersion get Rancher server version
func (c *Config) GetRancherVersion() (string, error) {
	if len(c.RancherVersion) > 0 {
		return c.RancherVersion, nil
	}

	if c.Client.Management == nil {
		_, err := c.ManagementClient()
		if err != nil {
			return "", err
		}
	}

	version, err := c.Client.Management.Setting.ByID("server-version")
	if err != nil {
		return "", fmt.Errorf("[ERROR] Getting Rancher version: %s", err)
	}
	c.RancherVersion = version.Value

	return c.RancherVersion, nil
}

func (c *Config) isRancherReady() error {
	var err error
	var resp []byte
	url := c.URL + "/ping"
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		resp, err = DoGet(url, "", "", "", c.CACerts, c.Insecure)
		if err == nil && rancher2ReadyAnswer == string(resp) {
			return nil
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("Rancher is not ready: %v", err)
		}
	}
}

func (c *Config) waitForRancherLocalActive() error {
	client, err := c.ManagementClient()
	if err != nil {
		return err
	}
	clusterLocal, err := client.Cluster.ByID(rancher2DefaultLocalClusterID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			return nil
		}
		return fmt.Errorf("Waiting For local cluster. Getting cluster: %s", err)
	}
	_, err = c.WaitForClusterState(clusterLocal.ID, clusterActiveCondition, c.Timeout)
	return err
}

func (c *Config) getK8SDefaultVersion() (string, error) {
	if len(c.K8SDefaultVersion) > 0 {
		return c.K8SDefaultVersion, nil
	}
	var err error
	if c.Client.Management == nil {
		_, err = c.ManagementClient()
		if err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		k8sVer, err := c.Client.Management.Setting.ByID("k8s-version")
		if err == nil {
			c.K8SDefaultVersion = k8sVer.Value
			return c.K8SDefaultVersion, nil
		}
		if !IsServerError(err) && !IsForbidden(err) {
			return "", err
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return "", err
		}
	}
}

func (c *Config) getK8SVersions() ([]string, error) {
	if len(c.K8SSupportedVersions) > 0 {
		return c.K8SSupportedVersions, nil
	}

	if c.Client.Management == nil {
		_, err := c.ManagementClient()
		if err != nil {
			return nil, err
		}
	}

	if ok, _ := c.IsRancherVersionLessThan(rancher2RKEK8sSystemImageVersion); ok {
		return nil, nil
	}

	RKEK8sSystemImageCollection, err := c.Client.Management.RkeK8sSystemImage.ListAll(NewListOpts(nil))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Listing RKE K8s System Images: %s", err)
	}
	versions := make([]*version.Version, 0, len(RKEK8sSystemImageCollection.Data))
	for _, RKEK8sSystem := range RKEK8sSystemImageCollection.Data {
		v, _ := version.NewVersion(RKEK8sSystem.Name)
		versions = append(versions, v)

	}
	sort.Sort(sort.Reverse(version.Collection(versions)))
	for i := range versions {
		c.K8SSupportedVersions = append(c.K8SSupportedVersions, "v"+versions[i].String())
	}
	return c.K8SSupportedVersions, nil
}

// Fix breaking API change https://github.com/rancher/rancher/pull/23718
func (c *Config) fixNodeTemplateID(id string) string {
	if ok, _ := c.IsRancherVersionGreaterThanOrEqual(rancher2NodeTemplateChangeVersion); ok && len(id) > 0 {
		if !strings.HasPrefix(id, rancher2NodeTemplateNewPrefix) {
			id = strings.Replace(id, ":", "-", -1)
			id = rancher2NodeTemplateNewPrefix + id
		}
	}
	return id
}

func (c *Config) IsRancherVersionGreaterThanOrEqualAndLessThan(ver1, ver2 string) (bool, error) {
	_, err := c.GetRancherVersion()
	if err != nil {
		return false, fmt.Errorf("[ERROR] getting rancher server version")
	}
	greaterOrEqualThan, err := IsVersionGreaterThanOrEqual(c.RancherVersion, ver1)
	if err != nil {
		return false, err
	}
	lessThan, err := IsVersionLessThan(c.RancherVersion, ver2)
	if err != nil {
		return false, err
	}
	return (greaterOrEqualThan && lessThan), nil
}

func (c *Config) IsRancherVersionLessThan(ver string) (bool, error) {
	if len(ver) == 0 {
		return false, fmt.Errorf("[ERROR] version is nil")
	}
	_, err := c.GetRancherVersion()
	if err != nil {
		return false, fmt.Errorf("[ERROR] getting rancher server version")
	}
	return IsVersionLessThan(c.RancherVersion, ver)
}

func (c *Config) IsRancherVersionGreaterThanOrEqual(ver string) (bool, error) {
	if len(ver) == 0 {
		return false, fmt.Errorf("[ERROR] version is nil")
	}
	_, err := c.GetRancherVersion()
	if err != nil {
		return false, fmt.Errorf("[ERROR] getting rancher server version")
	}
	return IsVersionGreaterThanOrEqual(c.RancherVersion, ver)
}

// UpdateToken update tokenkey and restart client connections
func (c *Config) UpdateToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("token is nil")
	}
	c.TokenKey = token

	return c.RestartClients()
}

// RestartClients connections
func (c *Config) RestartClients() error {
	c.Sync.Lock()
	if c.Client.Management != nil {
		c.Client.Management = nil
	}
	c.Sync.Unlock()

	_, err := c.ManagementClient()
	if err != nil {
		return err
	}
	c.Client.Cluster = map[string]*clusterClient.Client{}
	c.Client.Project = map[string]*projectClient.Client{}
	c.Client.CatalogV2 = map[string]*clientbase.APIBaseClient{}

	return nil
}

// ManagementClient creates a Rancher client scoped to the management API
func (c *Config) ManagementClient() (*managementClient.Client, error) {
	c.Sync.Lock()
	defer c.Sync.Unlock()

	if c.Client.Management != nil {
		return c.Client.Management, nil
	}

	err := c.isRancherReady()
	if err != nil {
		return nil, err
	}

	// Setup the management client
	options := c.CreateClientOpts()
	options.URL = options.URL + rancher2ClientAPIVersion
	mClient, err := managementClient.NewClient(options)
	if err != nil {
		return nil, err
	}
	c.Client.Management = mClient

	rancher2ClusterRKEK8SDefaultVersion, err = c.getK8SDefaultVersion()
	if err != nil {
		return nil, err
	}
	rancher2ClusterRKEK8SVersions, err = c.getK8SVersions()
	if err != nil {
		return nil, err
	}

	return c.Client.Management, nil
}

// CatalogV2Client creates a Rancher client scoped to a Cluster API
func (c *Config) CatalogV2Client(id string) (*clientbase.APIBaseClient, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Catalog V2 Client: cluster ID is nil")
	}

	c.Sync.Lock()
	defer c.Sync.Unlock()

	if c.Client.CatalogV2 == nil {
		c.Client.CatalogV2 = map[string]*clientbase.APIBaseClient{}
	}

	if c.Client.CatalogV2[id] != nil {
		return c.Client.CatalogV2[id], nil
	}

	err := c.isRancherReady()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	// Setup the cluster client
	options := c.CreateClientOpts()
	options.URL = options.URL + "/k8s/clusters/" + id + rancher2CatalogAPIVersion
	for {
		cli, err := clientbase.NewAPIClient(options)
		if err == nil {
			c.Client.CatalogV2[id] = &cli
			return c.Client.CatalogV2[id], nil
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, err
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout getting Catalog V2 Client at cluster ID %s: %v", id, err)
		}
	}
}

// ClusterClient creates a Rancher client scoped to a Cluster API
func (c *Config) ClusterClient(id string) (*clusterClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Cluster Client: cluster ID is nil")
	}

	c.Sync.Lock()
	defer c.Sync.Unlock()

	if c.Client.Cluster == nil {
		c.Client.Cluster = map[string]*clusterClient.Client{}
	}

	if c.Client.Cluster[id] != nil {
		return c.Client.Cluster[id], nil
	}

	err := c.isRancherReady()
	if err != nil {
		return nil, err
	}

	// Setup the cluster client
	options := c.CreateClientOpts()
	options.URL = options.URL + rancher2ClientAPIVersion + "/clusters/" + id
	cClient, err := clusterClient.NewClient(options)
	if err != nil {
		return nil, err
	}
	c.Client.Cluster[id] = cClient

	return c.Client.Cluster[id], nil
}

// ProjectClient creates a Rancher client scoped to a Project API
func (c *Config) ProjectClient(id string) (*projectClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Project Client: project ID is nil")
	}

	c.Sync.Lock()
	defer c.Sync.Unlock()

	if c.Client.Project == nil {
		c.Client.Project = map[string]*projectClient.Client{}
	}

	if c.Client.Project[id] != nil {
		return c.Client.Project[id], nil
	}

	err := c.isRancherReady()
	if err != nil {
		return nil, err
	}

	// Setup the project client
	options := c.CreateClientOpts()
	options.URL = options.URL + rancher2ClientAPIVersion + "/projects/" + id
	pClient, err := projectClient.NewClient(options)
	if err != nil {
		return nil, err
	}

	c.Client.Project[id] = pClient

	return c.Client.Project[id], nil
}

func (c *Config) NormalizeURL() error {
	c.Sync.Lock()
	defer c.Sync.Unlock()
	url, err := NormalizeURL(c.URL)
	if err != nil {
		return err
	}
	c.URL = url
	return nil
}

func (c *Config) CreateClientOpts() *clientbase.ClientOpts {
	options := &clientbase.ClientOpts{
		URL:      c.URL,
		TokenKey: c.TokenKey,
		CACerts:  c.CACerts,
		Insecure: c.Insecure,
	}
	return options
}

func (c *Config) GetProjectRoleTemplateBindingsByProjectID(projectID string) ([]managementClient.ProjectRoleTemplateBinding, error) {
	if projectID == "" {
		return nil, fmt.Errorf("[ERROR] Project ID is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"ProjectID": projectID}
	listOpts := NewListOpts(filters)

	collection, err := client.ProjectRoleTemplateBinding.List(listOpts)
	if err != nil {
		return nil, err
	}

	data := collection.Data

	// Paginating data if needed
	if collection.Pagination.Partial {
		for collection, err = collection.Next(); err != nil && collection != nil; collection, err = collection.Next() {
			data = append(data, collection.Data...)
		}
	}

	return data, err
}

func (c *Config) IsProjectDefault(project *managementClient.Project) bool {
	if project == nil {
		return false
	}

	for k, v := range project.Labels {
		if k == projectDefaultLabel && v == "true" {
			return true
		}
	}

	return false
}

func (c *Config) IsProjectSystem(project *managementClient.Project) bool {
	if project == nil {
		return false
	}

	for k, v := range project.Labels {
		if k == projectSystemLabel && v == "true" {
			return true
		}
	}

	return false
}

func (c *Config) GetProjectByName(name, clusterID string) (*managementClient.Project, error) {
	if name == "" {
		return nil, fmt.Errorf("[ERROR] Project name is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"clusterId": clusterID, "name": name}
	listOpts := NewListOpts(filters)

	collection, err := client.Project.List(listOpts)
	if err != nil {
		return nil, err
	}

	// Returning first project name matching
	if len(collection.Data) > 0 {
		return &collection.Data[0], nil
	}
	return nil, fmt.Errorf("[ERROR] Project %s on cluster %s not found", name, clusterID)
}

func (c *Config) GetProjectIDByName(name, clusterID string) (string, error) {
	if name == "" {
		return "", nil
	}

	project, err := c.GetProjectByName(name, clusterID)
	if err != nil {
		return "", err
	}
	return project.ID, nil
}

func (c *Config) GetProjectNameByID(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	project, err := c.GetProjectByID(id)
	if err != nil {
		return "", err
	}

	return project.Name, nil
}

func (c *Config) GetProjectByID(id string) (*managementClient.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("Project id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.Project.ByID(id)
}

func (c *Config) ProjectExist(id string) error {
	_, err := c.GetProjectByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetGlobalRoleByID(id string) (*managementClient.GlobalRole, error) {
	if id == "" {
		return nil, fmt.Errorf("Global role id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.GlobalRole.ByID(id)
}

func (c *Config) GlobalRoleExist(id string) error {
	_, err := c.GetGlobalRoleByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetRoleTemplateByID(id string) (*managementClient.RoleTemplate, error) {
	if id == "" {
		return nil, fmt.Errorf("Role template id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.RoleTemplate.ByID(id)
}

func (c *Config) RoleTemplateExist(id string) error {
	_, err := c.GetRoleTemplateByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetClusterProjects(id string) ([]managementClient.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Cluster id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"clusterId": id}
	listOpts := NewListOpts(filters)

	collection, err := client.Project.List(listOpts)
	if err != nil {
		return nil, err
	}

	return collection.Data, nil
}

func (c *Config) GetClusterSpecialProjectsID(id string) (string, string, error) {
	if id == "" {
		return "", "", fmt.Errorf("[ERROR] Cluster id is nil")
	}

	projects, err := c.GetClusterProjects(id)
	if err != nil {
		return "", "", err
	}

	found := 0
	defaultProjectID := ""
	systemProjectID := ""
	for _, project := range projects {
		if c.IsProjectDefault(&project) {
			found++
			defaultProjectID = project.ID
		}
		if c.IsProjectSystem(&project) {
			found++
			systemProjectID = project.ID
		}
		if found == 2 {
			return defaultProjectID, systemProjectID, nil
		}
	}

	return defaultProjectID, systemProjectID, nil
}

func (c *Config) GetClusterNodes(id string) ([]managementClient.Node, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Cluster id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"clusterId": id}
	listOpts := NewListOpts(filters)

	collection, err := client.Node.List(listOpts)
	if err != nil {
		return nil, err
	}

	return collection.Data, nil
}

func (c *Config) GetClusterByName(name string) (*managementClient.Cluster, error) {
	if name == "" {
		return nil, fmt.Errorf("[ERROR] Cluster name is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"name": name}
	listOpts := NewListOpts(filters)

	collection, err := client.Cluster.List(listOpts)
	if err != nil {
		return nil, err
	}

	// Cluster names are globally unique
	if len(collection.Data) > 0 {
		return &collection.Data[0], nil
	}

	return nil, fmt.Errorf("[ERROR] Cluster %s not found", name)
}

func (c *Config) GetClusterIDByName(name string) (string, error) {
	if name == "" {
		return "", nil
	}

	cluster, err := c.GetClusterByName(name)
	if err != nil {
		return "", err
	}
	return cluster.ID, nil
}

func (c *Config) GetClusterByID(id string) (*managementClient.Cluster, error) {
	if id == "" {
		return nil, fmt.Errorf("Cluster id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.Cluster.ByID(id)
}

func (c *Config) WaitForClusterState(clusterID, state string, interval time.Duration) (*managementClient.Cluster, error) {
	if clusterID == "" || state == "" {
		return nil, fmt.Errorf("Cluster ID and/or state is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for {
		obj, err := c.GetClusterByID(clusterID)
		if err != nil {
			if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
				return nil, fmt.Errorf("Getting cluster ID (%s): %v", clusterID, err)
			}
		}
		if obj != nil {
			for i := range obj.Conditions {
				if obj.Conditions[i].Type == state {
					// Status of the condition, one of True, False, Unknown.
					if obj.Conditions[i].Status == "Unknown" {
						break
					}
					if obj.Conditions[i].Status == "True" {
						return obj, nil
					}
					// When cluster condition is false, retrying if it has been updated for last rancher2WaitFalseCond seconds
					lastUpdate, err := time.Parse(time.RFC3339, obj.Conditions[i].LastUpdateTime)
					if err == nil && time.Since(lastUpdate) < rancher2WaitFalseCond*time.Second {
						break
					}
					return nil, fmt.Errorf("Cluster ID %s: %s", clusterID, obj.Conditions[i].Message)
				}
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout waiting for cluster ID %s", clusterID)
		}
	}
}

func (c *Config) getObjectV2ByID(clusterID, id, APIType string, resp interface{}) error {
	if id == "" {
		return fmt.Errorf("Object V2 id is nil")
	}
	if resp == nil {
		return fmt.Errorf("Object V2 response is nil")
	}
	if len(APIType) == 0 {
		return fmt.Errorf("Object API V2 type is nil")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err = client.ByID(APIType, id, resp)
		if err == nil {
			return nil
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) {
			return err
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("Timeout getting object V2 ID %s at cluster ID %s: %v", id, clusterID, err)
		}
	}
}

func (c *Config) GetSettingV2ByID(id string) (*SettingV2, error) {
	resp := &SettingV2{}
	err := c.getObjectV2ByID("local", id, settingV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting Setting V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func (c *Config) createObjectV2(clusterID string, APIType string, obj, resp interface{}) error {
	if resp == nil || obj == nil {
		return fmt.Errorf("Object V2 and/or response is nil")
	}
	if len(APIType) == 0 {
		return fmt.Errorf("Object API V2 type is nil")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return err
	}
	return client.Create(APIType, obj, resp)
}

func (c *Config) deleteObjectV2(clusterID string, resource *types.Resource) error {
	if resource == nil {
		return fmt.Errorf("Object V2 id is nil")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return err
	}
	return client.Delete(resource)
}

func (c *Config) updateObjectV2(clusterID, id, APIType string, update, resp interface{}) error {
	if id == "" {
		return fmt.Errorf("Object V2 id is nil")
	}
	if update == nil {
		return fmt.Errorf("Object V2 update is nil")
	}
	if len(APIType) == 0 {
		return fmt.Errorf("Object API V2 type is nil")
	}

	if resp == nil {
		resp = map[string]interface{}{}
	}

	resource := &types.Resource{}
	err := c.getObjectV2ByID(clusterID, id, APIType, resource)
	if err != nil {
		return err
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return err
	}
	return client.Update(APIType, resource, update, resp)
}

func (c *Config) UpdateClusterByID(cluster *managementClient.Cluster, update map[string]interface{}) (*managementClient.Cluster, error) {
	if cluster == nil {
		return nil, fmt.Errorf("[ERROR] Updating cluster: Cluster is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.Cluster.Update(cluster, update)
}

func (c *Config) checkClusterCondition(id, condition string) (bool, *managementClient.Cluster, error) {
	if len(condition) == 0 {
		return false, nil, fmt.Errorf("checkClusterCondition: Cluster condition is empty")
	}
	obj, err := c.GetClusterByID(id)
	if err != nil {
		return false, nil, err
	}

	for i := range obj.Conditions {
		if obj.Conditions[i].Type == condition {
			if obj.Conditions[i].Status == "True" {
				return true, obj, nil
			}
			return false, obj, nil
		}
	}

	return false, obj, nil
}

func (c *Config) isClusterActive(id string) (bool, *managementClient.Cluster, error) {
	return c.checkClusterCondition(id, clusterActiveCondition)
}

func (c *Config) isClusterConnected(id string) (bool, *managementClient.Cluster, error) {
	return c.checkClusterCondition(id, clusterConnectedCondition)
}

func (c *Config) isClusterMonitoringEnabledCondition(id string) (bool, *managementClient.Cluster, error) {
	return c.checkClusterCondition(id, clusterMonitoringEnabledCondition)
}

func (c *Config) isClusterAlertingEnabledCondition(id string) (bool, *managementClient.Cluster, error) {
	return c.checkClusterCondition(id, clusterAlertingEnabledCondition)
}

func (c *Config) ClusterExist(id string) error {
	_, err := c.GetClusterByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetClusterRegistrationTokenByID(id string) (*managementClient.ClusterRegistrationToken, error) {
	if id == "" {
		return nil, fmt.Errorf("Cluster Regitration Token id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.ClusterRegistrationToken.ByID(id)
}

func (c *Config) ClusterRegistrationTokenExist(id string) error {
	_, err := c.GetClusterRegistrationTokenByID(id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) CheckAuthConfigEnabled(id string) error {
	if id == "" {
		return fmt.Errorf("Auth config id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	listOpts := NewListOpts(nil)
	auths, err := client.AuthConfig.List(listOpts)
	if err != nil {
		return err
	}

	for _, auth := range auths.Data {
		if auth.Enabled {
			if auth.ID != id && auth.ID != "local" {
				return fmt.Errorf("%s provider is already enabled", auth.ID)
			}
		}
	}

	return nil
}

func (c *Config) GetAuthConfig(in *managementClient.AuthConfig) (interface{}, error) {
	resp, err := getAuthConfigObject(in.Type)
	if err != nil {
		return nil, err
	}

	link := "self"

	resource := types.Resource{}
	resource.Links = in.Links

	err = c.Client.Management.GetLink(resource, link, resp)
	if err != nil {
		return nil, fmt.Errorf("Error getting Auth Config [%s] %s", resource.Links[link], err)
	}

	return resp, nil
}

func (c *Config) UpdateAuthConfig(url string, createObj interface{}, respObject interface{}) error {
	return c.Client.Management.Ops.DoModify("PUT", url, createObj, respObject)
}

func (c *Config) GetUserByName(name string) (*managementClient.User, error) {
	if name == "" {
		return nil, fmt.Errorf("[ERROR] Username is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"username": name}
	listOpts := NewListOpts(filters)

	collection, err := client.User.List(listOpts)
	if err != nil {
		return nil, err
	}

	// Usernames are globally unique
	if len(collection.Data) > 0 {
		return &collection.Data[0], nil
	}
	return nil, fmt.Errorf("[ERROR] Username %s not found", name)
}

func (c *Config) GetUserIDByName(name string) (string, error) {
	if name == "" {
		return "", nil
	}

	user, err := c.GetUserByName(name)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (c *Config) activateDriver(id string, interval time.Duration) error {
	if id == googleConfigDriver {
		return c.activateKontainerDriver(id, interval)
	}

	return c.activateNodeDriver(id, interval)
}

func (c *Config) activateKontainerDriver(id string, interval time.Duration) error {
	if id == "" {
		return fmt.Errorf("[ERROR] Node Driver id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for updated := false; ; {
		driver, err := client.KontainerDriver.ByID(id)
		if err != nil && !IsServerError(err) && !IsUnknownSchemaType(err) && !IsNotFound(err) && !IsForbidden(err) {
			return fmt.Errorf("[ERROR] Getting Node Driver %s: %v", id, err)
		}
		if driver != nil {
			if driver.State == "active" {
				return nil
			}

			if !updated {
				err = client.KontainerDriver.ActionActivate(driver)
				if err == nil {
					updated = true
					continue
				}
				if !IsServerError(err) && !IsUnknownSchemaType(err) {
					return fmt.Errorf("[ERROR] Activating Node Driver %s: %v", id, err)
				}
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("Timeout activating Node Driver %s: %v", id, err)
		}
	}
}

func (c *Config) activateNodeDriver(id string, interval time.Duration) error {
	if id == "" {
		return fmt.Errorf("[ERROR] Node Driver id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for updated := false; ; {
		driver, err := client.NodeDriver.ByID(id)
		if err != nil && !IsServerError(err) && !IsUnknownSchemaType(err) && !IsNotFound(err) && !IsForbidden(err) {
			return fmt.Errorf("[ERROR] Getting Node Driver %s: %v", id, err)
		}
		if driver != nil {
			if driver.State == "active" {
				return nil
			}

			if !updated {
				driver, err = client.NodeDriver.ActionActivate(driver)
				if err == nil {
					updated = true
					continue
				}
				if !IsServerError(err) && !IsUnknownSchemaType(err) {
					return fmt.Errorf("[ERROR] Activating Node Driver %s: %v", id, err)
				}
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return fmt.Errorf("Timeout activating Node Driver %s: %v", id, err)
		}
	}
}

func (c *Config) UserPasswordChanged(user *managementClient.User, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	// Password has changed
	if err != nil {
		return true
	}

	return false
}

func (c *Config) SetUserPasswordByName(username, pass string) (bool, string, *managementClient.User, error) {
	if len(username) == 0 {
		return false, "", nil, fmt.Errorf("[ERROR] Setting user password: Username is nil")
	}

	// Generating rancdom password if nil
	if len(pass) == 0 {
		pass = GetRandomPass(passDefaultLen)
	}

	user, err := c.GetUserByName(username)
	if err != nil {
		return false, "", nil, err
	}

	changed, newUser, err := c.SetUserPassword(user, pass)

	return changed, pass, newUser, err
}

func (c *Config) SetUserPassword(user *managementClient.User, pass string) (bool, *managementClient.User, error) {
	changed := c.UserPasswordChanged(user, pass)
	if !changed {
		return changed, user, nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return false, nil, err
	}

	password := &managementClient.SetPasswordInput{NewPassword: pass}
	newUser, err := client.User.ActionSetpassword(user, password)
	if err != nil {
		return false, nil, fmt.Errorf("[ERROR] Setting %s password: %s", user.Username, err)
	}

	return changed, newUser, nil
}

// GenerateUserToken generates token with ttl measured in seconds
func (c *Config) GenerateUserToken(username, desc string, ttl int) (string, string, error) {
	user, err := c.GetUserByName(username)
	if err != nil {
		return "", "", err
	}

	client, err := c.ManagementClient()
	if err != nil {
		return "", "", err
	}

	tokenTTL := int64(ttl)
	if tokenTTL > 0 {
		tokenTTL = tokenTTL * 1000
	}

	token := &managementClient.Token{
		UserID:      user.ID,
		Description: desc,
		TTLMillis:   tokenTTL,
	}

	newToken, err := client.Token.Create(token)
	if err != nil {
		return "", "", fmt.Errorf("[ERROR] Creating Admin token: %s", err)
	}

	return newToken.ID, newToken.Token, nil
}

func (c *Config) IsTokenExpired(id string) (bool, error) {
	if len(id) == 0 {
		return true, nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return false, err
	}

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return token.Expired, nil
}

func (c *Config) DeleteToken(id string) error {
	if len(id) == 0 {
		return nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			return nil
		}
		return err
	}

	return client.Token.Delete(token)
}

func (c *Config) GetSetting(name string) (*managementClient.Setting, error) {
	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	return client.Setting.ByID(name)
}

func (c *Config) SetSetting(name, value string) error {
	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	setting, err := c.GetSetting(name)
	if err != nil {
		return fmt.Errorf("[ERROR] Getting setting %s: %s", name, err)
	}

	update := map[string]interface{}{
		"value": value,
	}

	_, err = client.Setting.Update(setting, update)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating setting %s: %s", name, err)
	}

	return nil
}

func (c *Config) GetSettingValue(name string) (string, error) {
	setting, err := c.GetSetting(name)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Getting setting %s: %s", name, err)
	}

	return setting.Value, nil
}

func (c *Config) GetCatalogByName(name, scope string) (interface{}, error) {
	if len(name) == 0 || len(scope) == 0 {
		return nil, fmt.Errorf("[ERROR] Name nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"name": name}
	listOpts := NewListOpts(filters)

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.List(listOpts)
	case catalogScopeGlobal:
		return client.Catalog.List(listOpts)
	case catalogScopeProject:
		return client.ProjectCatalog.List(listOpts)
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func (c *Config) GetCatalog(id, scope string) (interface{}, error) {
	if len(id) == 0 || len(scope) == 0 {
		return nil, fmt.Errorf("[ERROR] Id nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.ByID(id)
	case catalogScopeGlobal:
		return client.Catalog.ByID(id)
	case catalogScopeProject:
		return client.ProjectCatalog.ByID(id)
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func (c *Config) CreateCatalog(scope string, catalog interface{}) (interface{}, error) {
	if catalog == nil || len(scope) == 0 {
		return nil, fmt.Errorf("[ERROR] Catalog nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.Create(catalog.(*managementClient.ClusterCatalog))
	case catalogScopeGlobal:
		return client.Catalog.Create(catalog.(*managementClient.Catalog))
	case catalogScopeProject:
		return client.ProjectCatalog.Create(catalog.(*managementClient.ProjectCatalog))
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func (c *Config) UpdateCatalog(scope string, catalog interface{}, update map[string]interface{}) (interface{}, error) {
	if catalog == nil || len(scope) == 0 {
		return nil, fmt.Errorf("[ERROR] Catalog nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.Update(catalog.(*managementClient.ClusterCatalog), update)
	case catalogScopeGlobal:
		return client.Catalog.Update(catalog.(*managementClient.Catalog), update)
	case catalogScopeProject:
		return client.ProjectCatalog.Update(catalog.(*managementClient.ProjectCatalog), update)
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func (c *Config) DeleteCatalog(scope string, catalog interface{}) error {
	if catalog == nil || len(scope) == 0 {
		return fmt.Errorf("[ERROR] Catalog nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.Delete(catalog.(*managementClient.ClusterCatalog))
	case catalogScopeGlobal:
		return client.Catalog.Delete(catalog.(*managementClient.Catalog))
	case catalogScopeProject:
		return client.ProjectCatalog.Delete(catalog.(*managementClient.ProjectCatalog))
	default:
		return fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func (c *Config) RefreshCatalog(scope string, catalog interface{}) (*managementClient.CatalogRefresh, error) {
	if catalog == nil || len(scope) == 0 {
		return nil, fmt.Errorf("[ERROR] Catalog nor scope can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	switch scope {
	case catalogScopeCluster:
		return client.ClusterCatalog.ActionRefresh(catalog.(*managementClient.ClusterCatalog))
	case catalogScopeGlobal:
		return client.Catalog.ActionRefresh(catalog.(*managementClient.Catalog))
	case catalogScopeProject:
		return client.ProjectCatalog.ActionRefresh(catalog.(*managementClient.ProjectCatalog))
	default:
		return nil, fmt.Errorf("[ERROR] Unsupported scope on catalog: %s", scope)
	}
}

func getAuthConfigObject(kind string) (interface{}, error) {
	switch kind {
	case managementClient.ActiveDirectoryConfigType:
		return &managementClient.ActiveDirectoryConfig{}, nil
	case managementClient.ADFSConfigType:
		return &managementClient.ADFSConfig{}, nil
	case managementClient.AzureADConfigType:
		return &managementClient.AzureADConfig{}, nil
	case managementClient.FreeIpaConfigType:
		return &managementClient.LdapConfig{}, nil
	case managementClient.GithubConfigType:
		return &managementClient.GithubConfig{}, nil
	case managementClient.KeyCloakConfigType:
		return &managementClient.KeyCloakConfig{}, nil
	case managementClient.OKTAConfigType:
		return &managementClient.OKTAConfig{}, nil
	case managementClient.OpenLdapConfigType:
		return &managementClient.LdapConfig{}, nil
	case managementClient.PingConfigType:
		return &managementClient.PingConfig{}, nil
	default:
		return nil, fmt.Errorf("[ERROR] Auth config type %s not supported", kind)
	}
}

func (c *Config) GetRegistryByFilters(filters map[string]interface{}) (interface{}, error) {
	if filters == nil || len(filters["name"].(string)) == 0 || len(filters["projectId"].(string)) == 0 {
		return nil, fmt.Errorf("[ERROR] Name nor project_id can't be nil")
	}

	client, err := c.ProjectClient(filters["projectId"].(string))
	if err != nil {
		return nil, err
	}

	listOpts := NewListOpts(filters)

	if filters["namespaceId"] != nil {
		return client.NamespacedDockerCredential.List(listOpts)
	}

	return client.DockerCredential.List(listOpts)
}

func (c *Config) GetRegistry(id, projectID, namespaceID string) (interface{}, error) {
	if len(id) == 0 || len(projectID) == 0 {
		return nil, fmt.Errorf("[ERROR] Id nor project id can't be nil")
	}

	client, err := c.ProjectClient(projectID)
	if err != nil {
		return nil, err
	}

	if len(namespaceID) > 0 {
		return client.NamespacedDockerCredential.ByID(id)
	}

	return client.DockerCredential.ByID(id)
}

func (c *Config) createDockerCredential(registry *projectClient.DockerCredential) (*projectClient.DockerCredential, error) {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.DockerCredential.Create(registry)
}

func (c *Config) createNamespacedDockerCredential(registry *projectClient.NamespacedDockerCredential) (*projectClient.NamespacedDockerCredential, error) {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedDockerCredential.Create(registry)
}

func (c *Config) CreateRegistry(registry interface{}) (interface{}, error) {
	if registry == nil {
		return nil, fmt.Errorf("[ERROR] Registry can't be nil")
	}

	switch t := registry.(type) {
	case *projectClient.NamespacedDockerCredential:
		return c.createNamespacedDockerCredential(registry.(*projectClient.NamespacedDockerCredential))
	case *projectClient.DockerCredential:
		return c.createDockerCredential(registry.(*projectClient.DockerCredential))
	default:
		return nil, fmt.Errorf("[ERROR] Registry type %s isn't supported", t)
	}
}

func (c *Config) updateDockerCredential(registry *projectClient.DockerCredential, update map[string]interface{}) (*projectClient.DockerCredential, error) {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.DockerCredential.Update(registry, update)
}

func (c *Config) updateNamespacedDockerCredential(registry *projectClient.NamespacedDockerCredential, update map[string]interface{}) (*projectClient.NamespacedDockerCredential, error) {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedDockerCredential.Update(registry, update)
}

func (c *Config) UpdateRegistry(registry interface{}, update map[string]interface{}) (interface{}, error) {
	if registry == nil {
		return nil, fmt.Errorf("[ERROR] Registry can't be nil")
	}

	switch t := registry.(type) {
	case *projectClient.NamespacedDockerCredential:
		return c.updateNamespacedDockerCredential(registry.(*projectClient.NamespacedDockerCredential), update)
	case *projectClient.DockerCredential:
		return c.updateDockerCredential(registry.(*projectClient.DockerCredential), update)
	default:
		return nil, fmt.Errorf("[ERROR] Registry type %s isn't supported", t)
	}
}

func (c *Config) deleteDockerCredential(registry *projectClient.DockerCredential) error {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return err
	}
	return client.DockerCredential.Delete(registry)
}

func (c *Config) deleteNamespacedDockerCredential(registry *projectClient.NamespacedDockerCredential) error {
	client, err := c.ProjectClient(registry.ProjectID)
	if err != nil {
		return err
	}
	return client.NamespacedDockerCredential.Delete(registry)
}

func (c *Config) DeleteRegistry(registry interface{}) error {
	if registry == nil {
		return fmt.Errorf("[ERROR] Registry can't be nil")
	}

	switch t := registry.(type) {
	case *projectClient.NamespacedDockerCredential:
		return c.deleteNamespacedDockerCredential(registry.(*projectClient.NamespacedDockerCredential))
	case *projectClient.DockerCredential:
		return c.deleteDockerCredential(registry.(*projectClient.DockerCredential))
	default:
		return fmt.Errorf("[ERROR] Registry type %s isn't supported", t)
	}
}

func (c *Config) GetSecretByFilters(filters map[string]interface{}) (interface{}, error) {
	if filters == nil || len(filters["name"].(string)) == 0 || len(filters["projectId"].(string)) == 0 {
		return nil, fmt.Errorf("[ERROR] Name nor project_id can't be nil")
	}

	client, err := c.ProjectClient(filters["projectId"].(string))
	if err != nil {
		return nil, err
	}

	listOpts := NewListOpts(filters)

	if filters["namespaceId"] != nil {
		return client.NamespacedSecret.List(listOpts)
	}

	return client.Secret.List(listOpts)
}

func (c *Config) GetSecret(id, projectID, namespaceID string) (interface{}, error) {
	if len(id) == 0 || len(projectID) == 0 {
		return nil, fmt.Errorf("[ERROR] Id nor project id can't be nil")
	}

	client, err := c.ProjectClient(projectID)
	if err != nil {
		return nil, err
	}

	if len(namespaceID) > 0 {
		return client.NamespacedSecret.ByID(id)
	}

	return client.Secret.ByID(id)
}

func (c *Config) createSecret(secret *projectClient.Secret) (*projectClient.Secret, error) {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.Secret.Create(secret)
}

func (c *Config) createNamespacedSecret(secret *projectClient.NamespacedSecret) (*projectClient.NamespacedSecret, error) {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedSecret.Create(secret)
}

func (c *Config) CreateSecret(secret interface{}) (interface{}, error) {
	if secret == nil {
		return nil, fmt.Errorf("[ERROR] Secret can't be nil")
	}

	switch t := secret.(type) {
	case *projectClient.NamespacedSecret:
		return c.createNamespacedSecret(secret.(*projectClient.NamespacedSecret))
	case *projectClient.Secret:
		return c.createSecret(secret.(*projectClient.Secret))
	default:
		return nil, fmt.Errorf("[ERROR] Secret type %s isn't supported", t)
	}
}

func (c *Config) updateSecret(secret *projectClient.Secret, update map[string]interface{}) (*projectClient.Secret, error) {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.Secret.Update(secret, update)
}

func (c *Config) updateNamespacedSecret(secret *projectClient.NamespacedSecret, update map[string]interface{}) (*projectClient.NamespacedSecret, error) {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedSecret.Update(secret, update)
}

func (c *Config) UpdateSecret(secret interface{}, update map[string]interface{}) (interface{}, error) {
	if secret == nil {
		return nil, fmt.Errorf("[ERROR] Secret can't be nil")
	}

	switch t := secret.(type) {
	case *projectClient.NamespacedSecret:
		return c.updateNamespacedSecret(secret.(*projectClient.NamespacedSecret), update)
	case *projectClient.Secret:
		return c.updateSecret(secret.(*projectClient.Secret), update)
	default:
		return nil, fmt.Errorf("[ERROR] Secret type %s isn't supported", t)
	}
}

func (c *Config) deleteSecret(secret *projectClient.Secret) error {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return err
	}
	return client.Secret.Delete(secret)
}

func (c *Config) deleteNamespacedSecret(secret *projectClient.NamespacedSecret) error {
	client, err := c.ProjectClient(secret.ProjectID)
	if err != nil {
		return err
	}
	return client.NamespacedSecret.Delete(secret)
}

func (c *Config) DeleteSecret(secret interface{}) error {
	if secret == nil {
		return fmt.Errorf("[ERROR] Secret can't be nil")
	}

	switch t := secret.(type) {
	case *projectClient.NamespacedSecret:
		return c.deleteNamespacedSecret(secret.(*projectClient.NamespacedSecret))
	case *projectClient.Secret:
		return c.deleteSecret(secret.(*projectClient.Secret))
	default:
		return fmt.Errorf("[ERROR] Secret type %s isn't supported", t)
	}
}

func (c *Config) GetCertificateByFilters(filters map[string]interface{}) (interface{}, error) {
	if filters == nil || len(filters["name"].(string)) == 0 || len(filters["projectId"].(string)) == 0 {
		return nil, fmt.Errorf("[ERROR] Name nor project_id can't be nil")
	}

	client, err := c.ProjectClient(filters["projectId"].(string))
	if err != nil {
		return nil, err
	}

	listOpts := NewListOpts(filters)

	if filters["namespaceId"] != nil {
		return client.NamespacedCertificate.List(listOpts)
	}

	return client.Certificate.List(listOpts)
}

func (c *Config) GetCertificate(id, projectID, namespaceID string) (interface{}, error) {
	if len(id) == 0 || len(projectID) == 0 {
		return nil, fmt.Errorf("[ERROR] Id nor project id can't be nil")
	}

	client, err := c.ProjectClient(projectID)
	if err != nil {
		return nil, err
	}

	if len(namespaceID) > 0 {
		return client.NamespacedCertificate.ByID(id)
	}

	return client.Certificate.ByID(id)
}

func (c *Config) createCertificate(cert *projectClient.Certificate) (*projectClient.Certificate, error) {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.Certificate.Create(cert)
}

func (c *Config) createNamespacedCertificate(cert *projectClient.NamespacedCertificate) (*projectClient.NamespacedCertificate, error) {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedCertificate.Create(cert)
}

func (c *Config) CreateCertificate(cert interface{}) (interface{}, error) {
	if cert == nil {
		return nil, fmt.Errorf("[ERROR] Certificate can't be nil")
	}

	switch t := cert.(type) {
	case *projectClient.NamespacedCertificate:
		return c.createNamespacedCertificate(cert.(*projectClient.NamespacedCertificate))
	case *projectClient.Certificate:
		return c.createCertificate(cert.(*projectClient.Certificate))
	default:
		return nil, fmt.Errorf("[ERROR] Certificate type %s isn't supported", t)
	}
}

func (c *Config) updateCertificate(cert *projectClient.Certificate, update interface{}) (*projectClient.Certificate, error) {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.Certificate.Update(cert, update)
}

func (c *Config) updateNamespacedCertificate(cert *projectClient.NamespacedCertificate, update interface{}) (*projectClient.NamespacedCertificate, error) {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return nil, err
	}
	return client.NamespacedCertificate.Update(cert, update)
}

func (c *Config) UpdateCertificate(cert interface{}, update interface{}) (interface{}, error) {
	if cert == nil {
		return nil, fmt.Errorf("[ERROR] Certificate can't be nil")
	}

	switch t := cert.(type) {
	case *projectClient.NamespacedCertificate:
		return c.updateNamespacedCertificate(cert.(*projectClient.NamespacedCertificate), update)
	case *projectClient.Certificate:
		return c.updateCertificate(cert.(*projectClient.Certificate), update)
	default:
		return nil, fmt.Errorf("[ERROR] Certificate type %s isn't supported", t)
	}
}

func (c *Config) deleteCertificate(cert *projectClient.Certificate) error {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return err
	}
	return client.Certificate.Delete(cert)
}

func (c *Config) deleteNamespacedCertificate(cert *projectClient.NamespacedCertificate) error {
	client, err := c.ProjectClient(cert.ProjectID)
	if err != nil {
		return err
	}
	return client.NamespacedCertificate.Delete(cert)
}

func (c *Config) DeleteCertificate(cert interface{}) error {
	if cert == nil {
		return fmt.Errorf("[ERROR] Certificate can't be nil")
	}

	switch t := cert.(type) {
	case *projectClient.NamespacedCertificate:
		return c.deleteNamespacedCertificate(cert.(*projectClient.NamespacedCertificate))
	case *projectClient.Certificate:
		return c.deleteCertificate(cert.(*projectClient.Certificate))
	default:
		return fmt.Errorf("[ERROR] Certificate type %s isn't supported", t)
	}
}

func (c *Config) GetRecipientByNotifier(id string) (*managementClient.Recipient, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("[ERROR] Notifier ID can't be nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	notifier, err := client.Notifier.ByID(id)
	if err != nil {
		return nil, err
	}

	out := &managementClient.Recipient{}

	out.NotifierID = notifier.ID
	if notifier.DingtalkConfig != nil {
		out.NotifierType = recipientTypeDingtalk
	} else if notifier.MSTeamsConfig != nil {
		out.NotifierType = recipientTypeMsTeams
	} else if notifier.PagerdutyConfig != nil {
		out.NotifierType = recipientTypePagerduty
		out.Recipient = notifier.PagerdutyConfig.ServiceKey
	} else if notifier.SlackConfig != nil {
		out.NotifierType = recipientTypeSlack
		out.Recipient = notifier.SlackConfig.DefaultRecipient
	} else if notifier.SMTPConfig != nil {
		out.NotifierType = recipientTypeSMTP
		out.Recipient = notifier.SMTPConfig.DefaultRecipient
	} else if notifier.WebhookConfig != nil {
		out.NotifierType = recipientTypeWebhook
		out.Recipient = notifier.WebhookConfig.URL
	} else if notifier.WechatConfig != nil {
		out.NotifierType = recipientTypeWechat
		out.Recipient = notifier.WechatConfig.DefaultRecipient
	}

	return out, nil
}
