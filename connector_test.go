package vcloud

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/julienschmidt/httprouter"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup() {
	router := httprouter.New()

	router.POST("/api/sessions", sessionHandler)
	router.GET("/test", getHandler)
	router.POST("/test", postHandler)
	router.PUT("/test", postHandler)
	router.DELETE("/test", deleteHandler)
	router.NotFound = notFoundHandler

	server = httptest.NewTLSServer(router)
}

func teardown() {
	server.Close()
}

func parseRequest(r *http.Request) *[]byte {
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	return &data
}

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

func sessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/xml")
	user, pass, _ := r.BasicAuth()
	if r.Header.Get("accept") == "application/*+xml;version=5.1" &&
		user == "test@test" &&
		pass == "test" {
		w.Header().Set("x-vcloud-authorization", "test")
	} else {
		message, _ := loadFixture("fixtures/autherror.xml")
		authErr := errors.New(string(message))
		http.Error(w, authErr.Error(), 403)
	}
}

func auth(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("x-vcloud-authorization") != "test" {
		message, _ := loadFixture("fixtures/autherror.xml")
		authErr := errors.New(string(message))
		http.Error(w, authErr.Error(), 403)
		return false
	}
	return true
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := loadFixture("fixtures/resourcenotfound.xml")
	http.Error(w, string(payload), 404)
}

func getHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !auth(w, r) {
		return
	}
	w.Write([]byte("OK"))
}

func postHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !auth(w, r) {
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	if string(body) == "test request" {
		w.WriteHeader(201)
		w.Write([]byte("OK"))
	} else {
		message, _ := loadFixture("fixtures/eoferror.xml")
		authErr := errors.New(string(message))
		http.Error(w, authErr.Error(), 400)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !auth(w, r) {
		return
	}
	w.WriteHeader(200)
}

func TestAuthenticate(t *testing.T) {
	setup()
	defer teardown()
	tsurl, _ := url.Parse(server.URL)

	Convey("Given an Authorization attempt", t, func() {
		Convey("When using valid credentials", func() {
			cf := Config{
				URL:           tsurl.Host,
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
				URL:           tsurl.Host,
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
	setup()
	defer teardown()
	tsurl, _ := url.Parse(server.URL)

	Convey("Given an HTTP Get Request", t, func() {
		cf := Config{
			URL:           tsurl.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}

		c := NewConnector(&cf)
		authErr := c.Authenticate()

		Convey("Given a valid request", func() {
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			resp, err := c.Get(href)
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
			Convey("There should be an body", func() {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				So(len(body), ShouldEqual, len([]byte("OK")))
			})
		})

		Convey("Given an invalid request", func() {
			href := fmt.Sprintf("https://%s/invalidtest", tsurl.Host)
			resp, err := c.Get(href)
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
	setup()
	defer teardown()
	tsurl, _ := url.Parse(server.URL)

	Convey("Given an HTTP Post Request", t, func() {
		cf := Config{
			URL:           tsurl.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}
		c := NewConnector(&cf)
		authErr := c.Authenticate()

		Convey("Given a valid request", func() {
			data := []byte("test request")
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			resp, err := c.Post(href, data, "application/xml")
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

		Convey("Given an invalid request", func() {
			data := []byte("invalid")
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			resp, err := c.Post(href, data, "application/xml")
			Convey("There should be an error", func() {
				var message string
				if err != nil {
					message = err.Error()
				}
				So(err, ShouldNotBeNil)
				So(message, ShouldContainSubstring, "Bad request")
			})
			Convey("There should be no response body", func() {
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestPut(t *testing.T) {
	setup()
	defer teardown()
	tsurl, _ := url.Parse(server.URL)

	Convey("Given an HTTP Put Request", t, func() {
		cf := Config{
			URL:           tsurl.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}
		c := NewConnector(&cf)
		authErr := c.Authenticate()

		Convey("Given a valid request", func() {
			data := []byte("test request")
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			resp, err := c.Put(href, data, "application/xml")
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

		Convey("Given an invalid request", func() {
			data := []byte("invalid")
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			resp, err := c.Post(href, data, "application/xml")
			Convey("There should be an error", func() {
				var message string
				if err != nil {
					message = err.Error()
				}
				So(err, ShouldNotBeNil)
				So(message, ShouldContainSubstring, "Bad request")
			})
			Convey("There should be no response body", func() {
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestDelete(t *testing.T) {
	setup()
	defer teardown()
	tsurl, _ := url.Parse(server.URL)

	Convey("Given an HTTP Delete Request", t, func() {
		cf := Config{
			URL:           tsurl.Host,
			Username:      "test@test",
			Password:      "test",
			SSLSkipVerify: true,
		}
		c := NewConnector(&cf)
		authErr := c.Authenticate()

		Convey("Given a valid request", func() {
			href := fmt.Sprintf("https://%s/test", tsurl.Host)
			err := c.Delete(href)
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
		})

		Convey("Given an invalid request", func() {
			data := []byte("invalid")
			href := fmt.Sprintf("https://%s/invalidtest", tsurl.Host)
			resp, err := c.Post(href, data, "application/xml")
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
