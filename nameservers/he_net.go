package nameservers

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HeNet struct {
	HTTPClient *http.Client
	Password   string
}

func (h *HeNet) SetRecord(ctx context.Context, domain string, overrideIp string) (string, error) {
	u, _ := url.Parse("https://dyn.dns.he.net/nic/update")
	values := url.Values{
		"hostname": {domain},
		"password": {h.Password},
	}
	if len(overrideIp) > 0 {
		values["myip"] = []string{overrideIp}
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("debug: [he.net ddns(%s)] Response: %s", domain, string(body))

	things := strings.Split(string(body), " ")
	if len(things) != 2 {
		return "", errors.New("response is weird")
	}

	if things[0] != "good" && things[0] != "nochg" {
		return "", errors.New("Failed: " + things[0])
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("Not nice status code: " + strconv.Itoa(resp.StatusCode) + " with body: " + string(body))
	}

	return things[1], nil
}
