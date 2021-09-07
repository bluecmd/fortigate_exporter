package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type FortiExporterParameter struct {
	AuthKeys      *AuthKeys
	Listen        *string
	ScrapeTimeout *int
	TLSTimeout    *int
	TLSInsecure   *bool
	TlsExtraCAs   *[]LocalCert
	MaxBGPPaths   *int
	MaxVPNUsers   *int
}

func (p FortiExporterParameter) ParseToConfig() FortiExporterConfig {
	savedConfig := FortiExporterConfig{
		AuthKeys:      *p.AuthKeys,
		Listen:        *p.Listen,
		ScrapeTimeout: *p.ScrapeTimeout,
		TLSTimeout:    *p.TLSTimeout,
		TLSInsecure:   *p.TLSInsecure,
		TlsExtraCAs:   *p.TlsExtraCAs,
		MaxBGPPaths:   *p.MaxBGPPaths,
		MaxVPNUsers:   *p.MaxVPNUsers,
	}
	return savedConfig
}

type FortiExporterConfig struct {
	AuthKeys      AuthKeys
	Listen        string
	ScrapeTimeout int
	TLSTimeout    int
	TLSInsecure   bool
	TlsExtraCAs   []LocalCert
	MaxBGPPaths   int
	MaxVPNUsers   int
}

type AuthKeys map[Target]TargetAuth

type Target string
type Token string
type ProbeList []string

type Probes struct {
	Include ProbeList
	Exclude ProbeList
}

type TargetAuth struct {
	Token  Token
	Probes Probes
}

type LocalCert struct {
	Path    string
	Content []byte
}

func setParameterFlags(flagSet *flag.FlagSet) *FortiExporterParameter {
	localParameter := FortiExporterParameter{
		Listen:        flagSet.String("listen", ":9710", "address to listen on"),
		ScrapeTimeout: flagSet.Int("scrape-timeout", 30, "max seconds to allow a scrape to take"),
		TLSTimeout:    flagSet.Int("https-timeout", 10, "TLS Handshake timeout in seconds"),
		TLSInsecure:   flagSet.Bool("insecure", false, "Allow insecure certificates"),
		MaxBGPPaths:   flagSet.Int("max-bgp-paths", 10000, "How many BGP Paths to receive when counting routes, needs to be higher then the number of routes or metrics will not be generated"),
		MaxVPNUsers:   flagSet.Int("max-vpn-users", 0, "How many VPN Users to receive when counting users, needs to be greater than or equal the number of users or metrics will not be generated (0 eq. none by default)"),
		// defaults
		TlsExtraCAs: new([]LocalCert),
		AuthKeys:    new(AuthKeys),
	}
	flagSet.Func("extra-ca-certs", "comma-separated files containing extra PEMs to trust for TLS connections in addition to the system trust store. Multiple flags will be concatenated", func(s string) error {
		// parse ExtraCAs
		for _, eca := range strings.Split(s, ",") {
			if eca == "" {
				continue
			}

			certs, err := ioutil.ReadFile(eca)
			if err != nil {
				return fmt.Errorf("failed to read extra CA file %q: %v", eca, err)
			}

			certObject := LocalCert{
				Path:    eca,
				Content: certs,
			}
			localCerts := append(*localParameter.TlsExtraCAs, certObject)
			localParameter.TlsExtraCAs = &localCerts
		}
		return nil
	})

	flagSet.Func("auth-file", "file containing the authentication map to use when connecting to a Fortigate device", func(s string) error {
		// parse AuthKeys
		af, err := ioutil.ReadFile(s)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(af, &localParameter.AuthKeys); err != nil {
			return err
		}
		return nil
	})

	return &localParameter
}

func Init(flagSet *flag.FlagSet, args []string) (FortiExporterConfig, error) {
	parameters := setParameterFlags(flagSet)
	err := flagSet.Parse(args[1:])
	if err != nil {
		return FortiExporterConfig{}, err
	}
	return parameters.ParseToConfig(), nil
}
