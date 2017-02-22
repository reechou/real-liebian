package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"sync"

	"github.com/mitchellh/mapstructure"
)

type XHttpServer struct {
	sync.Mutex
	logic *ControllerLogic
	hs    *HttpSrv
	rul   *RobotUserLogic
	acMap map[int64]*AutoCheckGroup
	rbExt *RobotExt
}

type HttpHandler func(rsp http.ResponseWriter, req *http.Request) (interface{}, error)

func NewXHttpServer(addr string, port int, logic *ControllerLogic, rul *RobotUserLogic) *XHttpServer {
	xhs := &XHttpServer{
		hs: &HttpSrv{
			HttpAddr: addr,
			HttpPort: port,
			Routers:  make(map[string]http.HandlerFunc),
		},
		logic: logic,
		rul: rul,
		acMap: make(map[int64]*AutoCheckGroup),
		rbExt: NewRobotExt(logic.cfg),
	}
	xhs.initAutoCheckGroupList()
	xhs.registerHandlers()

	return xhs
}

func (xhs *XHttpServer) Run() {
	xhs.hs.Run()
}

func (xhs *XHttpServer) registerHandlers() {
	xhs.hs.Route("/", xhs.Index)

	xhs.hs.Route("/liebian/add_qrcode_url", xhs.httpWrap(xhs.addQRCodeUrl))
	xhs.hs.Route("/liebian/add_qrcode_url_list", xhs.httpWrap(xhs.addQRCodeUrlList))
	xhs.hs.Route("/liebian/get_qrcode_url", xhs.httpWrap(xhs.getQRCodeUrl))
	xhs.hs.Route("/liebian/get_all_qrcode_url", xhs.httpWrap(xhs.getAllQRCodeUrlInfoList))
	xhs.hs.Route("/liebian/get_active_qrcode_url", xhs.httpWrap(xhs.getActiveQRCodeUrlInfoList))
	xhs.hs.Route("/liebian/expired_qrcode_url", xhs.httpWrap(xhs.expiredQRCodeUrl))
	xhs.hs.Route("/liebian/create_group_setting", xhs.httpWrap(xhs.createGroupSetting))
	xhs.hs.Route("/liebian/create_robot_msg_setting", xhs.httpWrap(xhs.createRobotMsgSetting))
	xhs.hs.Route("/liebian/del_robot_msg_setting", xhs.httpWrap(xhs.delRobotMsgSetting))
	
	xhs.hs.Route("/robot/receive_msg", xhs.httpWrap(xhs.RobotReceiveMsg))
}

func (xhs *XHttpServer) initAutoCheckGroupList() {
	list, err := GetTypeRobotMsgSettingList()
	if err != nil {
		plog.Errorf("get group check group list error: %v", err)
		return
	}
	xhs.Lock()
	for _, v := range list {
		plog.Infof("get robot msg setting: %v", v)
		acg := NewAutoCheckGroup(v, xhs.rbExt, xhs.rul)
		go acg.Run()
		xhs.acMap[v.ID] = acg
	}
	xhs.Unlock()
}

func (xhs *XHttpServer) httpWrap(handler HttpHandler) func(rsp http.ResponseWriter, req *http.Request) {
	f := func(rsp http.ResponseWriter, req *http.Request) {
		logURL := req.URL.String()
		start := time.Now()
		defer func() {
			plog.Debugf("[XHttpServer][httpWrap] http: request url[%s] use_time[%v]", logURL, time.Now().Sub(start))
		}()
		obj, err := handler(rsp, req)
		// check err
	HAS_ERR:
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Access-Control-Allow-Methods", "POST")
		rsp.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")

		if err != nil {
			plog.Debugf("[XHttpServer][httpWrap] http: request url[%s] error: %v", logURL, err)
			code := 500
			errMsg := err.Error()
			if strings.Contains(errMsg, "Permission denied") || strings.Contains(errMsg, "ACL not found") {
				code = 403
			}
			rsp.WriteHeader(code)
			rsp.Write([]byte(errMsg))
			return
		}

		// return json object
		if obj != nil {
			var buf []byte
			buf, err = json.Marshal(obj)
			if err != nil {
				goto HAS_ERR
			}
			rsp.Header().Set("Content-Type", "application/json")
			rsp.Write(buf)
		}
	}
	return f
}

func (xhs *XHttpServer) Index(rsp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		rsp.WriteHeader(404)
		return
	}
	rsp.Write([]byte("REAL TECH"))
}

func (xhs *XHttpServer) decodeBody(req *http.Request, out interface{}, cb func(interface{}) error) error {
	var raw interface{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&raw); err != nil {
		return err
	}

	if cb != nil {
		if err := cb(raw); err != nil {
			return err
		}
	}

	return mapstructure.Decode(raw, out)
}

type ClientInfo struct {
	IP        string
	UserAgent string
	Referrer  string
}

func (xhs *XHttpServer) GetClientInfo(req *http.Request) *ClientInfo {
	return &ClientInfo{
		IP:        strings.Split(req.RemoteAddr, ":")[0],
		UserAgent: req.UserAgent(),
		Referrer:  req.Referer(),
	}
}
