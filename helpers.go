package vcloud

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

// ParseResponse ...
func ParseResponse(resp *http.Response, i interface{}) error {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(data))

	err = xml.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	return err
}
