package controller

import (
	"time"
)

var (
	EVERY_RECORD_TIME int64 = 600
)

type WxGroupMonitor struct {
	RecordTime int64

	stop chan struct{}
	done chan struct{}
}

func NewWxGroupMonitor() *WxGroupMonitor {
	wgm := &WxGroupMonitor{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go wgm.run()

	return wgm
}

func (self *WxGroupMonitor) Stop() {
	close(self.stop)
	<-self.done
}

func (self *WxGroupMonitor) run() {
	plog.Debugf("wx group monitor run start.")
	for {
		select {
		case <-time.After(time.Minute):
			self.check()
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *WxGroupMonitor) check() {
	recordTime := self.calRecordTime()
	if self.RecordTime >= recordTime {
		return
	}
	typeList, err := GetQRCodeUrlRobotAllType()
	if err != nil {
		plog.Errorf("get group robot all type error: %v", err)
		return
	}
	plog.Debugf("Every stat time[%d] recordtime[%d] last recordtime[%d] typelist[%v].", EVERY_RECORD_TIME, recordTime, self.RecordTime, typeList)
	endTime := recordTime + EVERY_RECORD_TIME
	for _, v := range typeList {
		if v.Type == 0 {
			continue
		}
		groupList, err := GetQRCodeUrlInfoListFromTime(v.Type, recordTime, endTime)
		if err != nil {
			plog.Errorf("get qrcode info list form time: %v", err)
			return
		}
		var sharedNum int64
		for _, v2 := range groupList {
			sharedNum += v2.SharedNum
		}
		rgm := &RobotGroupMonitor{
			Type:      v.Type,
			TimePoint: recordTime,
			GroupNum:  int64(len(groupList)),
			ShareNum:  sharedNum,
		}
		CreateRobotGroupMonitor(rgm)
	}

	self.RecordTime = recordTime
}

func (self *WxGroupMonitor) calRecordTime() int64 {
	now := time.Now().Unix()
	now = now - EVERY_RECORD_TIME
	return now - (now % EVERY_RECORD_TIME)
}
