package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

const (
	edgeGatewayType = "application/vnd.vmware.admin.edgeGateway+xml"
)

// EdgeGateway ...
type EdgeGateway struct {
	Connector     *Connector `xml:"-"`
	XMLName       xml.Name   `xml:"EdgeGateway"`
	Name          string     `xml:"name,attr"`
	Href          string     `xml:"href,attr"`
	Configuration t.GatewayConfiguration
	Links         []t.Link `xml:"Link"`
}

// FindEdgeGateway ...
func FindEdgeGateway(c *Connector, dcHref string, name string) (*EdgeGateway, error) {
	q := Query{
		Connector: c,
		Type:      "edgeGateway",
		Format:    "records",
		Filter:    "vdc",
		FilterArg: dcHref,
	}

	resp, err := q.Run()
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	results := t.QueryResultRecords{}
	err = xml.Unmarshal(*data, &results)
	if err != nil {
		return nil, err
	}

	gwHref := ""

	for _, gwr := range results.EdgeGatewayRecords {
		if gwr.Name == name {
			gwHref = gwr.Href
		}
	}

	gw := NewEdgeGateway(c, gwHref)

	return gw, nil
}

// NewEdgeGateway ...
func NewEdgeGateway(c *Connector, href string) *EdgeGateway {
	resp, err := c.Get(href)
	if err != nil {
		fmt.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	gw := EdgeGateway{}
	err = xml.Unmarshal(*data, &gw)
	if err != nil {
		log.Println(err)
	}

	gw.Connector = c

	return &gw
}
