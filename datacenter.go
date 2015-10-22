package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Datacenter ...
type Datacenter struct {
	Connector         *Connector          `xml:"-"`
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
func (d *Datacenter) GetEdgeGateway(name string) (*EdgeGateway, error) {
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

// GetVApp ...
func (d *Datacenter) GetVApp(name string) *VApp {
	var href string
	for _, v := range d.VApps() {
		if v.Name == name {
			href = v.Href
		}
	}
	vapp := NewVApp(d.Connector, href)
	return vapp
}

// Networks ...
func (d *Datacenter) Networks() []t.Link {
	return d.AvailableNetworks.Networks
}

// GetNetwork ...
func (d *Datacenter) GetNetwork(name string) (*Network, error) {
	var href string
	for _, n := range d.AvailableNetworks.Networks {
		if n.Name == name {
			href = n.Href
		}
	}
	return NewNetwork(d.Connector, href)
}

// CreateNetwork ...
func (d *Datacenter) CreateNetwork(n *Network) (*Network, error) {
	links := d.findLinks(orgNetworkType)
	cnURL, err := url.Parse(links[0].Href)

	if err != nil {
		return nil, err
	}

	data, err := xml.Marshal(n)
	if err != nil {
		return nil, err
	}

	resp, err := d.Connector.Post(cnURL.RequestURI(), data, orgNetworkType)
	if err != nil {
		return nil, err
	}
	nw := Network{}
	nwdata, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(*nwdata, &nw)
	if err != nil {
		return nil, err
	}

	return &nw, nil
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
