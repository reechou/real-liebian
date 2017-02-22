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

func (self *RobotExt) SendMsgs(robotWx string, msg *SendMsgInfo) error {
	plog.Debugf("robot[%s] send msg: %v", robotWx, msg)
	
	reqBytes, err := json.Marshal(msg)
	if err != nil {
		plog.Errorf("json encode error: %v", err)
		return err
	}
	
	url := "http://" + self.cfg.WxRobotExt.Host + "/sendmsgs"
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
