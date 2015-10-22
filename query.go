package vcloud

import (
	"fmt"
	"net/http"
	"net/url"
)

// Query ...
type Query struct {
	Connector *Connector `xml:"-"`
	Type      string
	Format    string
	Filter    string
	FilterArg string
	Results   string
}

func (q *Query) buildQueryURL() string {
	query := url.Values{}
	query.Add("type", q.Type)
	query.Add("format", q.Format)
	query.Add("filter", fmt.Sprintf("%s==%s", q.Filter, q.FilterArg))
	fullQuery := fmt.Sprintf("/api/query?%s", query.Encode())
	href, _ := url.QueryUnescape(fullQuery)
	return href
}

// Run ...
func (q *Query) Run() *http.Response {
	href := q.buildQueryURL()
	resp, err := q.Connector.Get(href)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
