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

	orgList := parseOrgList(data)

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

	org := parseOrg(data)
	org.Connector = c

	return org
}

func parseOrgList(d *[]byte) *t.OrgList {
	org := t.OrgList{}
	err := xml.Unmarshal(*d, &org)
	if err != nil {
		log.Println(err)
	}
	return &org
}

func parseOrg(d *[]byte) *Org {
	org := Org{}
	err := xml.Unmarshal(*d, &org)
	if err != nil {
		log.Println(err)
	}
	return &org
}

// Datacenters ...
func (o *Org) Datacenters() []t.Link {
	datacenters := o.findLinks("application/vnd.vmware.vcloud.vdc+xml")
	return datacenters
}

// GetDatacenter ...
func (o *Org) GetDatacenter(name string) *Datacenter {
	url := o.findLink("application/vnd.vmware.vcloud.vdc+xml", name)
	datacenter := NewDatacenter(o.Connector, url)
	return datacenter
}

// Networks ...
func (o *Org) Networks() []t.Link {
	networks := o.findLinks("application/vnd.vmware.vcloud.orgNetwork+xml")
	return networks
}

// GetNetwork ...
func (o *Org) GetNetwork(name string) *Network {
	url := o.findLink("application/vnd.vmware.vcloud.vdc+xml", name)
	network := NewNetwork(o.Connector, url)
	return network
}

// Catalogs ...
func (o *Org) Catalogs() []t.Link {
	catalogs := o.findLinks("application/vnd.vmware.vcloud.catalog+xml")
	return catalogs
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

func (o *Org) findLinks(xt string) []t.Link {
	var links []t.Link
	for _, link := range o.Links {
		if link.Type == xt {
			links = append(links, link)
		}
	}
	return links
}

func (o *Org) findLink(xt string, name string) string {
	for _, link := range o.Links {
		if link.Type == xt && link.Name == name {
			return link.Href
		}
	}
	return ""
}
