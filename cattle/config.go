package cattle

import (
	"fmt"
	//"log"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

const clusterProjectIDSeparator = ":"

// Client are the client kind for a Rancher API
type Client struct {
	Management *managementClient.Client
	Cluster    *clusterClient.Client
	Project    *projectClient.Client
}

// Config is the configuration parameters for a Rancher API
type Config struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	URL       string `json:"url"`
	CACerts   string `json:"cacert"`
	ClusterId string `json:"clusterId"`
	ProjectId string `json:"projectId"`
	Client    Client
}

// ManagementClient creates a Rancher client scoped to the global API
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

// ClusterClient creates a Rancher client scoped to an Cluster's API
func (c *Config) ClusterClient(id string) (*clusterClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Cluster Client: cluster ID is nil")
	}

	if c.Client.Cluster != nil && id == c.ClusterId {
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
	c.ClusterId = id

	return c.Client.Cluster, nil
}

// ProjectClient creates a Rancher client scoped to an Project's API
func (c *Config) ProjectClient(id string) (*projectClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Project Client: project ID is nil")
	}

	if c.Client.Project != nil && id == c.ProjectId {
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
	c.ProjectId = id

	return c.Client.Project, nil
}

func (c *Config) NormalizeUrl() {
	c.URL = strings.TrimSuffix(c.URL, "/")

	if !strings.HasSuffix(c.URL, "/v3") {
		c.URL = c.URL + "/v3"
	}
}

func (c *Config) CreateClientOpts() *clientbase.ClientOpts {
	c.NormalizeUrl()

	options := &clientbase.ClientOpts{
		URL:       c.URL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		TokenKey:  c.TokenKey,
		CACerts:   c.CACerts,
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

	projectsRoles, err := client.ProjectRoleTemplateBinding.List(listOpts)
	if err != nil {
		return nil, err
	}

	return projectsRoles.Data, nil
}

func (c *Config) GetProjectByName(name, clusterID string) (*managementClient.Project, error) {
	if name == "" {
		return nil, fmt.Errorf("[ERROR] Project name is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"ClusterId": clusterID}
	listOpts := NewListOpts(filters)

	projects, err := client.Project.List(listOpts)
	if err != nil {
		return nil, err
	}

	for _, project := range projects.Data {
		if project.Name == name {
			return &project, nil
		}
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

	client, err := c.ManagementClient()
	if err != nil {
		return "", err
	}

	project, err := client.Project.ByID(id)
	if err != nil {
		return "", err
	}

	return project.Name, nil
}

func NewListOpts(filters map[string]interface{}) *types.ListOpts {
	listOpts := clientbase.NewListOpts()
	if filters != nil {
		listOpts.Filters = filters
	}

	return listOpts
}

func IsNotFound(err error) bool {
	return clientbase.IsNotFound(err)
}

func splitID(id string) (clusterID, resourceID string) {
	separator := ":"
	if strings.Contains(id, separator) {
		return id[0:strings.Index(id, separator)], id[strings.Index(id, separator)+1:]
	}
	return "", id
}

func splitProjectID(id string) (clusterID, projectID string) {
	id = strings.TrimSuffix(id, clusterProjectIDSeparator)

	if strings.Contains(id, clusterProjectIDSeparator) {
		return id[0:strings.Index(id, clusterProjectIDSeparator)], id
	}

	return id, ""
}

/* RegistryClient creates a Rancher client scoped to a Registry's API
func (c *Config) RegistryClient(id string) (*rancherClient.RancherClient, error) {
	client, err := c.GlobalClient()
	if err != nil {
		return nil, err
	}
	reg, err := client.Registry.ById(id)
	if err != nil {
		return nil, err
	}

	return c.EnvironmentClient(reg.AccountId)
}

// CatalogClient creates a Rancher client scoped to a Catalog's API
func (c *Config) CatalogClient() (*catalog.RancherClient, error) {
	return catalog.NewRancherClient(&catalog.ClientOpts{
		Url:       c.APIURL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
}*/
