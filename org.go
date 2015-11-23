package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Org ...
type Org struct {
	Connector   *Connector `xml:"-"`
	XMLName     xml.Name   `xml:"Org"`
	ID          string     `xml:"id,attr"`
	Name        string     `xml:"name,attr"`
	Href        string     `xml:"href,attr"`
	Links       []t.Link   `xml:"Link"`
	Description string     `xml:"Description,value"`
	FullName    string     `xml:"FullName,value"`
}

// OrgList ...
func OrgList(c *Connector) (*[]t.OrgListOrg, error) {
	href := fmt.Sprintf("https://%s/api/org", c.Config.URL)
	resp, err := c.Get(href)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	orgList := parseOrgList(data)

	return &orgList.Org, nil
}

// NewOrg ...
func NewOrg(c *Connector, name string) (*Org, error) {
	href := findOrgHref(c, name)
	resp, err := c.Get(href)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	org := parseOrg(data)
	org.Connector = c

	return org, nil
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
func (o *Org) GetDatacenter(name string) (*Datacenter, error) {
	url := o.findLink("application/vnd.vmware.vcloud.vdc+xml", name)
	return NewDatacenter(o.Connector, url)
}

// Networks ...
func (o *Org) Networks() []t.Link {
	networks := o.findLinks("application/vnd.vmware.vcloud.orgNetwork+xml")
	return networks
}

// GetNetwork ...
func (o *Org) GetNetwork(name string) (*Network, error) {
	url := o.findLink("application/vnd.vmware.vcloud.vdc+xml", name)
	return NewNetwork(o.Connector, url)
}

// Catalogs ...
func (o *Org) Catalogs() []t.Link {
	catalogs := o.findLinks("application/vnd.vmware.vcloud.catalog+xml")
	return catalogs
}

func findOrgHref(c *Connector, name string) string {
	orgs, _ := OrgList(c)
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
