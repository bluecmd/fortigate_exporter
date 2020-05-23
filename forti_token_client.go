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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type fortiTokenClient struct {
	tgt url.URL
	hc  *http.Client
	ctx context.Context
}

func (c *fortiTokenClient) Get(path string, obj interface{}) error {
	u := c.tgt
	u.Path = path

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req = req.WithContext(c.ctx)
	resp, err := c.hc.Do(req)
	if err != nil {
		// Strip request details to hide sensitive access tokens
		uerr := err.(*url.Error)
		return uerr.Err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code was %d, expected 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func newFortiTokenClient(ctx context.Context, tgt url.URL, hc *http.Client, token string) (*fortiTokenClient, error) {
	// TODO(bluecmd): Switch this to use the "Authorization: Bearer xx" header instead (less likely to be leaked in debug logs)
	u := tgt
	q := u.Query()
	q.Add("access_token", token)
	u.RawQuery = q.Encode()
	return &fortiTokenClient{u, hc, ctx}, nil
}
