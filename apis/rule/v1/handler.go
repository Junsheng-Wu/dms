package v1

import (
	"dms/client/rule"
	"dms/keystone"
	"dms/utils/http"

	mrule "dms/models/rule"

	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

type handler struct {
	operator   mrule.Operator
	ruleOption rule.Options
}

func newHandler(rule rule.Options) *handler {
	operator := mrule.NewOperator(rule.Path, rule.PrometheusEndpoint, rule.PeerAddress)
	return &handler{
		ruleOption: rule,
		operator:   operator,
	}
}

func (h handler) HandlerGetAllRules(req *restful.Request, resp *restful.Response) {
	//params := parseRequestParams(req)
	rules, err := h.operator.ListRules()
	if err != nil {
		keystone.Log("List rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	response := RulesAPIResponse{
		BaseAPIResponse: BaseAPIResponse{
			Status: SUCCESS,
			Extra: map[string]interface{}{
				"count": len(rules),
			},
		},
		RulesData: rules,
	}

	resp.WriteAsJson(response)
}

func (h handler) HandlerAddRules(req *restful.Request, resp *restful.Response) {
	group := req.QueryParameter("group")
	var rule mrule.Rule
	if err := req.ReadEntity(&rule); err != nil {
		keystone.Log("Add rule", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	err := h.operator.AddRule(group, rule)
	if err != nil {
		keystone.Log("Add rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	response := RulesAPIResponse{
		BaseAPIResponse: BaseAPIResponse{
			Status: SUCCESS,
			Extra: map[string]interface{}{
				"rule": rule,
			},
		},
	}

	resp.WriteAsJson(response)
}

func (h handler) HandlerUpdateRules(req *restful.Request, resp *restful.Response) {
	group := req.QueryParameter("group")
	id := req.PathParameter("rule_id")
	var rule mrule.Rule
	if err := req.ReadEntity(&rule); err != nil {
		keystone.Log("Update rule", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	err := h.operator.UpdateRule(group, id, rule)
	if err != nil {
		keystone.Log("List rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	response := RulesAPIResponse{
		BaseAPIResponse: BaseAPIResponse{
			Status: SUCCESS,
			Extra: map[string]interface{}{
				"rule": rule,
			},
		},
	}

	resp.WriteAsJson(response)
}

func (h handler) HandlerDeleteRules(req *restful.Request, resp *restful.Response) {
	group := req.QueryParameter("group")
	id := req.PathParameter("rule_id")

	err := h.operator.DeleteRule(group, id)
	if err != nil {
		keystone.Log("List rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	response := RulesAPIResponse{
		BaseAPIResponse: BaseAPIResponse{
			Status: SUCCESS,
			Extra: map[string]interface{}{
				"id": id,
			},
		},
	}

	resp.WriteAsJson(response)
}

func (h handler) HandlerUpdateAllRules(req *restful.Request, resp *restful.Response) {
	var rules mrule.AlertingRules
	if err := req.ReadEntity(&rules); err != nil {
		keystone.Log("Update all rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}
	err := h.operator.UpdateAllRules(rules)
	if err != nil {
		keystone.Log("Update all rules", err)
		jsonBytes := convertResponse(err)
		http.HandleBadRequest(resp, nil, errors.New(string(jsonBytes)))
		return
	}

	response := RulesAPIResponse{
		BaseAPIResponse: BaseAPIResponse{
			Status: SUCCESS,
			Extra: map[string]interface{}{
				"rules": rules,
			},
		},
	}

	resp.WriteAsJson(response)
}
