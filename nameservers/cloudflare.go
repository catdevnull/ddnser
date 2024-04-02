package nameservers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CloudflareV4 struct {
	HTTPClient *http.Client
	Key        string
	ZoneName   string
}

func (c *CloudflareV4) SetRecord(ctx context.Context, domain string, overrideIp string) (string, error) {
	zoneId, err := c.apiGetZoneId(ctx)
	if err != nil {
		return "", err
	}
	recordId, err := c.apiGetDnsRecordId(ctx, zoneId, domain)
	if err != nil {
		return "", err
	}

	content := overrideIp
	if len(content) == 0 {
		ip, err := c.getIp(ctx)
		if err != nil {
			return "", err
		}
		content = ip
	}

	result, err := c.apiUpdateDnsRecord(ctx, zoneId, domain, recordId, content)
	return result, err
}

func (c *CloudflareV4) apiReq(ctx context.Context, path string, method string, reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, "https://api.cloudflare.com"+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Key)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (c *CloudflareV4) apiGet(ctx context.Context, path string) ([]byte, error) {
	return c.apiReq(ctx, path, "GET", nil)
}

func (c *CloudflareV4) getIp(ctx context.Context) (string, error) {
	// icanhazip.com es manejado por Cloudflare https://blog.apnic.net/2021/06/17/how-a-small-free-ip-tool-survived/
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipv4.icanhazip.com/", nil)
	if err != nil {
		return "", err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type apiZoneSearchResponse struct {
	Result []struct {
		Id string `json:"id"`
	} `json:"result"`
}

func (c *CloudflareV4) apiGetZoneId(ctx context.Context) (string, error) {
	str, err := c.apiGet(ctx, "/client/v4/zones?name="+c.ZoneName)
	if err != nil {
		return "", err
	}
	var resp apiZoneSearchResponse
	err = json.Unmarshal(str, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Result) != 1 {
		return "", errors.New("expected result to be len()=1 but len is " + fmt.Sprintf("%d", len(resp.Result)))
	}
	return resp.Result[0].Id, nil
}

type apiDnsRecordSearchResponse struct {
	Result []struct {
		Id string `json:"id"`
	} `json:"result"`
}

func (c *CloudflareV4) apiGetDnsRecordId(ctx context.Context, zoneId, domain string) (string, error) {
	str, err := c.apiGet(ctx, fmt.Sprintf("/client/v4/zones/%s/dns_records?type=A&name=%s", zoneId, domain))
	if err != nil {
		return "", err
	}
	var resp apiDnsRecordSearchResponse
	err = json.Unmarshal(str, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Result) != 1 {
		return "", errors.New("expected result to be len()=1 but len is " + fmt.Sprintf("%d", len(resp.Result)))
	}
	return resp.Result[0].Id, nil
}

// https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-patch-dns-record
type apiUpdateDnsRecordRequest struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}
type apiUpdateDnsRecordResponse struct {
	Result []struct {
		Id      string `json:"id"`
		Content string `json:"content"`
	} `json:"result"`
}

func (c *CloudflareV4) apiUpdateDnsRecord(ctx context.Context, zoneId, zoneName, recordId, content string) (string, error) {
	body, err := json.Marshal(apiUpdateDnsRecordRequest{
		Content: content,
		Name:    zoneName,
		Type:    "A",
	})
	if err != nil {
		return "", err
	}

	str, err := c.apiReq(ctx,
		fmt.Sprintf("/client/v4/zones/%s/dns_records/%s", zoneId, recordId),
		"PATCH", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	log.Printf("debug: [cloudflareV4(%s)] Response: %s", zoneName, string(body))
	var resp apiUpdateDnsRecordResponse
	err = json.Unmarshal(str, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Result) != 1 {
		return "", errors.New("expected result to be len()=1 but len is " + fmt.Sprintf("%d", len(resp.Result)))
	}
	return resp.Result[0].Content, nil
}
