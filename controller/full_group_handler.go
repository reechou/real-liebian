package controller

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

// maybe control group 30 minute, and release
// setting with after
const (
	SETTING_FULL_GROUP_START = 100
)

// setting with timer notify
const (
	SETTING_FULL_GROUP_IMG_NOTIFY = 200
)

const (
	ControlGroupAllTime = 1800
)

type GroupFullHandler struct {
	qrCodeInfo *QRCodeUrlInfo
	rul        *RobotUserLogic
	robotExt   *RobotExt

	fgm FullGroupManagerInterface

	settingDoneMap    map[int64]int
	settingList       []TypeRobotMsgSetting
	notifySettingList []TypeRobotMsgSetting
	robotList         []QRCodeUrlRobot

	startTime  int64
	checkTimes int
	idx        int

	stop chan struct{}
	done chan struct{}
}

func NewGroupFullHandler(qrCodeInfo *QRCodeUrlInfo, robotExt *RobotExt, rul *RobotUserLogic, fgm FullGroupManagerInterface) *GroupFullHandler {
	gfh := &GroupFullHandler{
		qrCodeInfo:     qrCodeInfo,
		robotExt:       robotExt,
		rul:            rul,
		fgm:            fgm,
		settingDoneMap: make(map[int64]int),
		startTime:      time.Now().Unix(),
		stop:           make(chan struct{}),
		done:           make(chan struct{}),
	}
	ok := gfh.init()
	if !ok {
		return nil
	}
	go gfh.run()

	return gfh
}

func (self *GroupFullHandler) Stop() {
	close(self.stop)
	<-self.done
}

func (self *GroupFullHandler) AddGroupImgUser(user string) {
	self.rul.AddGroupImgUser(self.qrCodeInfo.ID, "@"+user)
}

func (self *GroupFullHandler) init() bool {
	list, err := GetTypeRobotMsgSettingListOfEnd(self.qrCodeInfo.Type)
	if err != nil {
		plog.Errorf("get type[%d] robot msg setting of end error: %v", self.qrCodeInfo.Type, err)
		return false
	}
	notifyList, err := GetTypeRobotMsgSettingListOfEndNotify(self.qrCodeInfo.Type)
	if err != nil {
		plog.Errorf("get type[%d] robot msg notify setting of end error: %v", self.qrCodeInfo.Type, err)
		return false
	}
	if (list == nil || len(list) == 0) && (notifyList == nil || len(notifyList) == 0) {
		plog.Debugf("type[%d] has no setting of end.", self.qrCodeInfo.Type)
		return false
	}
	robotList, err := GetQRCodeUrlRobotList(self.qrCodeInfo.ID)
	if err != nil {
		plog.Errorf("qrcodeid[%d] get robot list end.", self.qrCodeInfo.ID)
		return false
	}
	if len(robotList) == 0 {
		plog.Debugf("qrcodeid[%d] has no robot list.", self.qrCodeInfo.ID)
		return false
	}
	self.settingList = list
	self.notifySettingList = notifyList
	self.robotList = robotList
	plog.Debugf("full group[%d][%s] robot list: %v", self.qrCodeInfo.ID, self.qrCodeInfo.Name, self.robotList)
	plog.Debugf("type[%d] name[%s] init start setting: %v %v", self.qrCodeInfo.Type, self.qrCodeInfo.Name, self.settingList, self.notifySettingList)

	return true
}

func (self *GroupFullHandler) end() {
	plog.Debugf("full group qrcode[%v] end.", self.qrCodeInfo)
	self.rul.DelGroup(self.qrCodeInfo.ID)
	self.fgm.ControlEnd(self.qrCodeInfo.ID)
}

