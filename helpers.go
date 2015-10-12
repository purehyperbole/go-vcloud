package vcloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ParseResponse ...
func ParseResponse(resp *http.Response) (*[]byte, error) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
	return &data, err
}
