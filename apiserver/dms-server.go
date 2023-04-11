package apiserver

import (
	"context"
	apiconfig "dms/apis/config"
	rulev1 "dms/apis/rule/v1"
	"dms/apiserver/config"
	"github.com/emicklei/go-restful"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	"net/http"
)

type APIServer struct {
	// number of emla apiserver
	ServerCount int

	Server *http.Server

	Config *config.Config
	// webservice container, where all webservice defines
	container *restful.Container
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()
	s.container.Router(restful.CurlyRouter{})
	s.container.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/var/static/swagger-ui"))))

	s.installAPIGroups()

	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}

	s.Server.Handler = s.container

	//s.buildHandlerChain(stopCh)

	return nil
}

func (s *APIServer) installAPIGroups() {
	urlruntime.Must(apiconfig.AddToContainer(s.container, s.Config))
	urlruntime.Must(rulev1.AddToContainer(s.container, s.Config.RuleOptions))
}

func (s *APIServer) Run(stopCh <-chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stopCh
		_ = s.Server.Shutdown(ctx)
	}()

	klog.V(0).Infof("Start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return err
}
