package rancher2

import (
	"fmt"
	"strings"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
)

const clusterProjectIDSeparator = ":"

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
