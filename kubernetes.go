package goss

import (
	"fmt"
)

const kubernetesBaseURL = "k8s_clusters"

type (
	KubernetesClusterEntity struct {
		ID               string                       `json:"id,omitempty"`
		LocationID       string                       `json:"location_id,omitempty"`
		Name             string                       `json:"name,omitempty"`
		Version          string                       `json:"version,omitempty"`
		HighAvailability bool                         `json:"high_availability,omitempty"`
		NodeGroups       []*KubernetesNodeGroupEntity `json:"node_groups,omitempty"`
		State            string                       `json:"state,omitempty"`
		Tags             []string                     `json:"tags,omitempty"`
	}

	KubernetesNodeGroupEntity struct {
		ID            string   `json:"id,omitempty"`
		Name          string   `json:"name,omitempty"`
		CPUPerNode    int      `json:"cpu_per_node,omitempty"`
		RAMPerNode    int      `json:"ram_per_node,omitempty"`
		NumberOfNodes int      `json:"number_of_nodes,omitempty"`
		Nodes         []string `json:"nodes,omitempty"`
		Ingress       bool     `json:"ingress,omitempty"`
		State         string   `json:"state,omitempty"`
		Tags          []string `json:"tags,omitempty"`
	}

	kubernetesClusterResponseWrap struct {
		KubernetesCluster *KubernetesClusterEntity `json:"kubernetes_cluster,omitempty"`
	}

	kubernetesClusterListResponseWrap struct {
		KubernetesClusters []*KubernetesClusterEntity `json:"kubernetes_clusters,omitempty"`
	}

	kubernetesNodeGroupResponseWrap struct {
		NodeGroup *KubernetesNodeGroupEntity `json:"node_group,omitempty"`
	}

	kubernetesNodeGroupListResponseWrap struct {
		NodeGroups []*KubernetesNodeGroupEntity `json:"node_groups,omitempty"`
	}
	kubernetesVersionsResponseWrap struct {
		Versions []string `json:"versions,omitempty"`
	}
)

func (c *SSClient) GetKubernetesVersions() ([]string, error) {
	const kubernetesVersionURL string = "k8s_versions"
	resp, err := makeRequest(c.client, kubernetesVersionURL, methodGet, nil, &kubernetesVersionsResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesVersionsResponseWrap).Versions, nil
}

func (c *SSClient) GetAvailableKubernetesVersions(kubernetesClusterID string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/k8s_versions", kubernetesBaseURL, kubernetesClusterID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &kubernetesVersionsResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesVersionsResponseWrap).Versions, nil
}

func (c *SSClient) GetKubernetesCluster(kubernetesClusterID string) (*KubernetesClusterEntity, error) {
	url := fmt.Sprintf("%s/%s", kubernetesBaseURL, kubernetesClusterID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &kubernetesClusterResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) GetKubernetesClusterList() ([]*KubernetesClusterEntity, error) {
	resp, err := makeRequest(c.client, kubernetesBaseURL, methodGet, nil, &kubernetesClusterListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterListResponseWrap).KubernetesClusters, nil
}

func (c *SSClient) GetKubernetesNodeGroupList() ([]*KubernetesNodeGroupEntity, error) {
	resp, err := makeRequest(c.client, kubernetesBaseURL, methodGet, nil, &kubernetesNodeGroupListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesNodeGroupListResponseWrap).NodeGroups, nil
}

func (c *SSClient) GetKubernetesNodeGroup(kubernetesClusterID string, nodeGroupID string) (*KubernetesNodeGroupEntity, error) {
	url := fmt.Sprintf("%s/%s/node_groups/%s", kubernetesBaseURL, kubernetesClusterID, nodeGroupID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &kubernetesNodeGroupListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesNodeGroupResponseWrap).NodeGroup, nil
}

func (c *SSClient) DeleteKubernetesCluster(kubernetesClusterID string) error {

	url := fmt.Sprintf("%s/%s", kubernetesBaseURL, kubernetesClusterID)

	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})

	return err
}

func (c *SSClient) DeleteKubernetesNodeGroup(kubernetesClusterID string, nodeGroupID string) error {

	url := fmt.Sprintf("%s/%s/node_groups/%s", kubernetesBaseURL, kubernetesClusterID, nodeGroupID)

	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})

	return err
}

