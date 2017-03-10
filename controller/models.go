package controller

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/real-liebian/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.MysqlInfo.User,
		cfg.MysqlInfo.Pass,
		cfg.MysqlInfo.Host,
		cfg.MysqlInfo.DBName))
	if err != nil {
		plog.Fatalf("Fail to init new engine: %v", err)
	}
	//x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	// if need show raw sql in log
	x.ShowSQL(true)

	// sync tables
	if err = x.Sync2(new(QRCodeUrlInfo),
		new(UserQRCodeUrl),
		new(TypeGroupSetting),
		new(TypeRobotMsgSetting),
		new(QRCodeUrlRobot),
		new(RobotGroupMonitor)); err != nil {
		plog.Fatalf("Fail to sync database: %v", err)
	}
}
