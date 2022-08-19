package goss

import (
	"fmt"
)

const gatewayBaseURL = "gateways"

type ProtoType string
type NATRuleType string

const (
	FirewallActionAllow  string = "Allow"
	FirewallActionDeny   string = "Deny"
	FirewallDirectionIn  string = "In"
	FirewallDirectionOut string = "Out"

	NATRuleTypeSNAT  NATRuleType = "SNAT"
	NATRuleTypeDNAT  NATRuleType = "DNAT"
	NATRuleTypeBINAT NATRuleType = "BINAT"

	ProtocolIP   ProtoType = "IP"
	ProtocolTCP  ProtoType = "TCP"
	ProtocolUDP  ProtoType = "UDP"
	ProtocolICMP ProtoType = "ICMP"
)

type (
	FirewallRule struct {
		Action          string    `json:"action,omitempty"`
		Direction       string    `json:"direction,omitempty"`
		Protocol        ProtoType `json:"protocol,omitempty"`
		Source          string    `json:"source,omitempty"`
		SourcePort      int       `json:"source_port,omitempty"`
		Destination     string    `json:"destination,omitempty"`
		DestinationPort int       `json:"destination_port,omitempty"`
	}

	NATRule struct {
		RuleType        NATRuleType `json:"type,omitempty"`
		Protocol        ProtoType   `json:"protocol,omitempty"`
		Source          string      `json:"source,omitempty"`
		Destination     string      `json:"destination,omitempty"`
		DestinationPort int         `json:"destination_port,omitempty"`
		Translated      string      `json:"translated,omitempty"`
		TranslatedPort  int         `json:"translated_port,omitempty"`
	}

	GatewayEntity struct {
		ID            string          `json:"id,omitempty"`
		LocationID    string          `json:"location_id,omitempty"`
		Name          string          `json:"name,omitempty"`
		NICS          []*NICEntity    `json:"nics,omitempty"`
		NetworkIDs    []string        `json:"network_ids,omitempty"`
		NATRules      []*NATRule      `json:"nat_rules,omitempty"`
		FirewallRules []*FirewallRule `json:"firewall_rules,omitempty"`
	}

	gatewayResponseWrap struct {
		Gateway *GatewayEntity `json:"gateway,omitempty"`
	}

	gatewayListResponseWrap struct {
		Gateways []*GatewayEntity `json:"gateways,omitempty"`
	}

	firewallRuleListResponseWrap struct {
		FirewallRules *[]FirewallRule `json:"firewall_rules,omitempty"`
	}

	NATRuleListResponseWrap struct {
		NATRules *[]NATRule `json:"nat_rules,omitempty"`
	}
)

func (c *SSClient) GetGateway(gatewayID string) (*GatewayEntity, error) {

	url := fmt.Sprintf("%s/%s", gatewayBaseURL, gatewayID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &gatewayResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*gatewayResponseWrap).Gateway, nil
}

func (c *SSClient) GetGatewayList() ([]*GatewayEntity, error) {

	resp, err := makeRequest(c.client, gatewayBaseURL, methodGet, nil, &gatewayListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*gatewayListResponseWrap).Gateways, nil
}

func (c *SSClient) CreateGateway(
	locationID string,
	name string,
	bandwidthMbps int,
	networkIDs []string,
) (*TaskIDWrap, error) {

	payload := map[string]interface{}{
		"location_id":    locationID,
		"name":           name,
		"bandwidth_mbps": bandwidthMbps,
		"network_ids":    networkIDs,
	}

	resp, err := makeRequest(c.client, gatewayBaseURL, methodPost, payload, &TaskIDWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateGatewayAndWait(
	locationID string,
	name string,
	bandwidthMbps int,
	networkIDs []string,
) (*GatewayEntity, error) {

	taskWrap, err := c.CreateGateway(locationID, name, bandwidthMbps, networkIDs)

	if err != nil {
		return nil, err
	}
	return c.waitGateway(taskWrap.ID)
}

func (c *SSClient) RenameGateway(gatewayID string, name string) error {

	url := fmt.Sprintf("%s/%s", gatewayBaseURL, gatewayID)
	payload := map[string]interface{}{
		"name": name,
	}

	_, err := makeRequest(c.client, url, methodPut, payload, nil)

	return err
}

func (c *SSClient) DeleteGateway(gatewayID string) error {

	url := fmt.Sprintf("%s/%s", gatewayBaseURL, gatewayID)

	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})

	return err
}

func (c *SSClient) EditGatewayBandwidth(gatewayID string, bandwidthMbps int) (*TaskIDWrap, error) {

	url := fmt.Sprintf("%s/%s/bandwidth", gatewayBaseURL, gatewayID)
	payload := map[string]interface{}{
		"bandwidth_mbps": bandwidthMbps,
	}

	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) GetFirewallRules(gatewayID string) (*[]FirewallRule, error) {

	url := fmt.Sprintf("%s/%s/firewall", gatewayBaseURL, gatewayID)

	resp, err := makeRequest(c.client, url, methodGet, nil, &firewallRuleListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*firewallRuleListResponseWrap).FirewallRules, nil
}

func (c *SSClient) EditFirewallRules(gatewayID string, firewallRules []FirewallRule) (*TaskIDWrap, error) {

	url := fmt.Sprintf("%s/%s/firewall", gatewayBaseURL, gatewayID)
	payload := map[string]interface{}{
		"firewall_rules": firewallRules,
	}

	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) EditFirewallRulesAndWait(gatewayID string, firewallRules []FirewallRule) (*GatewayEntity, error) {

	taskWrap, err := c.EditFirewallRules(gatewayID, firewallRules)

	if err != nil {
		return nil, err
	}
	return c.waitGateway(taskWrap.ID)
}

func (c *SSClient) GetNATRules(gatewayID string) (*[]NATRule, error) {

	url := fmt.Sprintf("%s/%s/nat", gatewayBaseURL, gatewayID)

	resp, err := makeRequest(c.client, url, methodGet, nil, &NATRuleListResponseWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*NATRuleListResponseWrap).NATRules, nil
}

func (c *SSClient) EditNATRules(gatewayID string, NATRules []NATRule) (*TaskIDWrap, error) {

	url := fmt.Sprintf("%s/%s/nat", gatewayBaseURL, gatewayID)
	payload := map[string]interface{}{
		"nat_rules": NATRules,
	}

	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})

	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) EditNATRulesAndWait(gatewayID string, NATRules []NATRule) (*GatewayEntity, error) {

	taskWrap, err := c.EditNATRules(gatewayID, NATRules)

	if err != nil {
		return nil, err
	}
	return c.waitGateway(taskWrap.ID)
}

func (c *SSClient) waitGateway(taskID string) (*GatewayEntity, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetGateway(task.GatewayID)
}
