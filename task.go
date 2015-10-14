package vcloud

import (
	"encoding/xml"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// Task ...
type Task struct {
	XMLName       xml.Name `xml:"Task"`
	Name          string   `xml:"name,attr"`
	Href          string   `xml:"href,attr"`
	OperationName string   `xml:"operationName,attr"`
	Status        string   `xml:"status,attr"`
	StartTime     string   `xml:"startTime,attr"`
	ExpiryTime    string   `xml:"expiryTime,attr"`
	Cancel        bool     `xml:"cancelRequested"`
	Links         []t.Link `xml:"Link"`
	Owner         t.Link   `xml:"Owner"`
	User          t.Link   `xml:"User"`
	Organization  t.Link   `xml:"Organization"`
}
