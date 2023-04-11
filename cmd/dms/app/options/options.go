package options

import (
	"crypto/tls"
	"dms/server/options"
	"flag"
	"fmt"
	"net/http"
	"strings"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"

	"dms/apiserver"
	apiserverconfig "dms/apiserver/config"
)

type ServerRunOptions struct {
	ConfigFile              string
	GenericServerRunOptions *options.ServerRunOptions
	*apiserverconfig.Config
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		GenericServerRunOptions: options.NewServerRunOptions(),
		Config:                  apiserverconfig.New(),
	}
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	s.GenericServerRunOptions.AddFlags(fs, s.GenericServerRunOptions)
	s.RuleOptions.AddFlags(fss.FlagSet("rule"), s.RuleOptions)
	//s.AlertingOptions.AddFlags(fss.FlagSet("alerting"), s.AlertingOptions)
	//s.AuthenticationOptions.AddFlags(fss.FlagSet("authentication"), s.AuthenticationOptions)
	//s.AuthorizationOptions.AddFlags(fss.FlagSet("authorization"), s.AuthorizationOptions)

	//fs = fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})

	return fss
}

// NewAPIServer creates an APIServer instance using given options
func (s *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	apiServer := &apiserver.APIServer{
		Config: s.Config,
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.GenericServerRunOptions.InsecurePort),
	}

	if s.GenericServerRunOptions.SecurePort != 0 {
		certificate, err := tls.LoadX509KeyPair(s.GenericServerRunOptions.TlsCertFile, s.GenericServerRunOptions.TlsPrivateKey)
		if err != nil {
			return nil, err
		}
		server.TLSConfig.Certificates = []tls.Certificate{certificate}
	}

	apiServer.Server = server

	return apiServer, nil
}
