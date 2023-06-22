// Package porkbun provides a client for the PorkBun API, primarily for managing Dynamic DNS
// records.
package porkbun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hay-kot/porkbun-dyndns-client/app/lib/httpclient"
)

type PorkBun struct {
	client           *httpclient.Client
	PorkBunKey       string
	PorkBunSecret    string
	PorkBunSubDomain string
	PorkBunDomain    string
}

func New(
	client *httpclient.Client,
	key string,
	secret string,
	subDomain string,
	domain string,
) *PorkBun {
	return &PorkBun{
		client:           client,
		PorkBunKey:       key,
		PorkBunSecret:    secret,
		PorkBunSubDomain: subDomain,
		PorkBunDomain:    domain,
	}
}

type credentials struct {
	APIKey    string `json:"apikey"`
	APISecret string `json:"secretapikey"`
}

type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	TTL     string `json:"ttl"`
}

var ErrNoRecordsFound = fmt.Errorf("no records found")

type DNSResponse struct {
	Status  string      `json:"status"`
	Records []DNSRecord `json:"records"`
}

func (pb *PorkBun) First(ctx context.Context) (*DNSRecord, error) {
	body := credentials{
		APIKey:    pb.PorkBunKey,
		APISecret: pb.PorkBunSecret,
	}

	bits, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	path := pb.client.Pathf("/dns/retrieveByNameType/%s/A/%s", pb.PorkBunDomain, pb.PorkBunSubDomain)
	resp, err := pb.client.Post(path, bits)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SetPorkBunDNSRecord: unexpected status code: %d", resp.StatusCode)
	}

	defer func() { _ = resp.Body.Close() }()

	var records DNSResponse

	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}

	if len(records.Records) == 0 {
		return nil, ErrNoRecordsFound
	}

	record := &DNSRecord{
		ID:      records.Records[0].ID,
		Name:    records.Records[0].Name,
		Content: records.Records[0].Content,
		Type:    records.Records[0].Type,
		TTL:     records.Records[0].TTL,
	}

	return record, nil
}

type DNSPost struct {
	DNSRecord
	credentials
}

func (pb *PorkBun) SetRecord(ctx context.Context, id string, record DNSRecord) error {
	body := DNSPost{
		DNSRecord: record,
		credentials: credentials{
			APIKey:    pb.PorkBunKey,
			APISecret: pb.PorkBunSecret,
		},
	}

	bits, err := json.Marshal(body)
	if err != nil {
		return err
	}

	path := pb.client.Pathf("/dns/edit/%s/%s", pb.PorkBunDomain, id)
	resp, err := pb.client.Post(path, bits)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SetPorkBunDNSRecord: unexpected status code: %d", resp.StatusCode)
	}

	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (pb *PorkBun) CreateRecord(ctx context.Context, record DNSRecord) error {
	type DNSPost struct {
		DNSRecord
		credentials
	}

	body := DNSPost{
		DNSRecord: record,
		credentials: credentials{
			APIKey:    pb.PorkBunKey,
			APISecret: pb.PorkBunSecret,
		},
	}

	bits, err := json.Marshal(body)
	if err != nil {
		return err
	}

	path := pb.client.Pathf("/dns/create/%s", pb.PorkBunDomain)
	resp, err := pb.client.Post(path, bits)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
