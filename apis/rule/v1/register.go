package v1

import (
	"net/http"

	"dms/apiserver/runtime"
	"dms/client/rule"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "alerting"
	RespOK    = "ok"
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

func AddToContainer(c *restful.Container, m *rule.Options) error {
	ws := runtime.NewWebService(GroupVersion)

	h := newHandler(*m)

	ws.Route(ws.GET("/rules").
		To(h.HandlerGetAllRules).
		Doc("List all the rules.").
		Param(ws.QueryParameter("group", "A comma-separated list of rule group, default get all rules.").DataType("string")).
		Returns(http.StatusOK, RespOK, RulesAPIResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{DKYRULE}))

	ws.Route(ws.POST("/rules").
		To(h.HandlerAddRules).
		Doc("Add the rule.").
		Param(ws.QueryParameter("group", "rule group").DataType("string")).
		Returns(http.StatusOK, RespOK, RulesAPIResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{DKYRULE}))

	ws.Route(ws.PATCH("/rules/{rule_id}").
		To(h.HandlerUpdateRules).
		Doc("Update the rule.").
		Param(ws.QueryParameter("group", "rule group").DataType("string")).
		Returns(http.StatusOK, RespOK, RulesAPIResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{DKYRULE}))

	ws.Route(ws.DELETE("/rules/{rule_id}").
		To(h.HandlerDeleteRules).
		Doc("Delete the rule.").
		Param(ws.QueryParameter("group", "rule group").DataType("string")).
		Returns(http.StatusOK, RespOK, RulesAPIResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{DKYRULE}))

	ws.Route(ws.POST("/rules/update").
		To(h.HandlerUpdateAllRules).
		Doc("Update all rules.").
		Returns(http.StatusOK, RespOK, RulesAPIResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{DKYRULE}))

	c.Add(ws)
	return nil
}
