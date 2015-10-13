package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Datacenter ...
type Datacenter struct {
	Connector         *Connector
	XMLName           xml.Name             `xml:"Vdc"`
	Name              string               `xml:"name,attr"`
	ComputeCapacity   t.ComputeCapacity    `xml:"ComputeCapacity"`
	AvailableNetworks []t.AvailableNetwork `xml:"AvailableNetworks"`
	ResourceEntities  []t.ResourceEntity   `xml:"ResourceEntities"`
	Links             []t.Link             `xml:"Links"`
}

// NewDatacenter ...
func NewDatacenter(c *Connector, href string) *Datacenter {
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

	dc := parseDatacenter(data)
	dc.Connector = c

	return dc
}

func parseDatacenter(d *[]byte) *Datacenter {
	dc := Datacenter{}
	err := xml.Unmarshal(*d, &dc)
	if err != nil {
		log.Println(err)
	}
	return &dc
}

// VApps ...
func (d *Datacenter) VApps() []t.ResourceEntity {
	var vapps []t.ResourceEntity
	for _, e := range d.ResourceEntities {
		if e.Type == "application/vnd.vmware.vcloud.vApp+xml" {
			vapps = append(vapps, e)
		}
	}
	return vapps
}

// GetVapp ...
//func (d *Datacenter) GetVapp(name string) *VApp {

//}

// Networks ...
func (d *Datacenter) Networks() []t.AvailableNetwork {
	return d.AvailableNetworks
}

// GetNetwork ...
func (d *Datacenter) GetNetwork(name string) *Network {
	var href string
	for _, n := range d.AvailableNetworks {
		if n.Name == name {
			href = n.Href
		}
	}
	network := NewNetwork(d.Connector, href)
	return network
}

func (d *Datacenter) findLinks(xt string) []t.Link {
	var links []t.Link
	for _, link := range d.Links {
		if link.Type == xt {
			links = append(links, link)
		}
	}
	return links
}

func (d *Datacenter) findLink(xt string, name string) string {
	for _, link := range d.Links {
		if link.Type == xt && link.Name == name {
			return link.Href
		}
	}
	return ""
}
