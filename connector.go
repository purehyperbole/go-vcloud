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
	Client    *http.Client
	AuthToken string
}

// NewConnector ...
func NewConnector(config *Config) *Connector {
	connector := Connector{}
	connector.Config = config
	connector.Client = &http.Client{}
	return &connector
}

// Authenticate ...
func (c *Connector) Authenticate() error {
	url := fmt.Sprintf("https://%s/api/sessions", c.Config.URL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Config.Username, c.Config.Password)
	req.Header.Set("accept", "application/*+xml;version=5.1")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	c.AuthToken = resp.Header.Get("x-vcloud-authorization")
	return nil
}

// Get ...
func (c *Connector) Get(uri string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s/api/%s", c.Config.URL, uri)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/*+xml;version=5.5")
	req.Header.Set("x-vcloud-authorization", c.AuthToken)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Post ...
func (c *Connector) Post() {

}

// Put ...
func (c *Connector) Put() {

}

// Delete ...
func (c *Connector) Delete() {

}
