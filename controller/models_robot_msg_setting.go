package controller

import (
	"time"
)

type MsgInfo struct {
	MsgType string `json:"msgType"`
	Msg     string `json:"msg"`
}
type TypeGroupSetting struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	Type         int64  `xorm:"not null default 0 int unique" json:"type"`
	GroupNum     int64  `xorm:"not null default 0 int" json:"groupNum"`
	RobotHost    string `xorm:"not null default '' varchar(128)" json:"robotHost"`
	Desc         string `xorm:"not null default '' varchar(128)" json:"desc"`
	IfHasExpired int    `xorm:"not null default 0 int" json:"ifHasExpired"`
	CreatedAt    int64  `xorm:"not null default 0 int" json:"createdAt"`
}
type TypeRobotMsgSetting struct {
	ID          int64  `xorm:"pk autoincr" json:"id"`
	Type        int64  `xorm:"not null default 0 int index" json:"type"`
	SettingType int    `xorm:"not null default 0 int index" json:"settingType"`
	Robot       string `xorm:"not null default '' varchar(128)" json:"robot"`
	Msg         string `xorm:"not null default '' varchar(2048)" json:"msg"`
	Interval    int64  `xorm:"not null default 0 int" json:"interval"`
	After       int64  `xorm:"not null default 0 int" json:"after"`
	CreatedAt   int64  `xorm:"not null default 0 int" json:"createdAt"`
}

func CreateTypeGroupSetting(info *TypeGroupSetting) error {
	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create type group setting setting error: %v", err)
		return err
	}
	plog.Infof("create type group setting[%v] success.", info)

	return nil
}

func GetTypeGroupSetting(info *TypeGroupSetting) (bool, error) {
	has, err := x.Where("type = ?", info.Type).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find type group setting from info[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetTypeGroupSettingList() ([]TypeGroupSetting, error) {
	var list []TypeGroupSetting
	err := x.Find(&list)
	if err != nil {
		plog.Errorf("get robot group setting list error: %v", err)
		return nil, err
	}
	return list, nil
}

func UpdateTypeGroupSetting(info *TypeGroupSetting) error {
	_, err := x.ID(info.ID).Cols("group_num", "robot_host", "desc").Update(info)
	if err != nil {
		plog.Errorf("update type robot group setting error: %v", err)
	}
	return err
}

func GetTypeRobotMsgSettingAll(t int64) ([]TypeRobotMsgSetting, error) {
	var list []TypeRobotMsgSetting
	err := x.Where("type = ?", t).Find(&list)
	if err != nil {
		plog.Errorf("get robot msg setting list all from type error: %v", err)
		return nil, err
	}
	return list, nil
}

func UpdateTypeRobotMsgSetting(info *TypeRobotMsgSetting) error {
	_, err := x.ID(info.ID).Cols("msg", "interval", "after").Update(info)
	if err != nil {
		plog.Errorf("update type robot msg setting error: %v", err)
	}
	return err
}

func GetTypeRobotMsgSettingList() ([]TypeRobotMsgSetting, error) {
	var list []TypeRobotMsgSetting
	err := x.Where("setting_type < ?", SETTING_FULL_GROUP_START).Find(&list)
	if err != nil {
		plog.Errorf("get robot msg setting list error: %v", err)
		return nil, err
	}
	return list, nil
}

func GetTypeRobotMsgSettingListOfEnd(t int64) ([]TypeRobotMsgSetting, error) {
	var list []TypeRobotMsgSetting
	err := x.Where("type = ?", t).And("setting_type >= ?", SETTING_FULL_GROUP_START).And("setting_type < ?", SETTING_FULL_GROUP_IMG_NOTIFY).Find(&list)
	if err != nil {
		plog.Errorf("get robot msg setting list error: %v", err)
		return nil, err
	}
	return list, nil
}

func GetTypeRobotMsgSettingListOfEndNotify(t int64) ([]TypeRobotMsgSetting, error) {
	var list []TypeRobotMsgSetting
	err := x.Where("type = ?", t).And("setting_type >= ?", SETTING_FULL_GROUP_IMG_NOTIFY).Find(&list)
	if err != nil {
		plog.Errorf("get robot msg setting list error: %v", err)
		return nil, err
	}
	return list, nil
}

func CreateTypeRobotMsgSetting(info *TypeRobotMsgSetting) error {
	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create robot msg setting error: %v", err)
		return err
	}
	plog.Infof("create robot msg setting[%v] success.", info)

	return nil
}

func DelTypeRobotMsgSetting(id int64) error {
	info := &TypeRobotMsgSetting{ID: id}
	_, err := x.Where("id = ?", id).Delete(info)
	if err != nil {
		plog.Errorf("id[%d] wx circle delete error: %v", id, err)
		return err
	}
	return nil
}

func GetTypeRobotMsgSetting(info *TypeRobotMsgSetting) (bool, error) {
	has, err := x.Where("id = ?", info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find robot msg setting from info[%v]", info)
		return false, nil
	}
	return true, nil
}
