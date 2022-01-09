package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hack-as-a-service/caddy_dnsimple/dnsimple"
	"github.com/libdns/libdns"
)

// Provider facilitates DNS record manipulation with DNSimple.
type Provider struct {
	APIToken  string `json:"api_token,omitempty"`
	AccountID string `json:"account_id,omitempty"`
}

func (p *Provider) client() dnsimple.Client {
	return dnsimple.NewClient(p.APIToken, p.AccountID)
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	zone = normalizeZone(zone)
	client := p.client()

	records, err := client.GetRecords(zone, ctx)
	if err != nil {
		return nil, err
	}

	libdnsRecords := []libdns.Record{}

	for _, r := range records {
		libdnsRecords = append(libdnsRecords, r.ToLibDns())
	}

	return libdnsRecords, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	zone = normalizeZone(zone)
	client := p.client()

	addedRecords := []libdns.Record{}

	for _, r := range records {
		record, err := client.CreateRecord(zone, dnsimple.Record{
			Name:     r.Name,
			Type:     r.Type,
			Content:  r.Value,
			TTL:      int(r.TTL),
			Priority: r.Priority,
		}, ctx)
		if err != nil {
			return nil, err
		}

		addedRecords = append(addedRecords, record.ToLibDns())
	}

	return addedRecords, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return nil, fmt.Errorf("TODO: not implemented")
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	zone = normalizeZone(zone)
	client := p.client()

	deletedRecords := []libdns.Record{}

	for _, r := range records {
		id, err := strconv.Atoi(r.ID)
		if err != nil {
			return nil, err
		}

		err = client.DeleteRecord(zone, id, ctx)
		if err != nil {
			return nil, err
		}

		deletedRecords = append(deletedRecords, r)
	}
	return deletedRecords, nil
}

func normalizeZone(zone string) string {
	return strings.TrimSuffix(zone, ".")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