func (c *SSClient) CreateKubernetesCluster(
	locationID string,
	name string,
	version string,
	highAvailability bool,
	tags []string,
	nodeGroups []*KubernetesNodeGroupEntity,
) (*KubernetesClusterEntity, error) {

	payload := map[string]interface{}{
		"location_id":       locationID,
		"name":              name,
		"version":           name,
		"high_availability": highAvailability,
		"tags":              tags,
		"node_groups":       nodeGroups,
	}

	resp, err := makeRequest(c.client, gatewayBaseURL, methodPost, payload, &kubernetesClusterResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) CreateKubernetesClusterAndWait(
	locationID string,
	name string,
	version string,
	highAvailability bool,
	tags []string,
	nodeGroups []*KubernetesNodeGroupEntity,
) (*KubernetesClusterEntity, error) {

	taskWrap, err := c.CreateKubernetesCluster(locationID, name, version, highAvailability, tags, nodeGroups)

	if err != nil {
		return nil, err
	}
	return c.waitKubernetesCluster(taskWrap.ID)
}

func (c *SSClient) CreateKubernetesNodeGroups(kubernetesClusterID string, nodeGroups []*KubernetesNodeGroupEntity) (*KubernetesClusterEntity, error) {

	url := fmt.Sprintf("%s/%s/node_groups", kubernetesBaseURL, kubernetesClusterID)
	payload := map[string]interface{}{
		"node_groups": nodeGroups,
	}

	resp, err := makeRequest(c.client, url, methodPost, payload, &kubernetesClusterResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) CreateKubernetesNodeGroupsAndWait(kubernetesClusterID string, nodeGroups []*KubernetesNodeGroupEntity) (*KubernetesClusterEntity, error) {

	taskWrap, err := c.CreateKubernetesNodeGroups(kubernetesClusterID, nodeGroups)

	if err != nil {
		return nil, err
	}
	return c.waitKubernetesCluster(taskWrap.ID)
}

func (c *SSClient) ScaleKubernetesNodeGroup(kubernetesClusterID string, nodeGroupID string, nodeReplicas int) (*KubernetesClusterEntity, error) {

	url := fmt.Sprintf("%s/%s/node_groups/%s", kubernetesBaseURL, kubernetesClusterID, nodeGroupID)
	payload := map[string]interface{}{
		"number_of_nodes": nodeReplicas,
	}
	resp, err := makeRequest(c.client, url, methodPut, payload, &kubernetesClusterResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) ScaleKubernetesNodeGroupAndWait(kubernetesClusterID string, nodeGroupID string, nodeReplicas int) (*KubernetesClusterEntity, error) {

	taskWrap, err := c.ScaleKubernetesNodeGroup(kubernetesClusterID, nodeGroupID, nodeReplicas)

	if err != nil {
		return nil, err
	}
	return c.waitKubernetesCluster(taskWrap.ID)
}

func (c *SSClient) DeployIngressController(kubernetesClusterID string, nodeGroupID string) (*KubernetesClusterEntity, error) {

	url := fmt.Sprintf("%s/%s/node_groups/%s/ingress", kubernetesBaseURL, kubernetesClusterID, nodeGroupID)

	resp, err := makeRequest(c.client, url, methodPost, nil, &kubernetesClusterResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) DeployIngressControllerAndWait(kubernetesClusterID string, nodeGroupID string) (*KubernetesClusterEntity, error) {

	taskWrap, err := c.DeployIngressController(kubernetesClusterID, nodeGroupID)

	if err != nil {
		return nil, err
	}
	return c.GetKubernetesCluster(taskWrap.ID)
}

func (c *SSClient) UpgradeKubernetesCluster(kubernetesClusterID string, version string) (*KubernetesClusterEntity, error) {
	url := fmt.Sprintf("%s/%s", kubernetesBaseURL, kubernetesClusterID)
	payload := map[string]interface{}{
		"version": version,
	}

	resp, err := makeRequest(c.client, url, methodPut, payload, &kubernetesClusterResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*kubernetesClusterResponseWrap).KubernetesCluster, nil
}

func (c *SSClient) waitKubernetesCluster(taskID string) (*KubernetesClusterEntity, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetKubernetesCluster(task.KubernetesClusterID)
}
