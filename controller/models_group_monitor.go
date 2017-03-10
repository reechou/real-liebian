package controller

import (
	"time"
)

type RobotGroupMonitor struct {
	ID        int64 `xorm:"pk autoincr" json:"id"`
	Type      int64 `xorm:"not null default 0 int unique(uni_type_time)" json:"type"` // 渠道
	TimePoint int64 `xorm:"not null default 0 int unique(uni_type_time)" json:"timePoint"`
	GroupNum  int64 `xorm:"not null default 0 int" json:"groupNum"` // 群数量
	ShareNum  int64 `xorm:"not null default 0 int" json:"shareNum"` // 分享截图数量
	CreatedAt int64 `xorm:"not null default 0 int" json:"createdAt"`
}

func CreateRobotGroupMonitor(info *RobotGroupMonitor) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	
	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create robot group monitor error: %v", err)
		return err
	}
	plog.Infof("create obot group monitor[%v] success.", info)
	
	return nil
}

func GetRobotGroupMonitorFromTime(t, startTime, endTime int64) ([]RobotGroupMonitor, error) {
	var list []RobotGroupMonitor
	err := x.Where("type = ?", t).And("time_point >= ?", startTime).And("time_point <= ?", endTime).Find(&list)
	if err != nil {
		plog.Errorf("get robot group monitor error: %v", err)
		return nil, err
	}
	return list, nil
}
