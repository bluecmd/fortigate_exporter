package config

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type FortiExporterParameter struct {
	AuthFile      *string
	Listen        *string
	ScrapeTimeout *int
	TLSTimeout    *int
	TLSInsecure   *bool
	TlsExtraCAs   *string
	MaxBGPPaths   *int
	MaxVPNUsers   *int
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

var (
	parameter = FortiExporterParameter{
		AuthFile:      flag.String("auth-file", "fortigate-key.yaml", "file containing the authentication map to use when connecting to a Fortigate device"),
		Listen:        flag.String("listen", ":9710", "address to listen on"),
		ScrapeTimeout: flag.Int("scrape-timeout", 30, "max seconds to allow a scrape to take"),
		TLSTimeout:    flag.Int("https-timeout", 10, "TLS Handshake timeout in seconds"),
		TLSInsecure:   flag.Bool("insecure", false, "Allow insecure certificates"),
		TlsExtraCAs:   flag.String("extra-ca-certs", "", "comma-separated files containing extra PEMs to trust for TLS connections in addition to the system trust store"),
		MaxBGPPaths:   flag.Int("max-bgp-paths", 10000, "How many BGP Paths to receive when counting routes, needs to be greater than or equal to the number of routes or metrics will not be generated"),
		MaxVPNUsers:   flag.Int("max-vpn-users", 0, "How many VPN Users to receive when counting users, needs to be greater than or equal the number of users or metrics will not be generated (0 eq. none by default)"),
	}

	savedConfig *FortiExporterConfig
)

func Init() error {
	// check if already parsed
	if savedConfig != nil {
		return nil
	}
	return ReInit()
}

func MustReInit() {
	if err := ReInit(); err != nil {
		log.Fatalf("config.ReInit failed: %+v", err)
	}
}
func ReInit() error {
	flag.Parse()

	savedConfig = &FortiExporterConfig{
		Listen:        *parameter.Listen,
		ScrapeTimeout: *parameter.ScrapeTimeout,
		TLSTimeout:    *parameter.TLSTimeout,
		TLSInsecure:   *parameter.TLSInsecure,
		MaxBGPPaths:   *parameter.MaxBGPPaths,
		MaxVPNUsers:   *parameter.MaxVPNUsers,
	}

	// parse AuthKeys
	af, err := ioutil.ReadFile(*parameter.AuthFile)
	if err != nil {
		log.Fatalf("Failed to read API authentication map file: %v", err)
		return err
	}

	if err := yaml.Unmarshal(af, &savedConfig.AuthKeys); err != nil {
		log.Fatalf("Failed to parse API authentication map file: %v", err)
		return err
	}

	log.Printf("Loaded %d API keys", len(savedConfig.AuthKeys))

	// parse ExtraCAs
	for _, eca := range strings.Split(*parameter.TlsExtraCAs, ",") {
		if eca == "" {
			continue
		}

		certs, err := ioutil.ReadFile(eca)
		if err != nil {
			log.Fatalf("Failed to read extra CA file %q: %v", eca, err)
			return err
		}

		certObject := LocalCert{
			Path:    eca,
			Content: certs,
		}
		savedConfig.TlsExtraCAs = append(savedConfig.TlsExtraCAs, certObject)
	}

	return nil
}

func GetConfig() FortiExporterConfig {
	return *savedConfig
}
