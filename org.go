package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Org ...
type Org struct {
	Connector *Connector
	XMLName   xml.Name `xml:"Org"`
	Name      string   `xml:"name,attr"`
	Href      string   `xml:"href,attr"`
}

// OrgList ...
func OrgList(c *Connector) *[]t.Org {
	orgList := t.OrgList{}

	resp, err := c.Get("org")
	if err != nil {
		log.Println(err)
	}

	err = ParseResponse(resp, orgList)
	if err != nil {
		log.Println(err)
	}

	return &orgList.Org
}

// NewOrg ...
func NewOrg(c *Connector, name string) *Org {
	org := Org{Connector: c}
	href := findOrgHref(c, name)
	url, err := url.Parse(href)

	if href != "" && err == nil {
		resp, err := c.Get(url.RequestURI())
		if err != nil {
			log.Println(err)
		}
		err = ParseResponse(resp, org)
		if err != nil {
			log.Println(err)
		}
	}

	return &org
}

// Datacenters ...
func (o *Org) Datacenters() *[]string {
	var datacenters []string
	return &datacenters
}

// GetDatacenter ...
func (o *Org) GetDatacenter(name string) *Datacenter {
	datacenter := NewDatacenter(o.Connector, "")
	return datacenter
}

func findOrgHref(c *Connector, name string) string {
	orgs := OrgList(c)
	for _, org := range *orgs {
		log.Println(org)
		if org.Name == name {
			return org.Href
		}
	}
	return ""
}
