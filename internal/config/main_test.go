package config

import (
	"flag"
	"os"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/utils/test"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	validAuthFile      = test.GetFixturePathPanic("validAuthFile.yaml")
	otherValidAuthFile = test.GetFixturePathPanic("otherValidAuthFile.yaml")
	invalidAuthFile    = test.GetFixturePathPanic("invalidAuthFile.yaml")
	certFile           = test.GetRelativeFixturePathPanic("cert.pem")
)

func TestInit(t *testing.T) {
	// default values
	t.Run("default values", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("fortigate_exporter", flag.ExitOnError)
		os.Args = []string{"fortigate_exporter"}
		defaultConfig := Init()
		test.Snapshotting.SnapshotT(t, defaultConfig)
	})

	// auth file test
	authFileTests := []struct {
		name          string
		authFileFlags []string
		expected      AuthKeys
		expectToFail  bool
	}{
		{
			name:          "no auth-file",
			authFileFlags: []string{},
			expected:      nil,
		},
		{
			name:          "single auth-file",
			authFileFlags: []string{"-auth-file=" + validAuthFile},
			expected: AuthKeys{
				"https://test.example.com": {
					Token: "abcd",
				},
			},
		},
		{
			name:          "single invalid auth-file",
			authFileFlags: []string{"-auth-file=" + invalidAuthFile},
			expected:      nil,
			expectToFail:  true,
		},
		{
			name:          "overload auth-file",
			authFileFlags: []string{"-auth-file=" + validAuthFile, "-auth-file=" + otherValidAuthFile},
			expected: AuthKeys{
				"https://test.example.com": {
					Token: "abcd",
				},
				"https://other.example.com": {
					Token: "def",
				},
			},
		},
		{
			name:          "fail if file is not found",
			authFileFlags: []string{"-auth-file=test.yaml"},
			expected:      nil,
			expectToFail:  true,
		},
	}
	for _, aTest := range authFileTests {
		t.Run(aTest.name, func(t *testing.T) {
			req := require.New(t)

			// reset flags
			flag.CommandLine = flag.NewFlagSet("fortigate_exporter", flag.PanicOnError)

			os.Args = append([]string{"fortigate_exporter"}, aTest.authFileFlags...)
			if aTest.expectToFail {
				assert.Panics(t, func() { Init() })
			} else {
				config := Init()
				req.Equal(aTest.expected, config.AuthKeys)
				test.Snapshotting.SnapshotT(t, config)
			}

		})
	}

	// caCerts
	caCertFlagTests := []struct {
		name          string
		certFlag      []string
		numberOfCerts int
		shouldBeNil   bool
		expectToFail  bool
	}{
		{
			name:          "no ca cert",
			certFlag:      []string{},
			numberOfCerts: 0,
			shouldBeNil:   true,
		},
		{
			name:          "ca cert not found",
			certFlag:      []string{"-extra-ca-certs=test.pem"},
			numberOfCerts: 0,
			expectToFail:  true,
		},
		{
			name:          "single ca cert",
			certFlag:      []string{"-extra-ca-certs=" + certFile},
			numberOfCerts: 1,
		},
		{
			name:          "multiple ca certs with comma separation",
			certFlag:      []string{"-extra-ca-certs=" + certFile + "," + certFile},
			numberOfCerts: 2,
		},
		{
			name:          "multiple ca certs with separate flags",
			certFlag:      []string{"-extra-ca-certs=" + certFile, "-extra-ca-certs=" + certFile},
			numberOfCerts: 2,
		},
		{
			name:          "multiple ca certs with separate flags and comma separate",
			certFlag:      []string{"-extra-ca-certs=" + certFile, "-extra-ca-certs=" + certFile + "," + certFile},
			numberOfCerts: 3,
		},
	}
	for _, aTest := range caCertFlagTests {
		t.Run(aTest.name, func(t *testing.T) {
			req := require.New(t)

			// reset flags
			flag.CommandLine = flag.NewFlagSet("fortigate_exporter", flag.PanicOnError)

			os.Args = append([]string{"fortigate_exporter", "-auth-file=" + validAuthFile}, aTest.certFlag...)
			if aTest.expectToFail {
				assert.Panics(t, func() { Init() })
				return
			}

			config := Init()
			if aTest.shouldBeNil {
				req.Nil(config.TlsExtraCAs)
			} else {
				req.NotNil(config.TlsExtraCAs)
				req.Exactly(len(config.TlsExtraCAs), aTest.numberOfCerts)
			}
			test.Snapshotting.SnapshotT(t, config)
		})
	}
}
