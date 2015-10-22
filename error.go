package vcloud

import (
	"encoding/xml"
	"log"

	t "git.r3labs.io/libraries/go-vcloud/types"
)

// ParseError ...
func ParseError(d *[]byte) *t.Error {
	e := t.Error{}
	err := xml.Unmarshal(*d, &e)
	if err != nil {
		log.Println(err)
	}
	return &e
}
