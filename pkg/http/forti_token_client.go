// HTTP client for Fortigate API using token authentication
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bluecmd/fortigate_exporter/internal/config"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type fortiTokenClient struct {
	tgt url.URL
	hc  HTTPClient
	ctx context.Context
	tok config.Token
}

func (c *fortiTokenClient) newGetRequest(url string) (*http.Request, error) {
	r, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.tok))
	return r, nil
}

func (c *fortiTokenClient) Get(path string, query string, obj interface{}) error {
	u := c.tgt
	u.Path = path
	u.RawQuery = query

	req, err := c.newGetRequest(u.String())
	if err != nil {
		return err
	}

	req = req.WithContext(c.ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code was %d, expected 200 (path: %q)", resp.StatusCode, path)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func (c *fortiTokenClient) String() string {
	return c.tgt.String()
}

func newFortiTokenClient(ctx context.Context, tgt url.URL, hc HTTPClient, token config.Token) (*fortiTokenClient, error) {
	return &fortiTokenClient{tgt, hc, ctx, token}, nil
}
