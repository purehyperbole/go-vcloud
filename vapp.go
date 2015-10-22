package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// VApp ...
type VApp struct {
	Connector *Connector `xml:"-"`
	XMLName   xml.Name   `xml:"VApp"`
	Name      string     `xml:"name,attr"`
	Href      string     `xml:"href,attr"`
	Status    string     `xml:"status,attr"`
	Deployed  bool       `xml:"deployed,attr"`
	Links     []t.Link   `xml:"Link"`
	//NetworkConfig t.
}

// NewVApp ...
func NewVApp(c *Connector, href string) *VApp {
	vURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(vURL.RequestURI())
	if err != nil {
		log.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		log.Println(err)
	}

	v := parseVApp(data)
	v.Connector = c

	return v
}

func parseVApp(d *[]byte) *VApp {
	v := VApp{}
	err := xml.Unmarshal(*d, &v)
	if err != nil {
		log.Println(err)
	}
	return &v
}
