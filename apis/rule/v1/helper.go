package v1

import (
	"encoding/json"
	"time"

	"dms/models/rule"

	"github.com/emicklei/go-restful"
)

var (
	SUCCESS     = "success"
	FAILURE     = "failure"
	DefaultStep = 10 * time.Minute
)

const (
	DKYRULE        = "dkyrule"
	AlertsCritical = "critical"
	AlertsWarning  = "warning"
	AlertsInfo     = "info"
	AlertsService  = "service"
	AlertsStorage  = "storage"
	AlertsHost     = "host"
	AlertsLog      = "logging"
)

type reqParams struct {
	group  string
	status string
	level  string
}

func parseRequestParams(req *restful.Request) reqParams {
	var r reqParams
	r.group = req.QueryParameter("group")
	r.status = req.QueryParameter("status")
	r.level = req.QueryParameter("level")
	return r
}

type BaseAPIResponse struct {
	Status    string                 `json:"status,omitempty" description:"api status, success or failure"`
	Error     string                 `json:"error,omitempty" description:"show error reason when status is failure"`
	ErrorList []string               `json:"errorList,omitempty" description:"show more than one error reason when status is failure"`
	Extra     map[string]interface{} `json:"extra,omitempty" description:"show extra information"`
}

type RulesAPIResponse struct {
	BaseAPIResponse
	RulesData map[string][]rule.Rule `json:"rulesData,omitempty" description:"prometheus rules struct"`
}

func convertResponse(err error) []byte {
	response := BaseAPIResponse{
		Status: FAILURE,
		Error:  err.Error(),
	}
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		response.Error = err.Error()
	}
	return jsonBytes
}
