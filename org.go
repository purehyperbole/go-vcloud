package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Org ...
type Org struct {
	Connector   *Connector
	XMLName     xml.Name `xml:"Org"`
	ID          string   `xml:"id,attr"`
	Name        string   `xml:"name,attr"`
	Href        string   `xml:"href,attr"`
	Links       []t.Link `xml:"Link"`
	Description string   `xml:"Description,value"`
	FullName    string   `xml:"FullName,value"`
}

// OrgList ...
func OrgList(c *Connector) *[]t.OrgListOrg {
	resp, err := c.Get("/api/org")
	if err != nil {
		log.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		log.Println(err)
	}

	orgList := ParseOrgList(data)

	return &orgList.Org
}

// NewOrg ...
func NewOrg(c *Connector, name string) *Org {
	href := findOrgHref(c, name)
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

	org := ParseOrg(data)
	org.Connector = c

	return org
}

// ParseOrgList ...
func ParseOrgList(d *[]byte) *t.OrgList {
	org := t.OrgList{}
	err := xml.Unmarshal(*d, &org)
	if err != nil {
		log.Println(err)
	}
	return &org
}

// ParseOrg ...
func ParseOrg(d *[]byte) *Org {
	org := Org{}
	err := xml.Unmarshal(*d, &org)
	if err != nil {
		log.Println(err)
	}
	return &org
}

// Datacenters ...
func (o *Org) Datacenters() []string {
	var datacenters []string
	return datacenters
}

// GetDatacenter ...
func (o *Org) GetDatacenter(name string) *Datacenter {
	datacenter := NewDatacenter(o.Connector, "")
	return datacenter
}

// Networks ...
func Networks() []string {
	var networks []string
	return networks
}

// GetNetwork ...
func GetNetwork(name string) {

}

func findOrgHref(c *Connector, name string) string {
	orgs := OrgList(c)
	for _, org := range *orgs {
		if org.Name == name {
			return org.Href
		}
	}
	return ""
}
