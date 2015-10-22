package vcloud

import (
	"encoding/xml"
	"errors"
	"log"
	"net/url"
	"time"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Task ...
type Task struct {
	Connector     *Connector `xml:"-"`
	XMLName       xml.Name   `xml:"Task"`
	Name          string     `xml:"name,attr"`
	Href          string     `xml:"href,attr"`
	OperationName string     `xml:"operationName,attr"`
	Status        string     `xml:"status,attr"`
	StartTime     string     `xml:"startTime,attr"`
	ExpiryTime    string     `xml:"expiryTime,attr"`
	Cancel        bool       `xml:"cancelRequested"`
	Links         []t.Link   `xml:"Link"`
	Owner         t.Link     `xml:"Owner"`
	User          t.Link     `xml:"User"`
	Error         *t.Error   `xml:"Error"`
	Organization  t.Link     `xml:"Organization"`
}

// NewTask ...
func NewTask(c *Connector, href string) (*Task, error) {
	tURL, err := url.Parse(href)

	if err != nil {
		return nil, err
	}

	resp, err := c.Get(tURL.RequestURI())
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	t := ParseTask(data)
	t.Connector = c

	return t, nil
}

// ParseTask ...
func ParseTask(d *[]byte) *Task {
	t := Task{}
	err := xml.Unmarshal(*d, &t)
	if err != nil {
		log.Println(err)
	}
	return &t
}

// Wait ...
func (t *Task) Wait() error {
	for {
		if t.Status == "running" {
			t, _ = NewTask(t.Connector, t.Href)
		} else {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if t.Status == "error" {
		return errors.New(t.Error.Message)
	}
	return nil
}
