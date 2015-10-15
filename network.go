package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Network ...
type Network struct {
	Connector     *Connector
	XMLName       xml.Name `xml:"http://www.vmware.com/vcloud/v1.5 OrgVdcNetwork"`
	Type          string   `xml:"type,attr,omitempty"`
	Name          string   `xml:"name,attr,omitempty"`
	Href          string   `xml:"href,attr,omitempty"`
	Status        string   `xml:"status,attr,omitempty"`
	Description   string   `xml:"Description,value"`
	Configuration struct {
		IPScopes      t.IPScopes `xml:"IpScopes"`
		FenceMode     string     `xml:"FenceMode"`
		RetainNetInfo bool       `xml:"RetainNetInfoAcrossDeployments"`
	} `xml:"Configuration"`
	EdgeGateway struct {
		Href string `xml:"href,attr"`
		Name string `xml:"name,attr"`
		Type string `xml:"type,attr"`
	} `xml:"EdgeGateway"`
	IsShared bool `xml:"IsShared,value,omitempty"`
}

// NewNetwork ...
func NewNetwork(c *Connector, href string) *Network {
	nURL, err := url.Parse(href)

	if href == "" && err != nil {
		log.Println(err)
	}

	resp, err := c.Get(nURL.RequestURI())
	if err != nil {
		log.Println(err)
	}

	data, err := ParseResponse(resp)
	if err != nil {
		log.Println(err)
	}

	n := parseNetwork(data)
	n.Connector = c

	return n
}

func parseNetwork(d *[]byte) *Network {
	n := Network{}
	err := xml.Unmarshal(*d, &n)
	if err != nil {
		log.Println(err)
	}
	return &n
}

// SetIsInherited ...
func (n *Network) SetIsInherited(inherited bool) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].IsInherited = inherited
}

// Netmask ...
func (n *Network) Netmask() string {
	n.configureIPScope()
	return n.Configuration.IPScopes.IPScope[0].Netmask
}

// SetNetmask ...
func (n *Network) SetNetmask(netmask string) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].Netmask = netmask
}

// Gateway ...
func (n *Network) Gateway() string {
	n.configureIPScope()
	return n.Configuration.IPScopes.IPScope[0].Gateway
}

// SetGateway ...
func (n *Network) SetGateway(gateway string) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].Gateway = gateway
}

// SetEdgeGateway ...
func (n *Network) SetEdgeGateway(href string, name string) {
	n.EdgeGateway.Type = "application/vnd.vmware.admin.edgeGateway+xml"
	n.EdgeGateway.Href = href
	n.EdgeGateway.Name = name
}

// SetIsEnabled ...
func (n *Network) SetIsEnabled(enabled bool) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].IsEnabled = enabled
}

// SetDNS1 ...
func (n *Network) SetDNS1(ns string) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].DNS1 = ns
}

// SetDNS2 ...
func (n *Network) SetDNS2(ns string) {
	n.configureIPScope()
	n.Configuration.IPScopes.IPScope[0].DNS2 = ns
}

// SetStartAddress ...
func (n *Network) SetStartAddress(start string) {
	n.configureIPRange()
	n.Configuration.IPScopes.IPScope[0].IPRanges.IPRange[0].StartAddress = start
}

// SetEndAddress ...
func (n *Network) SetEndAddress(end string) {
	n.configureIPRange()
	n.Configuration.IPScopes.IPScope[0].IPRanges.IPRange[0].EndAddress = end
}

// SetRetainNetInfo ...
func (n *Network) SetRetainNetInfo(retained bool) {
	n.configureIPRange()
	n.Configuration.RetainNetInfo = retained
}

// SetFenceMode ...
func (n *Network) SetFenceMode(mode string) {
	n.configureIPRange()
	n.Configuration.FenceMode = mode
}

// SetIsShared ...
func (n *Network) SetIsShared(shared bool) {
	n.IsShared = shared
}

func (n *Network) configureIPScope() {
	if len(n.Configuration.IPScopes.IPScope) < 1 {
		n.Configuration.IPScopes.IPScope = make([]t.IPScope, 1)
	}
}

func (n *Network) configureIPRange() {
	n.configureIPScope()
	if len(n.Configuration.IPScopes.IPScope[0].IPRanges.IPRange) < 1 {
		n.Configuration.IPScopes.IPScope[0].IPRanges.IPRange = make([]t.IPRange, 1)
	}
}
