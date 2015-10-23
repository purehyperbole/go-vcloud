package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

const (
	instantiateVAppTemplateParamsType = "application/vnd.vmware.vcloud.instantiateVAppTemplateParams+xml"
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
func NewDatacenter(c *Connector, href string) (*Datacenter, error) {
	dcURL, err := url.Parse(href)
	if err != nil {
		return nil, err
	}

	resp, err := c.Get(dcURL.RequestURI())
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	dc := parseDatacenter(data)
	dc.Connector = c

	return dc, nil
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
func (d *Datacenter) GetVApp(name string) (*VApp, error) {
	var href string
	for _, v := range d.VApps() {
		if v.Name == name {
			href = v.Href
		}
	}
	return NewVApp(d.Connector, href)
}

// CreateVApp ...
func (d *Datacenter) CreateVApp(request *t.InstantiateVApp) {
	links := d.findLinks(instantiateVAppTemplateParamsType)
	cvURL, err := url.Parse(links[0].Href)
	if err != nil {
		fmt.Println(cvURL)
	}

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

	nw.Connector = d.Connector

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
