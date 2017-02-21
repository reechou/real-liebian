package controller

const (
	QRCODE_STATUS_OK = iota
	QRCODE_STATUS_OFFLINE
)

type GetQRCodeUrlReq struct {
	AppId  string `json:"appId,omitempty"`
	OpenId string `json:"openId,omitempty"`
	T      int64  `json:"type"`
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
