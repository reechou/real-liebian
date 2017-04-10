package controller

import (
	"fmt"
	"math/rand"
	"time"
)

// type: 渠道号 0-REAL群二维码 1-达人二维码 2-普通用户二维码
type QRCodeUrlInfo struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	Name      string `xorm:"not null default '' varchar(128) index" json:"name"`
	Url       string `xorm:"not null default '' varchar(256)" json:"url"`
	Type      int64  `xorm:"not null default 0 int index" json:"type"`
	UserName  string `xorm:"not null default '' varchar(128) index" json:"userName"`
	IfMod     int64  `xorm:"not null default 0 int" json:"ifMod"`
	RobotWx   string `xorm:"not null default '' varchar(128) index" json:"robotWx"`
	GroupNum  int64  `xorm:"not null default 0 int" json:"groupNum"`
	Status    int64  `xorm:"not null default 0 int index" json:"status"`
	EndTime   int64  `xorm:"not null default 0 int index" json:"endTime"`
	SharedNum int64  `xorm:"not null default 0 int" json:"sharedNum"`
	CreatedAt int64  `xorm:"not null default 0 int index" json:"createdAt"`
}

func CreateQRCodeUrlInfo(info *QRCodeUrlInfo) error {
	if info.Url == "" {
		return fmt.Errorf("url cannot be nil.")
	}

	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create qrcode url error: %v", err)
		return err
	}
	plog.Infof("create qrcode url[%s] type[%d] success.", info.Url, info.Type)

	return nil
}

func CreateQRCodeUrlInfoList(list []QRCodeUrlInfo) error {
	if len(list) == 0 {
		return nil
	}

	now := time.Now().Unix()
	for i := 0; i < len(list); i++ {
		list[i].CreatedAt = now
	}

	_, err := x.Insert(&list)
	if err != nil {
		plog.Errorf("create qrcode url info list error: %v", err)
		return err
	}

	return nil
}

func GetAllQRCodeUrlInfoFromTypeCount(t int64) (int64, error) {
	count, err := x.Where("type = ?", t).And("status = 0").Count(&QRCodeUrlInfo{})
	if err != nil {
		plog.Errorf("get all qrcode url info from type count error: %v", err)
		return 0, err
	}
	return count, nil
}

func GetAllQRCodeUrlInfoFromType(t, offset, num int64) ([]QRCodeUrlInfo, error) {
	var list []QRCodeUrlInfo
	err := x.Where("type = ?", t).And("status = 0").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		plog.Errorf("get all qrcode url info from type list error: %v", err)
		return nil, err
	}
	return list, nil
}

func GetQRCodeUrlInfoFromId(info *QRCodeUrlInfo) (bool, error) {
	has, err := x.Where("id = ?", info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url info from id[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetQRCodeUrlInfo(info *QRCodeUrlInfo) (bool, error) {
	has, err := x.Where("name = ?", info.Name).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url info from info[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetQRCodeUrlInfoActive(info *QRCodeUrlInfo) (bool, error) {
	has, err := x.Where("status = 0").And("name = ?", info.Name).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url info from info[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetQRCodeUrlInfoFromRobotUserName(info *QRCodeUrlInfo) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("user_name = ?", info.UserName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find qrcode url info from robot username info[%v]", info)
		return false, nil
	}
	return true, nil
}

func GetQRCodeUrlInfoListCount(t int64) (int64, error) {
	count, err := x.Where("type = ?", t).Count(&QRCodeUrlInfo{})
	if err != nil {
		plog.Errorf("type[%d] get qrcode list count error: %v", t, err)
		return 0, err
	}
	return count, nil
}

func GetQRCodeUrlInfoList(num, t int64) ([]QRCodeUrlInfo, error) {
	createdTime := time.Now().Unix() - 518400
	var list []QRCodeUrlInfo
	err := x.Where("type = ?", t).And("status = 0").And("created_at > ?", createdTime).Limit(int(num)).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetQRCodeUrlInfoListWithoutExpired(num, t int64) ([]QRCodeUrlInfo, error) {
	var list []QRCodeUrlInfo
	err := x.Where("type = ?", t).And("status = 0").Limit(int(num)).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetQRCodeUrlInfoListFromTime(t, startTime, endTime int64) ([]QRCodeUrlInfo, error) {
	var list []QRCodeUrlInfo
	err := x.Where("type = ?", t).And("end_time >= ?", startTime).And("end_time <= ?", endTime).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetQRCodeUrlInfoOfRandom(count, t int64) (*QRCodeUrlInfo, error) {
	rand.Seed(time.Now().UnixNano())
	offset := rand.Intn(int(count))
	createdTime := time.Now().Unix() - 518400
	var list []QRCodeUrlInfo
	err := x.Where("type = ?", t).And("created_at > ?", createdTime).Limit(1, int(offset)).Find(&list)
	if err != nil {
		plog.Errorf("type[%d] get qr code info of random error: %v", t, err)
		return nil, err
	}
	if len(list) > 0 {
		return &list[0], nil
	}
	return nil, nil
}

func UpdateQRCodeUrlInfoStatus(info *QRCodeUrlInfo) error {
	_, err := x.ID(info.ID).Cols("status").Update(info)
	return err
}

func UpdateQRCodeUrlInfoSharedNum(info *QRCodeUrlInfo) error {
	_, err := x.ID(info.ID).Cols("shared_num").Update(info)
	return err
}

func UpdateQRCodeUrlInfoStatusFromName(info *QRCodeUrlInfo) error {
	_, err := x.Cols("status", "end_time").Update(info, &QRCodeUrlInfo{Name: info.Name})
	if err != nil {
		plog.Errorf("update qrcode url info from name error: %v", err)
	}
	return err
}

func UpdateQRCodeUrlInfoIfModAndGroupNum(info *QRCodeUrlInfo) error {
	_, err := x.ID(info.ID).Cols("if_mod", "group_num").Update(info)
	return err
}

func UpdateQRCodeUrlInfoIfMod(info *QRCodeUrlInfo) error {
	_, err := x.ID(info.ID).Cols("if_mod").Update(info)
	return err
}

func UpdateQRCodeUrlInfoRobotWx(info *QRCodeUrlInfo) error {
	_, err := x.ID(info.ID).Cols("user_name", "robot_wx").Update(info)
	return err
}
