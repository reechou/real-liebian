package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"math/rand"
	"time"
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

func (xhs *XHttpServer) addQRCodeUrlList(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info []QRCodeUrlInfo
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	
	err := CreateQRCodeUrlInfoList(info)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("add qrcode url list failed: %v", err)
		return response, nil
	}
	
	return response, nil
}

func (xhs *XHttpServer) getAllQRCodeUrlInfoList(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info GetAllQRCodeFromType
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	
	count, err := GetAllQRCodeUrlInfoFromTypeCount(info.Type)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get all qrcode count failed: %v", err)
		return response, nil
	}
	list, err := GetAllQRCodeUrlInfoFromType(info.Type, info.Offset, info.Num)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get all qrcode url list failed: %v", err)
		return response, nil
	}
	result := &GetAllQRCodeRsp{
		Count: count,
		List: list,
	}
	response.Data = result
	
	return response, nil
}
func (xhs *XHttpServer) getActiveQRCodeUrlInfoList(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info GetFromType
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	plog.Debugf("getActiveQRCodeUrlInfoList req: %v", info)
	
	list, err := GetQRCodeUrlListFromType(info.Type)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get active qrcode url info list failed: %v", err)
		return response, nil
	}
	response.Data = list
	
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
		AppId:  info.AppId,
		OpenId: info.OpenId,
		Type:   info.T,
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
		resResult.Result = &QRCodeUrlInfo{Url: userQRCode.Url}
		response.Data = resResult
		return response, nil
	}

	list, err := GetQRCodeUrlListFromType(info.T)
	if err != nil {
		plog.Errorf("get qrcode list from type error: %v", err)
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get qrcode list from type error: %v", err)
		return response, nil
	}
	if len(list) == 0 {
		plog.Errorf("get qrcode list from type is nil")
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("get qrcode list from type is nil")
		return response, nil
	}
	rand.Seed(time.Now().UnixNano())
	offset := rand.Intn(len(list))
	resResult.Status = GET_URL_STATUS_OK
	resResult.Result = &list[offset]
	response.Data = resResult

	userQRCode.Url = resResult.Result.Url
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

func (xhs *XHttpServer) createGroupSetting(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info TypeGroupSetting
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}

	err := CreateTypeGroupSetting(&info)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("create type group setting failed: %v", err)
		return response, nil
	}

	return response, nil
}

func (xhs *XHttpServer) createRobotMsgSetting(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info CreateRobotMsgSettingReq
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}

	msg, _ := json.Marshal(info.Msg)
	setting := TypeRobotMsgSetting{
		Type:        info.Type,
		SettingType: info.SettingType,
		Robot:       info.Robot,
		Msg:         string(msg),
		Interval:    info.Interval,
	}
	err := CreateTypeRobotMsgSetting(&setting)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("create type robot msg setting failed: %v", err)
		return response, nil
	}
	acg := NewAutoCheckGroup(setting, xhs.rbExt, xhs.rul)
	go acg.Run()
	xhs.Lock()
	xhs.acMap[setting.ID] = acg
	xhs.Unlock()

	return response, nil
}

func (xhs *XHttpServer) delRobotMsgSetting(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info DelRobotMsgSettingReq
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	
	err := DelTypeRobotMsgSetting(info.Id)
	if err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("del type robot msg setting failed: %v", err)
		return response, nil
	}
	xhs.Lock()
	defer xhs.Unlock()
	acg := xhs.acMap[info.Id]
	if acg != nil {
		acg.Stop()
		delete(xhs.acMap, info.Id)
	}
	
	return response, nil
}

func (xhs *XHttpServer) RobotReceiveMsg(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	response := &Response{Code: RES_OK}
	var info ReceiveMsgInfo
	if err := xhs.decodeBody(req, &info, nil); err != nil {
		response.Code = RES_ERR
		response.Msg = fmt.Sprintf("Request decode failed: %v", err)
		return response, nil
	}
	plog.Debugf("receive robot msg: %v", info)
	qrCodeInfo, ok, err := xhs.getNowActiveGroup(info.BaseInfo.FromGroupName)
	if err != nil {
		plog.Errorf("get now active group error: %v", err)
		return response, nil
	}
	if ok {
		xhs.handleRobotMsg(qrCodeInfo, &info)
	}

	return response, nil
}

func (xhs *XHttpServer) handleRobotMsg(qrCodeUrlInfo *QRCodeUrlInfo, msg *ReceiveMsgInfo) {
	switch msg.BaseInfo.ReceiveEvent {
	case RECEIVE_EVENT_MSG:
		if msg.MsgType == MSG_TYPE_IMG {
			xhs.rul.AddGroupImgUser(qrCodeUrlInfo.ID, "@"+msg.BaseInfo.FromNickName)
		}
	case RECEIVE_EVENT_MOD_GROUP_ADD:
		maxNum := xhs.logic.cfg.GroupMaxNum
		if maxNum == 0 {
			maxNum = 80
		}
		if msg.GroupMemberNum >= maxNum {
			xhs.changeActiveGroup(qrCodeUrlInfo)
			return
		}
		qrCodeUrlInfo.IfMod = 1
		UpdateQRCodeUrlInfoIfMod(qrCodeUrlInfo)
	}
}

func (xhs *XHttpServer) changeActiveGroup(qrCodeUrlInfo *QRCodeUrlInfo) {
	plog.Infof("qrcodeurlinfo[%v] change active group", qrCodeUrlInfo)
	qrCodeUrlInfo.Status = 1
	UpdateQRCodeUrlInfoStatus(qrCodeUrlInfo)
	xhs.rul.DelGroup(qrCodeUrlInfo.ID)
}

func (xhs *XHttpServer) getNowActiveGroup(group string) (*QRCodeUrlInfo, bool, error) {
	qrcodeUrlInfo := &QRCodeUrlInfo{
		Name: group,
	}
	has, err := GetQRCodeUrlInfo(qrcodeUrlInfo)
	if err != nil {
		plog.Errorf("get qrcode info error: %v", err)
		return nil, false, err
	}
	if !has {
		return nil, false, nil
	}
	groupList, err := GetQRCodeUrlListFromType(qrcodeUrlInfo.Type)
	if err != nil {
		plog.Errorf("get qrcode list from type error: %v", err)
		return nil, false, err
	}
	for _, v := range groupList {
		if v.Name == qrcodeUrlInfo.Name {
			return qrcodeUrlInfo, true, nil
		}
	}
	return nil, false, nil
}
