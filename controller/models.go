package controller

const (
	QRCODE_STATUS_OK = iota
	QRCODE_STATUS_OFFLINE
)

type QRCodeUrl struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Status     int64  `json:"-"`
	CreateTime int64  `json:"-"`
}

type QRCodeUrlList struct {
	List       []*QRCodeUrl
	UpdateTime int64
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
