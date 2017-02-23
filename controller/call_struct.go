package controller

const (
	QRCODE_STATUS_OK = iota
	QRCODE_STATUS_ACTIVE
)

type GetQRCodeUrlReq struct {
	AppId  string `json:"appId,omitempty"`
	OpenId string `json:"openId,omitempty"`
	T      int64  `json:"type"`
}

type CreateRobotMsgSettingReq struct {
	Type        int64     `json:"type"`
	SettingType int       `json:"settingType"`
	Robot       string    `json:"robot"`
	Msg         []MsgInfo `json:"msg"`
	Interval    int64     `json:"interval"`
}

type DelRobotMsgSettingReq struct {
	Id int64 `json:"id"`
}

type GetFromType struct {
	Type int64 `json:"type"`
}

type GetAllQRCodeFromType struct {
	Type int64 `json:"type"`
	Offset int64 `json:"offset"`
	Num int64 `json:"num"`
}

type GetAllQRCodeRsp struct {
	Count int64 `json:"count"`
	List  []QRCodeUrlInfo `json:"list"`
}

const (
	GET_URL_STATUS_OK = iota
	GET_URL_STATUS_HAS_EXIST
)

type GetQRCodeUrlRsp struct {
	Status int            `json:"status"`
	Result *QRCodeUrlInfo `json:"result,omitempty"`
}

type ExpiredQRCodeReq struct {
	Id int64 `json:"id"`
}

const (
	RES_OK = iota
	RES_ERR
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}