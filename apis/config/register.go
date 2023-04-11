package config

import (
	"dms/apiserver/config"
	"dms/apiserver/runtime"

	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func AddToContainer(c *restful.Container, config *config.Config) error {
	ws := runtime.NewWebService(schema.GroupVersion{Version: "configz"})
	ws.Route(ws.GET("/").
		To(func(request *restful.Request, response *restful.Response) {
			response.WriteAsJson(config.ToMap())
		}).
		Doc("Information about the server configuration"))

	c.Add(ws)
	return nil
}
