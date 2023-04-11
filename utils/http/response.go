package http

import (
	"net/http"
	"runtime"

	"github.com/emicklei/go-restful"
	"k8s.io/klog"
)

func HandleInternalError(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	_ = response.WriteError(http.StatusInternalServerError, err)
}

// HandleBadRequest writes http.StatusBadRequest and log error
func HandleBadRequest(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	_ = response.WriteError(http.StatusBadRequest, err)
}

func HandleNotFound(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	_ = response.WriteError(http.StatusNotFound, err)
}

func HandleForbidden(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	_ = response.WriteError(http.StatusForbidden, err)
}

func HandleConflict(response *restful.Response, req *restful.Request, err error) {
	_, fn, line, _ := runtime.Caller(1)
	klog.Errorf("%s:%d %v", fn, line, err)
	_ = response.WriteError(http.StatusConflict, err)
}
