package options

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"dms/utils/net"
)

type ServerRunOptions struct {
	BindAddress   string `json:"bindAddress" yaml:"bindAddress"`
	InsecurePort  int    `json:"insecurePort" yaml:"insecurePort"`
	SecurePort    int    `json:"securePort,omitempty" yaml:"securePort,omitempty"`
	TlsCertFile   string `json:"tlsCertFile,omitempty" yaml:"tlsCertFile,omitempty"`
	TlsPrivateKey string `json:"tlsPrivateKey,omitempty" yaml:"tlsPrivateKey,omitempty"`
}

func NewServerRunOptions() *ServerRunOptions {
	// create default server run options
	return &ServerRunOptions{
		BindAddress:  "0.0.0.0",
		InsecurePort: 9200,
	}
}

func (s *ServerRunOptions) Validate() []error {
	var errs []error

	if s.SecurePort == 0 && s.InsecurePort == 0 {
		errs = append(errs, fmt.Errorf("insecure and secure port can not be disabled at the same time"))
	}

	if net.IsValidPort(s.SecurePort) {
		if s.TlsCertFile == "" {
			errs = append(errs, fmt.Errorf("tls cert file is empty while secure serving"))
		} else {
			if _, err := os.Stat(s.TlsCertFile); err != nil {
				errs = append(errs, err)
			}
		}

		if s.TlsPrivateKey == "" {
			errs = append(errs, fmt.Errorf("tls private key file is empty while secure serving"))
		} else {
			if _, err := os.Stat(s.TlsPrivateKey); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet, c *ServerRunOptions) {
	fs.StringVar(&s.BindAddress, "bind-address", c.BindAddress, "server bind address")
	fs.IntVar(&s.InsecurePort, "insecure-port", c.InsecurePort, "insecure port number")
	fs.IntVar(&s.SecurePort, "secure-port", s.SecurePort, "secure port number")
	fs.StringVar(&s.TlsCertFile, "tls-cert-file", c.TlsCertFile, "tls cert file")
	fs.StringVar(&s.TlsPrivateKey, "tls-private-key", c.TlsPrivateKey, "tls private key")
}
