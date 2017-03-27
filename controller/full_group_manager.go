package controller

import (
	"sync"
	"time"
)

type FullGroupManagerInterface interface {
	ControlEnd(qrcodeId int64)
}

type FullGroupManager struct {
	sync.Mutex
	robotExt *RobotExt
	rul      *RobotUserLogic

	fghMap map[int64]*GroupFullHandler
	msgs   chan ReceiveMsgInfo

	stop chan struct{}
	done chan struct{}
}

func NewFullGroupManager(robotExt *RobotExt, rul *RobotUserLogic) *FullGroupManager {
	fgm := &FullGroupManager{
		robotExt: robotExt,
		rul:      rul,
		fghMap:   make(map[int64]*GroupFullHandler),
		msgs:     make(chan ReceiveMsgInfo, 1024),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	go fgm.run()

	return fgm
}

func (self *FullGroupManager) run() {
	plog.Infof("full group manager start run.")
	for {
		select {
		case msg := <-self.msgs:
			self.handleMsg(&msg)
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *FullGroupManager) handleMsg(msg *ReceiveMsgInfo) {
	qrcodeUrlInfo := &QRCodeUrlInfo{
		Name: msg.BaseInfo.FromGroupName,
	}
	has, err := GetQRCodeUrlInfo(qrcodeUrlInfo)
	if err != nil {
		plog.Errorf("get qrcode info error: %v", err)
		return
	}
	if !has {
		return
	}
	//if qrcodeUrlInfo.RobotWx != msg.BaseInfo.WechatNick {
	//	qrcodeUrlInfo.RobotWx = msg.BaseInfo.WechatNick
	//	UpdateQRCodeUrlInfoRobotWx(qrcodeUrlInfo)
	//}
	self.Lock()
	defer self.Unlock()
	fgh, ok := self.fghMap[qrcodeUrlInfo.ID]
	if ok {
		switch msg.BaseInfo.ReceiveEvent {
		case RECEIVE_EVENT_MSG:
			if msg.MsgType == MSG_TYPE_IMG {
				fgh.AddGroupImgUser(msg.BaseInfo.FromNickName)
			}
		}
	}
}

func (self *FullGroupManager) Stop() {
	for _, v := range self.fghMap {
		v.Stop()
	}

	close(self.stop)
	<-self.done
}

func (self *FullGroupManager) FilterReceiveMsg(msg ReceiveMsgInfo) {
	select {
	case self.msgs <- msg:
	case <-time.After(2 * time.Second):
		plog.Errorf("filter receive msg maybe channal is full.")
		return
	}
}

func (self *FullGroupManager) ControlEnd(qrcodeId int64) {
	self.Lock()
	defer self.Unlock()

	delete(self.fghMap, qrcodeId)
}

func (self *FullGroupManager) GroupFull(qrcodeInfo *QRCodeUrlInfo) {
	self.Lock()
	defer self.Unlock()

	_, ok := self.fghMap[qrcodeInfo.ID]
	if ok {
		return
	}

	fgh := NewGroupFullHandler(qrcodeInfo, self.robotExt, self.rul, self)
	if fgh == nil {
		self.rul.DelGroup(qrcodeInfo.ID)
		return
	}

	self.fghMap[qrcodeInfo.ID] = fgh
}
