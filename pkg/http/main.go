// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus-community/fortigate_exporter/internal/config"
)

type FortiHTTP interface {
	Get(path string, query string, obj interface{}) error
}

func NewFortiClient(ctx context.Context, tgt url.URL, hc *http.Client, aConfig config.FortiExporterConfig) (FortiHTTP, error) {

	auth, ok := aConfig.AuthKeys[config.Target(tgt.String())]
	if !ok {
		return nil, fmt.Errorf("no API authentication registered for %q", tgt.String())
	}

	if auth.Token != "" {
		if tgt.Scheme != "https" {
			return nil, fmt.Errorf("FortiOS only supports token for HTTPS connections")
		}
		c, err := newFortiTokenClient(ctx, tgt, hc, auth.Token)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, fmt.Errorf("invalid authentication data for %q", tgt.String())
}

func Configure(config config.FortiExporterConfig) error {
	roots, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Unable to fetch system CA store: %v", err)
		return err
	}
	for _, cert := range config.TlsExtraCAs {

		if ok := roots.AppendCertsFromPEM(cert.Content); !ok {
			return fmt.Errorf("failed to append certs from PEM %q, unknown error", cert.Path)
		}
	}
	tc := &tls.Config{RootCAs: roots}
	if config.TLSInsecure {
		tc.InsecureSkipVerify = true
	}
	http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout = time.Duration(config.TLSTimeout) * time.Second
	http.DefaultTransport.(*http.Transport).TLSClientConfig = tc
	return nil
}
