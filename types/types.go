package types

import "encoding/xml"

// Link ...
type Link struct {
	XMLName xml.Name `xml:"Link"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Href    string   `xml:"href,attr"`
}

// OrgList ...
type OrgList struct {
	XMLName xml.Name     `xml:"OrgList"`
	Org     []OrgListOrg `xml:"Org"`
}

// Error ...
type Error struct {
	XMLName        xml.Name `xml:"Error"`
	MinorErrorCode string   `xml:"minorErrorCode,attr"`
	MajorErrorCode string   `xml:"majorErrorCode,attr"`
	Message        string   `xml:"message,attr"`
}

// OrgListOrg ...
type OrgListOrg struct {
	XMLName xml.Name `xml:"Org"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Href    string   `xml:"href,attr"`
}

// IPRange ...
type IPRange struct {
	XMLName      xml.Name `xml:"IpRange"`
	StartAddress string   `xml:"StartAddress,value"`
	EndAddress   string   `xml:"EndAddress,value"`
}

// IPRanges ...
type IPRanges struct {
	XMLName xml.Name  `xml:"IpRanges"`
	IPRange []IPRange `xml:"IpRange"`
}

// EdgeGateway ...
type EdgeGateway struct {
	XMLName xml.Name `xml:"EdgeGateway"`
	Name    string   `xml:"name,attr"`
}

// SubAllocation ...
type SubAllocation struct {
	XMLName     xml.Name    `xml:"SubAllocation"`
	EdgeGateway EdgeGateway `xml:"EdgeGateway"`
	IPRanges    IPRanges    `xml:"IpRanges"`
}

// SubAllocations ...
type SubAllocations struct {
	XMLName       xml.Name        `xml:"SubAllocations"`
	SubAllocation []SubAllocation `xml:"SubAllocation"`
}

// AllocatedIPAddresses ...
type AllocatedIPAddresses struct {
	XMLName   xml.Name `xml:"AllocatedIpAddresses"`
	IPAddress []string `xml:"IpAddress,value"`
}

// IPScope ...
type IPScope struct {
	XMLName              xml.Name             `xml:"IpScope"`
	SubAllocations       SubAllocations       `xml:"SubAllocations"`
	IPRanges             IPRanges             `xml:"IpRanges"`
	AllocatedIPAddresses AllocatedIPAddresses `xml:"AllocatedIpAddresses"`
}

// IPScopes ...
type IPScopes struct {
	XMLName xml.Name  `xml:"IpScopes"`
	IPScope []IPScope `xml:"IpScope"`
}

// Configuration ...
type Configuration struct {
	XMLName  xml.Name `xml:"Configuration"`
	IPScopes IPScopes `xml:"IpScopes"`
}

// ExternalNetwork ...
type ExternalNetwork struct {
	XMLName       xml.Name      `xml:"ExternalNetwork"`
	Configuration Configuration `xml:"Configuration"`
}

// ComputeCapacity ...
type ComputeCapacity struct {
	XMLName xml.Name `xml:"ComputeCapacity"`
	CPU     struct {
		Units     string `xml:"Units,value"`
		Allocated int    `xml:"Allocated,value"`
		Limit     int    `xml:"Limit,value"`
		Reserved  int    `xml:"Reserved,value"`
		Used      int    `xml:"Used,value"`
		Overhead  int    `xml:"Overhead,value"`
	} `xml:"Cpu"`
	Memory struct {
		Units     string `xml:"Units,value"`
		Allocated int    `xml:"Allocated,value"`
		Limit     int    `xml:"Limit,value"`
		Reserved  int    `xml:"Reserved,value"`
		Used      int    `xml:"Used,value"`
		Overhead  int    `xml:"Overhead,value"`
	} `xml:"memory"`
}

// AvailableNetwork ...
type AvailableNetwork struct {
	XMLName xml.Name `xml:"Network"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Href    string   `xml:"href,attr"`
}

// ResourceEntity ...
type ResourceEntity struct {
	XMLName xml.Name `xml:"ResourceEntity"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name,attr"`
	Href    string   `xml:"href,attr"`
}
