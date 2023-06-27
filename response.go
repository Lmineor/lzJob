package lzJob

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

func BadRequestResp(response *restful.Response, err error) {
	badReq := msgInfo{Msg: err.Error(), Status: http.StatusBadRequest}
	response.WriteHeaderAndJson(badReq.Status, badReq, restful.MIME_JSON)
}

type msgInfo struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func InternalErrorResp(response *restful.Response, err error) {
	internal := msgInfo{
		Status: http.StatusInternalServerError,
		Msg:    err.Error(),
	}
	response.WriteHeaderAndJson(internal.Status, internal, restful.MIME_JSON)
}
func SuccessResp(response *restful.Response) {
	su := success{Msg: "ok", Status: http.StatusOK}
	response.WriteHeaderAndJson(su.Status, su, restful.MIME_JSON)
}
func SuccessRespWithData(response *restful.Response, data interface{}) {
	successWithData := successWithData{
		success: success{Msg: "ok", Status: http.StatusOK},
		Data:    data,
	}
	response.WriteHeaderAndJson(successWithData.Status, successWithData, restful.MIME_JSON)
}

type success struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type successWithData struct {
	success
	Data interface{} `json:"data"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
