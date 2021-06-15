// Tests of forti_token_client
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type fakeHTTPClient struct {
	status int
	body   string
}

func (c *fakeHTTPClient) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(c.body)),
		StatusCode: c.status,
	}, nil
}

func newClient(sc int, b string) (*fortiTokenClient, error) {
	return newFortiTokenClient(
		context.Background(),
		url.URL{Scheme: "https", Host: "localhost"},
		&fakeHTTPClient{sc, b},
		"TEST-TOKEN",
	)
}

func TestGetParse(t *testing.T) {
	c, _ := newClient(200, `{ "data": "test" }`)
	type D struct {
		Data string
	}
	var v D
	exp := D{"test"}
	if err := c.Get("test", "", &v); err != nil || !reflect.DeepEqual(v, exp) {
		t.Errorf("Get() %v, %v, expected %v, nil", v, err, exp)
	}
}

func TestGetFail(t *testing.T) {
	c, _ := newClient(404, `{}`)
	type D struct {
		Data string
	}
	err := c.Get("test", "", &D{})
	if err == nil {
		t.Errorf("Get() expected non-nil error, got nil error")
	}
}
