package controller

import (
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	SETTING_TYPE_GROUP_ADD          = iota // 群成员增加
	SETTING_TYPE_GROUP_USER_MSG_IMG        // 群成员发送图片
	SETTING_TYPE_GROUP_OTHER               // 其他
)

var (
	RANDOM_MSG_ADD = []string{
		".", "..", "↭",
		"★", "✔", "↧",
		"↩", "⇤", "⇜",
		"↞", "↜", "┄",
		"-", "--", "^", "^_^",
		"!", "!!", "↮",
		"！", "•", "“",
		"[机智]", "[机智][机智]",
		"♥", "♥♥", "♥♥♥",
		"─", "↕↕", "↕",
		"☈", "✓", "☑",
		"⊰", "⊱", "†",
		"↓", "ˉ", "﹀",
		"﹏", "˜", "ˆ",
		"﹡", "≑", "≐",
		"≍", "≎", "≏",
		"≖", "≗", "≡",
	}
)

type AutoCheckGroup struct {
	sync.Mutex
	setting  TypeRobotMsgSetting
	msgs     []MsgInfo
	robotExt *RobotExt
	rul      *RobotUserLogic
	idx      int

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

func (self *AutoCheckGroup) Refresh(setting *TypeRobotMsgSetting) {
	self.Lock()
	defer self.Unlock()

	self.setting.Type = setting.Type
	self.setting.SettingType = setting.SettingType
	self.setting.Robot = setting.Robot
	self.setting.Msg = setting.Msg
	self.setting.Interval = setting.Interval

	self.msgs = make([]MsgInfo, 0)
	err := json.Unmarshal([]byte(self.setting.Msg), &self.msgs)
	if err != nil {
		plog.Errorf("get setting msgs error: %v", err)
	}

	plog.Infof("refresh auto group check setting: %v.", self.setting)
	plog.Infof("refresh auto group check msg: %v.", self.msgs)
}

func (self *AutoCheckGroup) check() {
	self.Lock()
	defer self.Unlock()

	activeList, err := GetQRCodeUrlListFromType(self.setting.Type)
	if err != nil {
		plog.Errorf("get qrcode url list from type error: %v", err)
		return
	}
	plog.Debugf("auto check active list: %v", activeList)
	switch self.setting.SettingType {
	case SETTING_TYPE_GROUP_ADD:
		for _, v := range activeList {
			if v.IfMod != 0 {
				v.IfMod = 0
				UpdateQRCodeUrlInfoIfMod(&v)
				self.sendMsgs(&v)
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
	list, err := GetQRCodeUrlRobotList(info.ID)
	if err != nil {
		plog.Errorf("get qrcode robot list error: %v", err)
		return
	}
	if len(list) == 0 {
		plog.Errorf("cannot found robot from id[%d]", info.ID)
		return
	}
	listIdx := self.idx % len(list)
	self.idx++
	if self.idx == 100 {
		self.idx = 0
	}
	robot := list[listIdx].RobotWx
	//if robot == "" {
	//	robot = self.setting.Robot
	//}
	for _, v := range self.msgs {
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
			UserName:   list[listIdx].UserName,
			NickName:   info.Name,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
		self.robotExt.SendMsgs(info.Type, robot, &sendReq)
		time.Sleep(2 * time.Second)
	}
}

func (self *AutoCheckGroup) sendMsgsAddPrefix(prefix string, info *QRCodeUrlInfo) {
	list, err := GetQRCodeUrlRobotList(info.ID)
	if err != nil {
		plog.Errorf("get qrcode robot list error: %v", err)
		return
	}
	if len(list) == 0 {
		plog.Errorf("cannot found robot from id[%d]", info.ID)
		return
	}
	for _, v := range list {
		prefix = strings.Replace(prefix, "@"+v.RobotWx, "", -1)
	}
	if !strings.Contains(prefix, "@") {
		return
	}

	listIdx := self.idx % len(list)
	self.idx++
	if self.idx == 100 {
		self.idx = 0
	}
	robot := list[listIdx].RobotWx
	//if robot == "" {
	//	robot = self.setting.Robot
	//}
	var sendReq SendMsgInfo
	for _, v := range self.msgs {
		if v.MsgType == MSG_TYPE_TEXT {
			//rand.Seed(time.Now().UnixNano())
			offset := rand.Intn(len(RANDOM_MSG_ADD))
			plog.Debugf("sendMsgsAddPrefix get random msg add offset: %d", offset)
			v.Msg = v.Msg + RANDOM_MSG_ADD[offset]
		}
		sendReq.SendMsgs = append(sendReq.SendMsgs, SendBaseInfo{
			WechatNick: robot,
			ChatType:   CHAT_TYPE_GROUP,
			UserName:   list[listIdx].UserName,
			NickName:   info.Name,
			MsgType:    v.MsgType,
			Msg:        prefix + " " + v.Msg,
		})
	}
	self.robotExt.SendMsgs(info.Type, robot, &sendReq)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
