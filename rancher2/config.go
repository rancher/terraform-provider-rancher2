package rancher2

import (
	"fmt"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

// Client are the client kind for a Rancher v3 API
type Client struct {
	Management *managementClient.Client
	Cluster    *clusterClient.Client
	Project    *projectClient.Client
}

// Config is the configuration parameters for a Rancher v3 API
type Config struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	URL       string `json:"url"`
	CACerts   string `json:"cacert"`
	Insecure  bool   `json:"insecure"`
	ClusterID string `json:"clusterId"`
	ProjectID string `json:"projectId"`
	Client    Client
}

// ManagementClient creates a Rancher client scoped to the management API
func (c *Config) ManagementClient() (*managementClient.Client, error) {
	if c.Client.Management != nil {
		return c.Client.Management, nil
	}

	options := c.CreateClientOpts()

	// Setup the management client
	mClient, err := managementClient.NewClient(options)
	if err != nil {
		return nil, err
	}
	c.Client.Management = mClient

	return c.Client.Management, nil
}

// ClusterClient creates a Rancher client scoped to a Cluster API
func (c *Config) ClusterClient(id string) (*clusterClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Cluster Client: cluster ID is nil")
	}

	if c.Client.Cluster != nil && id == c.ClusterID {
		return c.Client.Cluster, nil
	}

	options := c.CreateClientOpts()
	options.URL = options.URL + "/clusters/" + id

	// Setup the project client
	cClient, err := clusterClient.NewClient(options)
	if err != nil {
		return nil, err
	}
	c.Client.Cluster = cClient
	c.ClusterID = id

	return c.Client.Cluster, nil
}

// ProjectClient creates a Rancher client scoped to a Project API
func (c *Config) ProjectClient(id string) (*projectClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Project Client: project ID is nil")
	}

	if c.Client.Project != nil && id == c.ProjectID {
		return c.Client.Project, nil
	}

	options := c.CreateClientOpts()
	options.URL = options.URL + "/projects/" + id

	// Setup the project client
	pClient, err := projectClient.NewClient(options)
	if err != nil {
		return nil, err
	}

	c.Client.Project = pClient
	c.ProjectID = id

	return c.Client.Project, nil
}

func (c *Config) NormalizeURL() {
	c.URL = NormalizeURL(c.URL)
}

func (c *Config) CreateClientOpts() *clientbase.ClientOpts {
	c.NormalizeURL()

	options := &clientbase.ClientOpts{
		URL:       c.URL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		TokenKey:  c.TokenKey,
		CACerts:   c.CACerts,
		Insecure:  c.Insecure,
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

func (c *Config) isClusterActive(id string) (bool, error) {
	clus, err := c.GetClusterByID(id)
	if err != nil {
		return false, err
	}

	if clus.State == "active" {
		return true, nil
	}

	return false, nil
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

func (c *Config) activateNodeDriver(id string) error {
	if id == "" {
		return fmt.Errorf("[ERROR] Node Driver id is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return err
	}

	driver, err := client.NodeDriver.ByID(id)
	if err != nil {
		return fmt.Errorf("[ERROR] Getting Node Driver %s: %v", id, err)
	}

	if driver.State == "active" {
		return nil
	}

	_, err = client.NodeDriver.ActionActivate(driver)
	if err != nil {
		return fmt.Errorf("[ERROR] Activating Node Driver %s: %v", id, err)
	}

	return nil
}
