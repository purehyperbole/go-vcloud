package vcloud

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

// Config ...
type Config struct {
	URL           string
	Username      string
	Password      string
	Debug         bool
	SSLSkipVerify bool
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.SSLSkipVerify},
	}
	connector.Config = config
	connector.Client = &http.Client{Transport: tr}
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

	if resp.StatusCode != 200 {
		return newError(resp)
	}

	c.AuthToken = resp.Header.Get("x-vcloud-authorization")
	return nil
}

// Get ...
func (c *Connector) Get(uri string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.URL, uri)

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

	if resp.StatusCode != 200 {
		return nil, newError(resp)
	}

	return resp, nil
}

// Post ...
func (c *Connector) Post(uri string, data []byte, contentType string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.URL, uri)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/*+xml;version=5.5")
	req.Header.Set("x-vcloud-authorization", c.AuthToken)
	req.Header.Set("Content-Type", contentType)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, newError(resp)
	}

	return resp, nil
}

// Put ...
func (c *Connector) Put(uri string, data []byte, contentType string) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", c.Config.URL, uri)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/*+xml;version=5.5")
	req.Header.Set("x-vcloud-authorization", c.AuthToken)
	req.Header.Set("Content-Type", contentType)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, newError(resp)
	}

	return resp, nil
}

// Delete ...
func (c *Connector) Delete(uri string) error {
	url := fmt.Sprintf("https://%s%s", c.Config.URL, uri)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/*+xml;version=5.5")
	req.Header.Set("x-vcloud-authorization", c.AuthToken)
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return newError(resp)
	}

	return nil
}

func newError(resp *http.Response) error {
	data, err := ParseResponse(resp)
	if err != nil {
		return err
	}
	vcloudErr := ParseError(data)
	return errors.New(vcloudErr.Message)
}
