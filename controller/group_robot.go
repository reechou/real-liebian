package controller

import (
	"time"
)

type QRCodeUrlRobot struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	QrcodeId  int64  `xorm:"not null default 0 int index" json:"qrcodeId"`
	GroupName string `xorm:"not null default '' varchar(128)" json:"groupName"`
	UserName  string `xorm:"not null default '' varchar(128) unique(robot)" json:"userName"`
	RobotWx   string `xorm:"not null default '' varchar(128) unique(robot)" json:"robotWx"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"updatedAt"`
}

func CreateQRCodeUrlRobot(info *QRCodeUrlRobot) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create qrcode url robot error: %v", err)
		return err
	}
	plog.Infof("create qrcode[%d] robot[%s] success.", info.QrcodeId, info.RobotWx)

	return nil
}

func GetQRCodeUrlRobotFromRobot(info *QRCodeUrlRobot) (bool, error) {
	has, err := x.Where("user_name = ?", info.UserName).And("robot_wx = ?", info.RobotWx).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url robot from robot[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetQRCodeUrlRobot(info *QRCodeUrlRobot) (bool, error) {
	has, err := x.Where("qrcode_id = ?", info.QrcodeId).And("robot_wx = ?", info.RobotWx).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url robot from qrcode_id[%v]", info)
		return false, nil
	}
	return true, nil
}

func UpdateQRCodeUrlRobot(info *QRCodeUrlRobot) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("user_name", "robot_wx", "updated_at").Update(info)
	return err
}

func GetQRCodeUrlRobotList(qrcodeId int64) ([]QRCodeUrlRobot, error) {
	var list []QRCodeUrlRobot
	err := x.Where("qrcode_id = ?", qrcodeId).Find(&list)
	if err != nil {
		plog.Errorf("get qrcodeid[%d] robot list error: %v", qrcodeId, err)
		return nil, err
	}
	return list, nil
}
