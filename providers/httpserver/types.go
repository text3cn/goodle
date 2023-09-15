package httpserver

import "net/http"

type ReqStruct struct {
	IRequest
	request *http.Request
}
type RespStruct struct {
	request        *ReqStruct
	responseWriter http.ResponseWriter
}
