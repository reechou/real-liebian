package controller

import (
	"fmt"
	"net/http"
)

func (xhs *XHttpServer) addQRCodeUrl(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info QRCodeUrl
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}

	err := xhs.logic.cdb.InsertQRCodeUrl(&info)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("add qrcode url failed: %v", err)
		return response, nil
	}

	return response, nil
}

func (xhs *XHttpServer) getQRCodeUrl(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}

	info, err := xhs.logic.GetQRCodeUrl()
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get qrcode url failed: %v", err)
		return response, nil
	}
	response.Data = info

	return response, nil
}
