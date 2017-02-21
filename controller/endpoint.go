package controller

import (
	"fmt"
	"net/http"
)

func (xhs *XHttpServer) addQRCodeUrl(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info QRCodeUrlInfo
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}

	err := CreateQRCodeUrlInfo(&info)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("add qrcode url failed: %v", err)
		return response, nil
	}

	return response, nil
}

func (xhs *XHttpServer) getQRCodeUrl(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	
	var info GetQRCodeUrlReq
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	
	userQRCode := &UserQRCodeUrl{
		AppId: info.AppId,
		OpenId: info.OpenId,
		Type: info.T,
	}
	has, err := GetUserQRCodeUrl(userQRCode)
	if err != nil {
		plog.Errorf("get user qrcode url error: %v", err)
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get user qrcode url failed: %v", err)
		return response, nil
	}
	resResult := &GetQRCodeUrlRsp{}
	if has {
		resResult.Status = GET_URL_STATUS_HAS_EXIST
		response.Data = resResult
		return response, nil
	}
	
	count, err := GetQRCodeUrlInfoListCount(info.T)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get qrcode url count failed: %v", err)
		return response, nil
	}
	result, err := GetQRCodeUrlInfoOfRandom(count, info.T)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get qrcode url info failed: %v", err)
		return response, nil
	}
	resResult.Status = GET_URL_STATUS_OK
	resResult.Result = result
	response.Data = resResult
	
	userQRCode.Url = result.Url
	err = CreateUserQRCodeUrl(userQRCode)
	if err != nil {
		plog.Errorf("create user qrcode url error: %v", err)
	}

	return response, nil
}

func (xhs *XHttpServer) expiredQRCodeUrl(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info ExpiredQRCodeReq
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}

	err := UpdateQRCodeUrlInfoStatus(&QRCodeUrlInfo{ID: info.Id, Status: 1})
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("update qrcode url status failed: %v", err)
		return response, nil
	}

	return response, nil
}
