package cattle

import (
	"fmt"
	"log"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	clusterClient "github.com/rancher/types/client/cluster/v3"
	managementClient "github.com/rancher/types/client/management/v3"
	projectClient "github.com/rancher/types/client/project/v3"
)

const clusterProjectIDSeparator = ":"

// Config is the configuration parameters for a Rancher API
type Config struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	URL       string `json:"url"`
	CACerts   string `json:"cacert"`
}

// ManagementClient creates a Rancher client scoped to the global API
func (c *Config) ManagementClient() (*managementClient.Client, error) {
	options := c.CreateClientOpts()

	// Setup the management client
	mClient, err := managementClient.NewClient(options)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Management Client configured for url: %s", c.URL)

	return mClient, nil
}

// ClusterClient creates a Rancher client scoped to an Cluster's API
func (c *Config) ClusterClient(id string) (*clusterClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Cluster Client: cluster ID is nil")
	}

	options := c.CreateClientOpts()
	options.URL = options.URL + "/clusters/" + id

	// Setup the project client
	cClient, err := clusterClient.NewClient(options)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Cluster Client configured for url: %s", options.URL)

	return cClient, nil
}

// ProjectClient creates a Rancher client scoped to an Project's API
func (c *Config) ProjectClient(id string) (*projectClient.Client, error) {
	if id == "" {
		return nil, fmt.Errorf("[ERROR] Rancher Project Client: project ID is nil")
	}

	options := c.CreateClientOpts()
	options.URL = options.URL + "/projects/" + id

	// Setup the project client
	pClient, err := projectClient.NewClient(options)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Project Client configured for url: %s", options.URL)

	return pClient, nil
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

func (c *Config) GetProjectByName(name, clusterID string) (*managementClient.Project, error) {
	if name == "" {
		return nil, fmt.Errorf("[ERROR] Project name is nil")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, err
	}

	filters := map[string]interface{}{"ClusterId": clusterID}
	listOpts := NewListOpts()
	listOpts.Filters = filters

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

func NewListOpts() *types.ListOpts {
	return clientbase.NewListOpts()
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
