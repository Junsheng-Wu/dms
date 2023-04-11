package rule

import (
	"github.com/spf13/pflag"
)

type Options struct {
	Path               string `json:"path,omitempty" yaml:"path,omitempty"`
	PeerAddress        string `json:"peer-address,omitempty" yaml:"peer-address,omitempty"`
	PrometheusEndpoint string `json:"prometheus-endpoint,omitempty" yaml:"prometheus-endpoint,omitempty"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Validate() []error {
	var errs []error

	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet, s *Options) {
	fs.StringVar(&o.Path, "rule-path", s.Path, ""+
		"prometheus rule file path.")
	fs.StringVar(&o.PeerAddress, "peer-address", s.PeerAddress, ""+
		"prometheus peer address.")
	fs.StringVar(&o.PrometheusEndpoint, "prometheus-endpoint", s.PrometheusEndpoint, ""+
		"prometheus endpoint.")
}
