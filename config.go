package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"gitea.nulo.in/Nulo/ddnser/nameservers"
)

type config struct {
	Ip      string   `json:"ip,omitempty"`
	Every   int      `json:"every,omitempty"`
	Domains []domain `json:"domains"`
}

type domain struct {
	Type string `json:"type"`
	Name string `json:"name"`
	// TODO: lograr que esto sea un coso de propiedades arbitrario
	Key string `json:"key"`
}

func parseConfig(path string) (config config, err error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &config)
	return
}

type State struct {
	HTTPClient http.Client

	Ip string
	// Every defines how often (in seconds) poll for DDNS.
	// -1 means never poll.
	Every   int
	Domains []Domain
}
type Domain struct {
	Name       string
	NameServer nameservers.NameServer
}

func LoadConfig(path string) (state State, err error) {
	parsed, err := parseConfig(path)
	if err != nil {
		return
	}
	state.Ip = parsed.Ip
	state.Every = parsed.Every
	// if not defined or 0, set to default
	if state.Every == 0 {
		state.Every = 2
	}
	for _, d := range parsed.Domains {
		switch d.Type {
		case "njalla ddns":
			state.Domains = append(state.Domains, Domain{
				Name:       d.Name,
				NameServer: &nameservers.Njalla{HTTPClient: &state.HTTPClient, Key: d.Key},
			})
		case "he.net ddns":
			state.Domains = append(state.Domains, Domain{
				Name:       d.Name,
				NameServer: &nameservers.HeNet{HTTPClient: &state.HTTPClient, Password: d.Key},
			})
		default:
			err = errors.New("I don't know the service type " + d.Type)
			return
		}
	}
	return
}
