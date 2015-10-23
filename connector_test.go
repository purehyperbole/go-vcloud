package vcloud

import (
	"bufio"
	"errors"
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

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(f)
	_, err = buffer.Read(bytes)

	return bytes, nil
}

func TestAuthenticate(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	Convey("Given a valid HTTP Get Request", t, func() {

	})
}
