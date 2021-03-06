package controller

import (
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/coreos/pkg/capnslog"
	"github.com/reechou/real-liebian/config"
)

var plog = capnslog.NewPackageLogger("github.com/reezhou/real-liebian", "controller")

type ControllerLogic struct {
	sync.Mutex

	cfg *config.Config

	aliyunOss *config.AliyunOss
	//cdb       *ControllerDB
	xServer *XHttpServer
	rul     *RobotUserLogic
	monitor *WxGroupMonitor

	qrCodeUpateTime int64
	qrCodeUrlIdx    int

	stop chan struct{}
	done chan struct{}
}

func NewControllerLogic(cfg *config.Config) *ControllerLogic {
	cl := &ControllerLogic{
		cfg:       cfg,
		aliyunOss: &cfg.AliyunOss,
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
	}
	aliyunClient, err := oss.New(cl.aliyunOss.Endpoint, cl.aliyunOss.AccessKeyId, cl.aliyunOss.AccessKeySecret)
	if err != nil {
		plog.Panicf("aliyun oss new error: %v\n", err)
	}
	cl.aliyunOss.AliyunClient = aliyunClient
	//db, err := NewControllerDB(&cfg.MysqlInfo)
	//if err != nil {
	//	plog.Panicf("db controller new error: %v\n", err)
	//}
	//cl.cdb = db
	//err = cl.Init()
	//if err != nil {
	//	plog.Panicf("logic init error: %v\n", err)
	//}
	//go cl.run()

	InitDB(cfg)
	cl.rul = NewRobotUserLogic()
	cl.monitor = NewWxGroupMonitor()
	cl.xServer = NewXHttpServer(cfg.ListenAddr, cfg.ListenPort, cl, cl.rul)
	setupLogging(cfg)

	return cl
}

func (cl *ControllerLogic) Start() {
	cl.xServer.Run()
}

func (cl *ControllerLogic) Stop() {
	close(cl.stop)
	<-cl.done
}

//func (cl *ControllerLogic) Init() error {
//	list, err := cl.cdb.GetQRCodeUrlList(0)
//	if err != nil {
//		return err
//	}
//	now := time.Now().Unix()
//	cl.qrCodeUpateTime = list.UpdateTime
//	for _, v := range list.List {
//		if v.Status != QRCODE_STATUS_OK {
//			continue
//		}
//		if now-v.CreateTime >= cl.cfg.QRCodeExpired {
//			continue
//		}
//		cl.qrCodeUrlMap[v.ID] = v
//	}
//
//	return nil
//}

//func (cl *ControllerLogic) run() {
//	for {
//		select {
//		case <-time.After(30 * time.Second):
//			cl.onRefresh()
//		case <-cl.stop:
//			close(cl.done)
//			return
//		}
//	}
//}

//func (cl *ControllerLogic) onRefresh() {
//	list, err := cl.cdb.GetQRCodeUrlList(cl.qrCodeUpateTime)
//	if err != nil {
//		plog.Errorf("get qrcodeurl list error: %v", err)
//		return
//	}
//
//	cl.Lock()
//	defer cl.Unlock()
//	cl.qrCodeUpateTime = list.UpdateTime
//	now := time.Now().Unix()
//	for _, v := range list.List {
//		vqr := cl.qrCodeUrlMap[v.ID]
//		if vqr != nil {
//			if v.Status != QRCODE_STATUS_OK || now-v.CreateTime >= cl.cfg.QRCodeExpired {
//				delete(cl.qrCodeUrlMap, v.ID)
//				continue
//			}
//			vqr.Name = v.Name
//			vqr.Url = v.Url
//			continue
//		}
//		if v.Status != QRCODE_STATUS_OK || now-v.CreateTime >= cl.cfg.QRCodeExpired {
//			continue
//		}
//		cl.qrCodeUrlMap[v.ID] = v
//	}
//	plog.Infof("on refresh get qrcode url success.\n")
//}

//func (cl *ControllerLogic) GetQRCodeUrl() (*QRCodeUrl, error) {
//	if len(cl.qrCodeUrlMap) == 0 {
//		return nil, fmt.Errorf("cannot find useful qrcode url.")
//	}
//	plog.Debugf("get qrcode url: all code: %v.\n", cl.qrCodeUrlMap)
//	// rand.Seed(time.Now().Unix())
//	// qIdx := rand.Intn(len(cl.qrCodeUrlMap))
//	// addIdx := 0
//	// if qIdx >= len(cl.qrCodeUrlMap) {
//	// 	qIdx = 0
//	// }
//
//	cl.Lock()
//	defer cl.Unlock()
//	now := time.Now().Unix()
//	for _, v := range cl.qrCodeUrlMap {
//		plog.Debugf("v(%v) now(%d) CreateTime(%d)", v, now, v.CreateTime)
//		if now-v.CreateTime >= cl.cfg.QRCodeExpired {
//			continue
//		}
//		return v, nil
//	}
//	return nil, fmt.Errorf("cannot find useful qrcode url.")
//}

func setupLogging(cfg *config.Config) {
	capnslog.SetGlobalLogLevel(capnslog.INFO)
	if cfg.Debug {
		capnslog.SetGlobalLogLevel(capnslog.DEBUG)
	}
}
