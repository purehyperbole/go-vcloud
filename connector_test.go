package vcloud

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func loadFixture(file string) ([]byte, error) {
	path, _ := filepath.Abs(file)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(f)
	_, err = buffer.Read(bytes)

	return bytes, nil
}

func parseRequest(r *http.Request) *[]byte {
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	return &data
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, _ := r.BasicAuth()
	if r.Header.Get("accept") == "application/*+xml;version=5.1" &&
		user == "test@test" &&
		pass == "test" {
		w.Header().Set("x-vcloud-authorization", "test")
	} else if r.Header.Get("x-vcloud-authorization") != "test" {
		message, _ := loadFixture("fixtures/autherror.xml")
		authErr := errors.New(string(message))
		http.Error(w, authErr.Error(), 403)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/api/org" {
		payload, _ := loadFixture("fixtures/orglist.xml")
		w.Write(payload)
	} else if r.RequestURI != "/api/sessions" {
		payload, _ := loadFixture("fixtures/resourcenotfound.xml")
		http.Error(w, string(payload), 404)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/api/vdc/test/action/instantiateVAppTemplate" {
		if r.Method == "POST" {
			body := parseRequest(r)
			if r.Header.Get("Content-Type") != "application/vnd.vmware.vcloud.instantiateVAppTemplateParams+xml" {
				payload, _ := loadFixture("fixtures/resourcenotfound.xml")
				http.Error(w, string(payload), 404)
			} else if len(*body) == 0 {
				payload, _ := loadFixture("fixtures/eoferror.xml")
				http.Error(w, string(payload), 400)
			} else {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("Accepted!"))
			}
		} else {
			payload, _ := loadFixture("fixtures/methodnotallowed.xml")
			http.Error(w, string(payload), 405)
		}
	} else if r.RequestURI != "/api/sessions" {
		payload, _ := loadFixture("fixtures/resourcenotfound.xml")
		http.Error(w, string(payload), 404)
	}
}

func TestAuthenticate(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		authHandler(w, r)
	}))
	defer ts.Close()

	tsURL, _ := url.Parse(ts.URL)

	Convey("Given an Authorization attempt", t, func() {
		Convey("When using valid credentials", func() {
			cf := Config{
				URL:           tsURL.Host,
				Username:      "test@test",
				Password:      "test",
				SSLSkipVerify: true,
			}
			c := NewConnector(&cf)
			err := c.Authenticate()
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Auth Token should be stored", func() {
				So(c.AuthToken, ShouldEqual, "test")
			})
		})
		Convey("When using invalid credentials", func() {
			cf := Config{
				URL:           tsURL.Host,
				Username:      "tset@tset",
				Password:      "tset",
				SSLSkipVerify: true,
			}
			c := NewConnector(&cf)
			err := c.Authenticate()
			Convey("Error should be returned", func() {
				So(err.Error(), ShouldEqual, "Access is forbidden")
			})
			Convey("Auth Token should not be stored", func() {
				So(c.AuthToken, ShouldBeBlank)
			})
		})
	})

}

func TestGet(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		authHandler(w, r)
		getHandler(w, r)
	}))
	defer ts.Close()
	tsURL, _ := url.Parse(ts.URL)

	Convey("Given an HTTP Get Request", t, func() {
		cf := Config{
			URL:           tsURL.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}
		c := NewConnector(&cf)
		authErr := c.Authenticate()
		Convey("Given a valid request", func() {
			resp, err := c.Get("/api/org")
			Convey("We should be authenticated", func() {
				var message string
				if authErr != nil {
					message = authErr.Error()
				}
				So(message, ShouldNotEqual, "Access is forbidden")
			})
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("There should be a payload returned", func() {
				So(resp, ShouldNotBeNil)
			})
			Convey("There should be an xml body", func() {
				payload, _ := loadFixture("fixtures/orglist.xml")
				respBody, _ := ParseResponse(resp)
				So(len(*respBody), ShouldEqual, len(payload))
			})
		})
		Convey("Given an invalid request", func() {
			resp, err := c.Get("/api/test")
			Convey("There should be an error", func() {
				var message string
				if err != nil {
					message = err.Error()
				}
				So(err, ShouldNotBeNil)
				So(message, ShouldEqual, "Resource not found")
			})
			Convey("There should be no response body", func() {
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestPost(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		authHandler(w, r)
		postHandler(w, r)
	}))
	defer ts.Close()
	tsURL, _ := url.Parse(ts.URL)

	Convey("Given an HTTP Post Request", t, func() {
		cf := Config{
			URL:           tsURL.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}
		c := NewConnector(&cf)
		authErr := c.Authenticate()

		Convey("Given a valid request", func() {
			data := []byte("test request")
			resp, err := c.Post("/api/vdc/test/action/instantiateVAppTemplate", data, "application/vnd.vmware.vcloud.instantiateVAppTemplateParams+xml")
			Convey("We should be authenticated", func() {
				var message string
				if authErr != nil {
					message = authErr.Error()
				}
				So(message, ShouldNotEqual, "Access is forbidden")
			})
			Convey("We should not receive an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("We should receive a payload", func() {
				So(resp, ShouldNotBeNil)
			})
		})
	})
}
