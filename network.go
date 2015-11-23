package vcloud

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

const (
	orgNetworkType = "application/vnd.vmware.vcloud.orgVdcNetwork+xml"
)

// Network ...
type Network struct {
	Connector     *Connector `xml:"-"`
	XMLName       xml.Name   `xml:"http://www.vmware.com/vcloud/v1.5 OrgVdcNetwork"`
	XMLNS1        string     `xml:"xmlns:xsi,attr,omitempty"`
	XMLNS2        string     `xml:"xsi:schemaLocation,attr,omitempty"`
	Type          string     `xml:"type,attr,omitempty"`
	Name          string     `xml:"name,attr,omitempty"`
	Href          string     `xml:"href,attr,omitempty"`
	ID            string     `xml:"id,attr,omitempty"`
	Status        string     `xml:"status,attr,omitempty"`
	Description   string     `xml:"Description,value"`
	Tasks         *Tasks     `xml:"Tasks"`
	Configuration struct {
		IPScopes      t.IPScopes `xml:"IpScopes"`
		FenceMode     string     `xml:"FenceMode"`
		RetainNetInfo bool       `xml:"RetainNetInfoAcrossDeployments"`
	} `xml:"Configuration"`
	EdgeGateway *t.NetworkGateway `xml:"EdgeGateway,omitempty"`
	IsShared    bool              `xml:"IsShared,value,omitempty"`
}

// NewNetwork ...
func NewNetwork(c *Connector, href string) (*Network, error) {
	resp, err := c.Get(href)
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	n := parseNetwork(data)
	n.Connector = c

	return n, nil
}

func parseNetwork(d *[]byte) *Network {
	n := Network{}
	err := xml.Unmarshal(*d, &n)
	if err != nil {
		log.Println(err)
	}
	return &n
}

// Reload ...
func (n *Network) Reload() error {
	nw, err := NewNetwork(n.Connector, n.Href)
	n = nw
	return err
}

// Update ...
func (n *Network) Update() (*Task, error) {
	data, err := xml.Marshal(n)
	if err != nil {
		return nil, err
	}

	resp, err := n.Connector.Put(n.getAdminHref(), data, orgNetworkType)
	if err != nil {
		return nil, err
	}

	tdata, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	n.Reload()
	task := ParseTask(tdata)

	return task, nil
}

// Delete ...
func (n *Network) Delete() error {
	err := n.Connector.Delete(n.getAdminHref())
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// GetTasks ...
func (n *Network) GetTasks() []Task {
	for i := 0; i < len(n.Tasks.Task); i++ {
		n.Tasks.Task[i].Connector = n.Connector
	}
	return n.Tasks.Task
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
	n.EdgeGateway = &t.NetworkGateway{}
	n.EdgeGateway.Type = edgeGatewayType
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

func (n *Network) getAdminHref() string {
	return strings.Replace(n.Href, "/network", "/admin/network", 1)
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
