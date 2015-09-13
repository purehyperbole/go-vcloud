package vcloud

import "encoding/xml"

// Datacenter ...
type Datacenter struct {
	Connector *Connector
	XMLName   xml.Name `xml:"VDC"`
}

// NewDatacenter ...
func NewDatacenter(c *Connector, url string) *Datacenter {
	datacenter := Datacenter{Connector: c}
	return &datacenter
}
