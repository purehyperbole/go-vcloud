package vcloud

import (
	"encoding/xml"
	"log"
	"net/url"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// VApp ...
type VApp struct {
	Connector *Connector `xml:"-"`
	XMLName   xml.Name   `xml:"VApp"`
	Name      string     `xml:"name,attr"`
	Href      string     `xml:"href,attr"`
	Status    string     `xml:"status,attr"`
	Deployed  bool       `xml:"deployed,attr"`
	Links     []t.Link   `xml:"Link"`
	Tasks     *Tasks     `xml:"Tasks"`
	//NetworkConfig t.
}

// NewVApp ...
func NewVApp(c *Connector, href string) (*VApp, error) {
	vURL, err := url.Parse(href)

	if err != nil {
		return nil, err
	}

	resp, err := c.Get(vURL.RequestURI())
	if err != nil {
		return nil, err
	}

	data, err := ParseResponse(resp)
	if err != nil {
		return nil, err
	}

	v := parseVApp(data)
	v.Connector = c

	return v, nil
}

func parseVApp(d *[]byte) *VApp {
	v := VApp{}
	err := xml.Unmarshal(*d, &v)
	if err != nil {
		log.Println(err)
	}
	return &v
}

// GetTasks ...
func (v *VApp) GetTasks() []Task {
	for i := 0; i < len(v.Tasks.Task); i++ {
		v.Tasks.Task[i].Connector = v.Connector
	}
	return v.Tasks.Task
}
