package controller

import (
	"fmt"
	"math/rand"
	"time"
)

// type: 0-群二维码 1-达人二维码 2-普通用户二维码
type QRCodeUrlInfo struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	Name      string `xorm:"not null default '' varchar(128)" json:"name"`
	Url       string `xorm:"not null default '' varchar(256)" json:"url"`
	Type      int64  `xorm:"not null default 0 int index" json:"type"`
	Status    int64  `xorm:"not null default 0 int index" json:"status"`
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

func GetQRCodeUrlInfoListCount(t int64) (int64, error) {
	count, err := x.Where("type = ?", t).Count(&QRCodeUrlInfo{})
	if err != nil {
		plog.Errorf("type[%d] get qrcode list count error: %v", t, err)
		return 0, err
	}
	return count, nil
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
