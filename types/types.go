package types

import "encoding/xml"

// Link ...
type Link struct {
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
	Href string `xml:"href,attr"`
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
	XMLName              xml.Name              `xml:"IpScope"`
	IsInherited          bool                  `xml:"IsInherited,value"`
	Gateway              string                `xml:"Gateway,value"`
	Netmask              string                `xml:"Netmask,value"`
	DNS1                 string                `xml:"Dns1,value,omitempty"`
	DNS2                 string                `xml:"Dns2,value,omitempty"`
	DNSSuffix            string                `xml:"DnsSuffix,value,omitempty"`
	IsEnabled            bool                  `xml:"IsEnabled,value,omitempty"`
	IPRanges             IPRanges              `xml:"IpRanges"`
	SubAllocations       *SubAllocations       `xml:"SubAllocations"`
	AllocatedIPAddresses *AllocatedIPAddresses `xml:"AllocatedIpAddresses"`
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

// AvailableNetworks ...
type AvailableNetworks struct {
	XMLName  xml.Name `xml:"AvailableNetworks"`
	Networks []Link   `xml:"Network"`
}

// ResourceEntities ...
type ResourceEntities struct {
	XMLName  xml.Name `xml:"ResourceEntities"`
	Entities []Link   `xml:"ResourceEntity"`
}

// QueryResultRecords ...
type QueryResultRecords struct {
	XMLName            xml.Name `xml:"QueryResultRecords"`
	Total              int      `xml:"total,attr"`
	PageSize           int      `xml:"pageSize,attr"`
	Page               int      `xml:"page,attr"`
	EdgeGatewayRecords []Link   `xml:"EdgeGatewayRecord"`
}

// GatewayConfiguration ...
type GatewayConfiguration struct {
	XMLName                     xml.Name `xml:"Configuration"`
	GatewayBackingConfiguration string   `xml:"GatewayBackingConfiguration,value"`
	GatewayInterfaces           struct {
		Interfaces []struct {
			Name                string `xml:"Name,value"`
			Network             Link   `xml:"Network"`
			SubnetParticipation struct {
				Gateway   string `xml:"Gateway,value"`
				Netmask   string `xml:"Netmask,value"`
				IPAddress string `xml:"IpAddress,value"`
			} `xml:"SubnetParticipation"`
		} `xml:"GatewayInterface"`
	} `xml:"GatewayInterfaces"`
}

// NetworkGateway ...
type NetworkGateway struct {
	XMLName xml.Name `xml:"EdgeGateway"`
	Href    string   `xml:"href,attr"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

// InstantiateVApp ...
type InstantiateVApp struct {
	XMLName xml.Name `xml:"http://www.vmware.com/vcloud/v1.5 InstantiateVAppTemplateParams"`
	Params  struct {
		NetworkConfig struct {
			NetworkName   string `xml:"networkName,attr"`
			ParentNetwork struct {
				Type string `xml:"type,attr"`
				Name string `xml:"name,attr"`
				Href string `xml:"href,attr"`
			} `xml:"Configuration> ParentNetwork"`
			FenceMode string `xml:"Configuration> FenceMode,value"`
		} `xml:"NetworkConfigSection> NetworkConfig"`
		Source struct {
			Type string `xml:"type,attr"`
			Name string `xml:"name,attr"`
			Href string `xml:"href,attr"`
		} `xml:"Source"`
	} `xml:"InstantiationParams"`
}
