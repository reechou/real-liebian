package controller

import (
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"io/ioutil"
	
	"github.com/reechou/real-liebian/config"
)

type RobotExt struct {
	cfg *config.Config
	client *http.Client
}

func NewRobotExt(cfg *config.Config) *RobotExt {
	return &RobotExt{
		client: &http.Client{},
		cfg: cfg,
	}
}

func (self *RobotExt) SendMsgs(t int64, robotWx string, msg *SendMsgInfo) error {
	host := self.getRobotHost(t)
	plog.Debugf("robot[%s] host[%s] send msg: %v", robotWx, msg, host)
	
	reqBytes, err := json.Marshal(msg)
	if err != nil {
		plog.Errorf("json encode error: %v", err)
		return err
	}
	
	url := "http://" + host + "/sendmsgs"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		plog.Errorf("http new request error: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		plog.Errorf("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		plog.Errorf("ioutil ReadAll error: %v", err)
		return err
	}
	var response SendMsgResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		plog.Errorf("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		plog.Errorf("send msg[%v] result code error: %d %s", msg, response.Code, response.Msg)
		return fmt.Errorf("send msg result error.")
	}
	
	return nil
}

func (self *RobotExt) getRobotHost(t int64) string {
	tgs := &TypeGroupSetting{
		Type: t,
	}
	has, err := GetTypeGroupSetting(tgs)
	if err != nil {
		plog.Errorf("get type group setting error: %v", err)
		return self.cfg.WxRobotExt.Host
	}
	if !has {
		plog.Errorf("has none this type[%d] of group setting", t)
		return self.cfg.WxRobotExt.Host
	}
	if tgs.RobotHost == "" {
		plog.Debugf("not found the type[%d] of robot host", t)
		return self.cfg.WxRobotExt.Host
	}
	return tgs.RobotHost
}
