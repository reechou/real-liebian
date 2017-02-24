package controller

import (
	"time"
)

type UserQRCodeUrl struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	AppId     string `xorm:"not null default '' varchar(128) unique(account_qrcode_url)" json:"appId"`
	OpenId    string `xorm:"not null default '' varchar(128) unique(account_qrcode_url)" json:"openId"`
	Type      int64  `xorm:"not null default 0 int unique(account_qrcode_url)" json:"type"`
	Url       string  `xorm:"not null default '' varchar(256)" json:"url"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createdAt"`
}

func CreateUserQRCodeUrl(info *UserQRCodeUrl) error {
	now := time.Now().Unix()
	info.CreatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		plog.Errorf("create user qrcode url error: %v", err)
		return err
	}
	plog.Infof("create user qrcode url[%v] success.", info)

	return nil
}

func GetUserQRCodeUrl(info *UserQRCodeUrl) (bool, error) {
	has, err := x.Where("app_id = ?", info.AppId).And("open_id = ?", info.OpenId).And("type = ?", info.Type).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		plog.Debugf("cannot find user qrcode url from info[%v]", info)
		return false, nil
	}
	return true, nil
}