func (self *GroupFullHandler) run() {
	defer func() {
		self.end()
	}()

	for {
		select {
		case <-time.After(31 * time.Second):
			ok := self.check()
			if ok {
				plog.Debugf("qrcodeinfo[%v] group full control run end.", self.qrCodeInfo)
				return
			}
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *GroupFullHandler) check() bool {
	ifEnd := true
	timeline := time.Now().Unix() - self.startTime
	plog.Debugf("full group[%s] check starttime[%d] timeline[%d]", self.qrCodeInfo.Name, self.startTime, timeline)
	if timeline > ControlGroupAllTime {
		return ifEnd
	}
	for _, v := range self.settingList {
		_, ok := self.settingDoneMap[v.ID]
		if ok {
			continue
		}
		if timeline >= v.After {
			self.settingDoneMap[v.ID] = 1
			msgs := make([]MsgInfo, 0)
			err := json.Unmarshal([]byte(v.Msg), &msgs)
			if err != nil {
				plog.Errorf("get setting msgs error: %v", err)
				continue
			}
			self.sendMsgs(msgs)
		} else {
			ifEnd = false
		}
	}
	if timeline < ControlGroupAllTime {
		ifEnd = false
		for _, v := range self.notifySettingList {
			switch v.SettingType {
			case SETTING_FULL_GROUP_IMG_NOTIFY:
				userList := self.rul.GetGroup(self.qrCodeInfo.ID)
				if userList != nil {
					msgs := make([]MsgInfo, 0)
					err := json.Unmarshal([]byte(v.Msg), &msgs)
					if err != nil {
						plog.Errorf("get notify setting msgs error: %v", err)
						continue
					}
					userStr := strings.Join(userList, " ")
					self.sendMsgsAddPrefix(userStr, msgs)
					self.rul.ClearGroup(self.qrCodeInfo.ID)
				}
			}
		}
	}
	return ifEnd
}

func (self *GroupFullHandler) sendMsgs(msgs []MsgInfo) {
	listIdx := self.idx % len(self.robotList)
	self.idx++
	if self.idx == 100 {
		self.idx = 0
	}
	robot := self.robotList[listIdx].RobotWx
	//robot := self.qrCodeInfo.RobotWx
	for _, v := range msgs {
		if v.MsgType == MSG_TYPE_TEXT {
			offset := rand.Intn(len(RANDOM_MSG_ADD))
			plog.Debugf("sendMsgs get random msg add offset: %d", offset)
			v.Msg = v.Msg + RANDOM_MSG_ADD[offset]
		}
		if v.MsgType == MSG_TYPE_IMG {
			urls := strings.Split(v.Msg, "|||")
			offset := rand.Intn(len(urls))
			v.Msg = urls[offset]
			plog.Debugf("sendMsgs get random img url offset: %d, %s", offset, v.Msg)
		}
		var sendReq SendMsgInfo
		sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
			WechatNick: robot,
			ChatType:   CHAT_TYPE_GROUP,
			UserName:   self.robotList[listIdx].UserName,
			NickName:   self.qrCodeInfo.Name,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
		self.robotExt.SendMsgs(self.qrCodeInfo.Type, robot, &sendReq)
		time.Sleep(2 * time.Second)
	}
}

func (self *GroupFullHandler) sendMsgsAddPrefix(prefix string, msgs []MsgInfo) {
	for _, v := range self.robotList {
		prefix = strings.Replace(prefix, "@"+v.RobotWx, "", -1)
	}
	if !strings.Contains(prefix, "@") {
		return
	}
	
	listIdx := self.idx % len(self.robotList)
	self.idx++
	if self.idx == 100 {
		self.idx = 0
	}
	robot := self.robotList[listIdx].RobotWx
	//robot := self.qrCodeInfo.RobotWx
	var sendReq SendMsgInfo
	for _, v := range msgs {
		if v.MsgType == MSG_TYPE_TEXT {
			//rand.Seed(time.Now().UnixNano())
			offset := rand.Intn(len(RANDOM_MSG_ADD))
			plog.Debugf("sendMsgsAddPrefix get random msg add offset: %d", offset)
			v.Msg = v.Msg + RANDOM_MSG_ADD[offset]
			sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
				WechatNick: robot,
				ChatType:   CHAT_TYPE_GROUP,
				UserName:   self.robotList[listIdx].UserName,
				NickName:   self.qrCodeInfo.Name,
				MsgType:    v.MsgType,
				Msg:        prefix + " " + v.Msg,
			})
		}
	}
	self.robotExt.SendMsgs(self.qrCodeInfo.Type, robot, &sendReq)
}
