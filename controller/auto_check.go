package controller

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	SETTING_TYPE_GROUP_ADD = iota   // 群成员增加
	SETTING_TYPE_GROUP_USER_MSG_IMG // 群成员发送图片
	SETTING_TYPE_GROUP_OTHER        // 其他
)

type AutoCheckGroup struct {
	setting  TypeRobotMsgSetting
	msgs     []MsgInfo
	robotExt *RobotExt
	rul      *RobotUserLogic

	stop chan struct{}
	done chan struct{}
}

func NewAutoCheckGroup(setting TypeRobotMsgSetting, robotExt *RobotExt, rul *RobotUserLogic) *AutoCheckGroup {
	acg := &AutoCheckGroup{
		setting:  setting,
		robotExt: robotExt,
		rul:      rul,
		msgs:     make([]MsgInfo, 0),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	err := json.Unmarshal([]byte(acg.setting.Msg), &acg.msgs)
	if err != nil {
		plog.Errorf("get setting msgs error: %v", err)
	}

	return acg
}

func (self *AutoCheckGroup) Stop() {
	plog.Infof("auto group check setting: %v stopped.", self.setting)
	close(self.stop)
	<-self.done
}

func (self *AutoCheckGroup) Run() {
	plog.Infof("auto group check setting: %v start.", self.setting)
	plog.Infof("auto group check msg: %v.", self.msgs)
	for {
		select {
		case <-time.After(time.Duration(self.setting.Interval) * time.Second):
			self.check()
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *AutoCheckGroup) check() {
	activeList, err := GetQRCodeUrlListFromType(self.setting.Type)
	if err != nil {
		plog.Errorf("get qrcode url list from type error: %v", err)
		return
	}
	switch self.setting.SettingType {
	case SETTING_TYPE_GROUP_ADD:
		for _, v := range activeList {
			if v.IfMod != 0 {
				self.sendMsgs(&v)
				v.IfMod = 0
				UpdateQRCodeUrlInfoIfMod(&v)
			}
		}
	case SETTING_TYPE_GROUP_USER_MSG_IMG:
		for _, v := range activeList {
			userList := self.rul.GetGroup(v.ID)
			if userList != nil {
				userStr := strings.Join(userList, " ")
				self.sendMsgsAddPrefix(userStr, &v)
				self.rul.ClearGroup(v.ID)
			}
		}
	case SETTING_TYPE_GROUP_OTHER:
		for _, v := range activeList {
			self.sendMsgs(&v)
		}
	}
}

func (self *AutoCheckGroup) sendMsgs(info *QRCodeUrlInfo) {
	var sendReq SendMsgInfo
	for _, v := range self.msgs {
		sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
			WechatNick: self.setting.Robot,
			ChatType:   CHAT_TYPE_GROUP,
			NickName:   info.Name,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
	}
	self.robotExt.SendMsgs(self.setting.Robot, &sendReq)
}

func (self *AutoCheckGroup) sendMsgsAddPrefix(prefix string, info *QRCodeUrlInfo) {
	var sendReq SendMsgInfo
	for _, v := range self.msgs {
		sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
			WechatNick: self.setting.Robot,
			ChatType:   CHAT_TYPE_GROUP,
			NickName:   info.Name,
			MsgType:    v.MsgType,
			Msg:        prefix + " " + v.Msg,
		})
	}
	self.robotExt.SendMsgs(self.setting.Robot, &sendReq)
}
