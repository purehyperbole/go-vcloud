package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

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
	Organization  t.Link     `xml:"Organization"`
}

// NewTask ...
func NewTask(c *Connector, href string) *Task {
	tURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(tURL.RequestURI())
	if err != nil {
		log.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		log.Println(err)
	}

	t := ParseTask(data)
	t.Connector = c

	return t
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
