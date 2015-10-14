package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Network ...
type Network struct {
	Connector     *Connector
	XMLName       xml.Name `xml:"OrgVdcNetwork"`
	Name          string   `xml:"name,attr"`
	Href          string   `xml:"href,attr"`
	Status        string   `xml:"status,attr"`
	Configuration struct {
		IPScopes t.IPScopes `xml:"IpScopes"`
	} `xml:"Configuration"`
	EdgeGateway struct {
		Href string `xml:"href,attr"`
	} `xml:"EdgeGateway"`
	IsShared string `xml:"IsShared,value"`
}

// NewNetwork ...
func NewNetwork(c *Connector, href string) *Network {
	nURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(nURL.RequestURI())
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

// Netmask ...
func (n *Network) Netmask() string {
	return n.Configuration.IPScopes.IPScope[0].Netmask
}

// Gateway ...
func (n *Network) Gateway() string {
	return n.Configuration.IPScopes.IPScope[0].Gateway
}
