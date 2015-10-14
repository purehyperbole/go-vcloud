package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Datacenter ...
type Datacenter struct {
	Connector         *Connector
	XMLName           xml.Name            `xml:"Vdc"`
	Name              string              `xml:"name,attr"`
	Href              string              `xml:"href,attr"`
	ComputeCapacity   t.ComputeCapacity   `xml:"ComputeCapacity"`
	AvailableNetworks t.AvailableNetworks `xml:"AvailableNetworks"`
	ResourceEntities  t.ResourceEntities  `xml:"ResourceEntities"`
	Links             []t.Link            `xml:"Link"`
}

// NewDatacenter ...
func NewDatacenter(c *Connector, href string) *Datacenter {
	dcURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(dcURL.RequestURI())
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

// GetEdgeGateway ...
func (d *Datacenter) GetEdgeGateway(name string) *EdgeGateway {
	return FindEdgeGateway(d.Connector, d.Href, name)
}

// VApps ...
func (d *Datacenter) VApps() []t.Link {
	var vapps []t.Link
	for _, e := range d.ResourceEntities.Entities {
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
func (d *Datacenter) Networks() []t.Link {
	return d.AvailableNetworks.Networks
}

// GetNetwork ...
func (d *Datacenter) GetNetwork(name string) *Network {
	var href string
	for _, n := range d.AvailableNetworks.Networks {
		if n.Name == name {
			href = n.Href
		}
	}
	network := NewNetwork(d.Connector, href)
	return network
}

// CreateNetwork ...
func (d *Datacenter) CreateNetwork(n *Network) (*Task, error) {
	task := Task{}

	links := d.findLinks("application/vnd.vmware.vcloud.orgVdcNetwork+xml")
	fmt.Println(links)
	cnURL, err := url.Parse(links[0].Href)

	if err != nil {
		log.Println(err)
	}

	data, err := xml.Marshal(n)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := d.Connector.Post(cnURL.RequestURI(), data, "application/vnd.vmware.vcloud.orgVdcNetwork+xml")
	if err != nil {
		fmt.Println(err)
	}
	tdata, err := ParseResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	err = xml.Unmarshal(*tdata, &task)
	if err != nil {
		fmt.Println(err)
	}

	return &task, nil
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
