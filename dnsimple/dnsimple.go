package dnsimple

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/libdns/libdns"
)

type Record struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Content   string    `json:"content,omitempty"`
	TTL       int       `json:"ttl,omitempty"`
	Type      string    `json:"type,omitempty"`
	Priority  int       `json:"priority,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (r Record) ToLibDns() libdns.Record {
	return libdns.Record{
		ID:       strconv.Itoa(r.ID),
		Type:     r.Type,
		Name:     r.Name,
		Value:    r.Content,
		TTL:      time.Duration(r.TTL),
		Priority: r.Priority,
	}
}

type Client struct {
	APIToken  string
	AccountID string
}

func NewClient(apiToken string, accountID string) Client {
	return Client{
		APIToken:  apiToken,
		AccountID: accountID,
	}
}

func (c Client) request(request http.Request, ctx context.Context, v interface{}) error {
	request.Header.Add("Authorization", "Bearer "+c.APIToken)

	resp, err := http.DefaultClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) GetRecords(zone string, ctx context.Context) ([]Record, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.dnsimple.com/v2/%s/zones/%s/records", c.AccountID, zone), nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []Record `json:"data"`
	}

	err = c.request(*request, ctx, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c Client) CreateRecord(zone string, record Record, ctx context.Context) (*Record, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("https://api.dnsimple.com/v2/%s/zones/%s/records", c.AccountID, zone), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")

	var resp struct {
		Data Record `json:"data"`
	}

	err = c.request(*request, ctx, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c Client) DeleteRecord(zone string, id int, ctx context.Context) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.dnsimple.com/v2/%s/zones/%s/records/%d", c.AccountID, zone, id), nil)
	if err != nil {
		return err
	}

	err = c.request(*request, ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
