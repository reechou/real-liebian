package controller

import (
	"strconv"
	"time"

	"github.com/reechou/real-liebian/utils"
)

type ControllerDB struct {
	db *utils.MysqlController
}

func NewControllerDB(cfg *utils.MysqlInfo) (*ControllerDB, error) {
	cdb := &ControllerDB{
		db: utils.NewMysqlController(),
	}
	err := cdb.db.InitMysql(cfg)
	if err != nil {
		plog.Errorf("Mysql init error: %v\n", err)
		return nil, err
	}

	return cdb, nil
}

func (cdb *ControllerDB) InsertQRCodeUrl(info *QRCodeUrl) error {
	now := time.Now().Unix()
	id, err := cdb.db.Insert("insert into qrcode_url(name,url,status,create_time) values(?,?,?,?)", info.Name, info.Url, info.Status, now)
	if err != nil {
		return err
	}
	info.ID = id
	return nil
}

func (cdb *ControllerDB) GetQRCodeUrlList(startTime int64) (*QRCodeUrlList, error) {
	rows, err := cdb.db.FetchRows("select id,name,url,status,create_time,UNIX_TIMESTAMP(time) as utime from qrcode_url where UNIX_TIMESTAMP(time)>?", startTime)
	if err != nil {
		return nil, err
	}
	list := &QRCodeUrlList{
		UpdateTime: startTime,
	}
	for _, v := range *rows {
		id, err := strconv.ParseInt(v["id"], 10, 0)
		if err != nil {
			continue
		}
		status, err := strconv.ParseInt(v["status"], 10, 0)
		if err != nil {
			continue
		}
		createTime, err := strconv.ParseInt(v["create_time"], 10, 0)
		if err != nil {
			continue
		}
		uTime, err := strconv.ParseInt(v["utime"], 10, 0)
		if err != nil {
			continue
		}
		if list.UpdateTime < uTime {
			list.UpdateTime = uTime
		}
		info := &QRCodeUrl{
			ID:         id,
			Name:       v["name"],
			Url:        v["url"],
			Status:     status,
			CreateTime: createTime,
		}
		list.List = append(list.List, info)
	}
	return list, nil
}

func (cdb *ControllerDB) UpdateQRCodeStatus(info *QRCodeUrl) error {
	var err error
	if info.ID != 0 {
		_, err = cdb.db.Exec("update qrcode_url set status=? where id=?", info.Status, info.ID)
	} else if info.Url != "" {
		_, err = cdb.db.Exec("update qrcode_url set status=? where url=?", info.Status, info.Url)
	}
	if err != nil {
		return err
	}
	return nil
}
