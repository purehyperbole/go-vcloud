package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"
)

// Network ...
type Network struct {
	Connector *Connector
	XMLName   xml.Name `xml:"OrgNetwork"`
}

// NewNetwork ...
func NewNetwork(c *Connector, href string) *Network {
	url, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(url.RequestURI())
	if err != nil {
		log.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		log.Println(err)
	}

	n := parseNetwork(data)
	n.Connector = c

	return n
}

func parseNetwork(d *[]byte) *Network {
	n := Network{}
	err := xml.Unmarshal(*d, &n)
	if err != nil {
		log.Println(err)
	}
	return &n
}
