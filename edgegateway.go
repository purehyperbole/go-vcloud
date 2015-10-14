package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// EdgeGateway ...
type EdgeGateway struct {
	Connector     *Connector
	XMLName       xml.Name `xml:"EdgeGateway"`
	Name          string   `xml:"name,attr"`
	Href          string   `xml:"href,attr"`
	Configuration t.GatewayConfiguration
	Links         []t.Link `xml:"Link"`
}

// FindEdgeGateway ...
func FindEdgeGateway(c *Connector, dcHref string, name string) *EdgeGateway {
	q := Query{
		Connector: c,
		Type:      "edgeGateway",
		Format:    "records",
		Filter:    "vdc",
		FilterArg: dcHref,
	}

	resp := q.Run()
	data, err := ParseResponse(resp)
	if err != nil {
		fmt.Println(err)
	}

	results := t.QueryResultRecords{}
	err = xml.Unmarshal(*data, &results)
	if err != nil {
		fmt.Println(err)
	}

	gwHref := ""

	for _, gwr := range results.EdgeGatewayRecords {
		if gwr.Name == name {
			gwHref = gwr.Href
		}
	}

	gw := NewEdgeGateway(c, gwHref)

	return gw
}

// NewEdgeGateway ...
func NewEdgeGateway(c *Connector, href string) *EdgeGateway {
	gwURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(gwURL.RequestURI())
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
