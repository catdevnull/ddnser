package nameservers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Njalla struct {
	HTTPClient *http.Client
	Key        string
}

type njallaValue struct {
	A string `json:"A"`
}

type njallaResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Value   njallaValue `json:"value"`
}

func (n *Njalla) SetRecord(ctx context.Context, domain string, overrideIp string) (string, error) {
	u, _ := url.Parse("https://njal.la/update/")
	values := url.Values{
		"h": {domain},
		"k": {n.Key},
	}
	if len(overrideIp) > 0 {

		values["a"] = []string{overrideIp}
	} else {
		values["auto"] = []string{""}
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := n.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	decoded := njallaResponse{}
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		return "", err
	}
	log.Printf("debug: [njalla ddns] Response: %s", string(body))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("Not nice status code: " + strconv.Itoa(resp.StatusCode) + " with body: " + string(body))
	}
	if decoded.Status < 200 || decoded.Status > 299 {
		return "", errors.New("Not nice JSON status: " + strconv.Itoa(decoded.Status) + " with body: " + string(body))
	}
	return decoded.Value.A, nil
}
