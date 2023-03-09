package goss

import (
	"fmt"
	"log"
	"time"
)

const domainBaseURL = "domains"

type RecordType string

const (
	ARecordType     RecordType = "A"
	AAAARecordType  RecordType = "AAAA"
	MXRecordType    RecordType = "MX"
	CNAMERecordType RecordType = "CNAME"
	NSRecordType    RecordType = "NS"
	TXTRecordType   RecordType = "TXT"
	SRVRecordType   RecordType = "SRV"
)

type ProtocolType string

const (
	TCPProtocol ProtocolType = "TCP"
	UDPProtocol ProtocolType = "UDP"
	TLSProtocol ProtocolType = "TLS"
)

type (
	DomainRecord struct {
		ID             int          `json:"id"`
		Name           string       `json:"name"`
		Type           RecordType   `json:"type,omitempty"`
		IP             string       `json:"ip,omitempty"`
		MailHost       string       `json:"mail_host,omitempty"`
		Priority       int          `json:"priority,omitempty"`
		CanonicalName  string       `json:"canonical_name,omitempty"`
		NameServerHost string       `json:"name_server_host,omitempty"`
		Text           string       `json:"text,omitempty"`
		Protocol       ProtocolType `json:"protocol,omitempty"`
		Service        string       `json:"service,omitempty"`
		Weight         int          `json:"weight,omitempty"`
		Port           int          `json:"port,omitempty"`
		Target         string       `json:"target,omitempty"`
		TTL            string       `json:"ttl,omitempty"`
	}

	DomainResponse struct {
		Name        string          `json:"name,omitempty"`
		IsDelegated bool            `json:"is_delegated,omitempty"`
		Records     []*DomainRecord `json:"records,omitempty"`
	}

	domainResponseWrap struct {
		Domain *DomainResponse `json:"domain,omitempty"`
	}

	domainListResponseWrap struct {
		Domains []*DomainResponse `json:"domains,omitempty"`
	}
)

func (c *SSClient) GetDomain(domainName string) (*DomainResponse, error) {
	url := fmt.Sprintf("%s/%s", domainBaseURL, domainName)
	resp, err := makeRequest(c.client, url, methodGet, nil, &domainResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*domainResponseWrap).Domain, nil
}

func (c *SSClient) CreateDomain(
	name string,
	migrateRecords bool,
) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":            name,
		"migrate_records": migrateRecords,
	}

	resp, err := makeRequest(c.client, domainBaseURL, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateDomainAndWait(
	name string,
	migrateRecords bool,
) (*DomainResponse, error) {
	taskWrap, err := c.CreateDomain(name, migrateRecords)
	if err != nil {
		return nil, err
	}
	return c.waitDomain(taskWrap.ID)
}

func (c *SSClient) UpdateDomain(domainName string, cpu int, ram int) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"cpu":    cpu,
		"ram_mb": ram,
	}
	url := fmt.Sprintf("%s/%s", domainBaseURL, domainName)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) DeleteDomain(domainName string) error {
	url := fmt.Sprintf("%s/%s", domainBaseURL, domainName)
	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})
	return err
}

// -------- DOMAIN RECORDS --------

func (c *SSClient) CreateRecord(
	domainName string,
	record *DomainRecord,
) (*TaskIDWrap, error) {
	url := fmt.Sprintf("%s/%s/records", domainBaseURL, domainName)
	resp, err := makeRequest(c.client, url, methodPost, record, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateRecordAndWait(
	DomainName string,
	record *DomainRecord,
) (*DomainResponse, error) {
	taskWrap, err := c.CreateRecord(DomainName, record)
	if err != nil {
		return nil, err
	}
	return c.waitDomain(taskWrap.ID)
}

func (c *SSClient) UpdateRecord(
	domainName string,
	record *DomainRecord,
) (*TaskIDWrap, error) {
	url := fmt.Sprintf("%s/%s/records", domainBaseURL, domainName)
	resp, err := makeRequest(c.client, url, methodPut, record, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdateRecordAndWait(
	domainName string,
	record *DomainRecord,
) (*DomainResponse, error) {
	taskWrap, err := c.UpdateRecord(domainName, record)
	if err != nil {
		return nil, err
	}
	return c.waitDomain(taskWrap.ID)
}

func (c *SSClient) DeleteRecord(domainName string, recordId int) error {
	url := fmt.Sprintf("%s/%s/records/%d", domainBaseURL, domainName, recordId)
	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})
	if err != nil {
		return err
	}
	c.waitRecordDelition(domainName, recordId)
	return err
}

func (c *SSClient) waitDomain(taskID string) (*DomainResponse, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetDomain(task.DomainName)
}

func (c *SSClient) waitRecordDelition(domainName string, recordId int) (*DomainResponse, error) {
	const duration = defaultTaskCompletionDuration
	begin := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var (
		domain *DomainResponse
		err    error
	)
	for range ticker.C {
		recordWadDeleted := true

		domain, err = c.GetDomain(domainName)
		if err != nil {
			return nil, err
		}
		for _, record := range domain.Records {
			if record.ID == recordId {
				recordWadDeleted = false
				break
			}
		}
		if recordWadDeleted {
			return domain, nil
		} else {
			log.Default().Printf("[TRACE] Record isn't removed: %#v", domain)
		}

	}
	if time.Since(begin) > duration {
		return nil, fmt.Errorf("domain record wasn't removed for %f secs", duration.Seconds())
	}
	return domain, err
}


func (c *SSClient) GetDomainList() ([]*DomainResponse, error) {
	resp, err := makeRequest(c.client, domainBaseURL, methodGet, nil, &domainListResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*domainListResponseWrap).Domains, nil
}
