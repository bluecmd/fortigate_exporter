package config

import (
	"flag"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/utils/test"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	validAuthFile      = test.GetRelativeFixturePath("validAuthFile.yaml")
	otherValidAuthFile = test.GetRelativeFixturePath("otherValidAuthFile.yaml")
	invalidAuthFile    = test.GetRelativeFixturePath("invalidAuthFile.yaml")
	certFile           = test.GetRelativeFixturePath("cert.pem")
)

func TestInitDefaultValues(t *testing.T) {
	defaultConfig, err := Init(flag.NewFlagSet("fortigate_exporter", flag.ExitOnError), []string{"fortigate_exporter"})
	assert.NoError(t, err)
	test.Snapshotting.SnapshotT(t, defaultConfig)
}

func TestInitAuthFile(t *testing.T) {
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

			flagSet := flag.NewFlagSet("fortigate_exporter", flag.PanicOnError)
			args := append([]string{"fortigate_exporter"}, aTest.authFileFlags...)

			if aTest.expectToFail {
				assert.Panics(t, func() { _, _ = Init(flagSet, args) })
			} else {
				config, err := Init(flagSet, args)
				req.NoError(err)
				req.Equal(aTest.expected, config.AuthKeys)
				test.Snapshotting.SnapshotT(t, config)
			}
		})
	}
}

func TestInitCaCerts(t *testing.T) {
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

			flagSet := flag.NewFlagSet("fortigate_exporter", flag.PanicOnError)
			args := append([]string{"fortigate_exporter", "-auth-file=" + validAuthFile}, aTest.certFlag...)

			if aTest.expectToFail {
				assert.Panics(t, func() { _, _ = Init(flagSet, args) })
				return
			}

			config, err := Init(flagSet, args)
			req.NoError(err)
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
