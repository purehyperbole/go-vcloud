package vcloud

import (
	"fmt"
	"net/http"
)

// Config ...
type Config struct {
	URL      string
	Username string
	Password string
}

// Connector ...
type Connector struct {
	Config    *Config
	AuthToken string
}

// NewConnector ...
func NewConnector(config *Config) *Connector {
	connector := Connector{}
	connector.Config = config
	return &connector
}

// Authenticate ...
func (c *Connector) Authenticate() error {
	url := fmt.Sprintf("https://%s/api/sessions", c.Config.URL)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Config.Username, c.Config.Password)
	req.Header.Set("accept", "application/*+xml;version=5.1")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	c.AuthToken = resp.Header.Get("x-vcloud-authorization")
	return nil
}
