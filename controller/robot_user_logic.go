package controller

import (
	"sync"
	"strings"
)

type RobotUserLogic struct {
	sync.Mutex
	UserRobotImgMap map[int64][]string
}

func NewRobotUserLogic() *RobotUserLogic {
	rul := &RobotUserLogic{
		UserRobotImgMap: make(map[int64][]string),
	}
	
	return rul
}

func (self *RobotUserLogic) GetGroup(id int64) []string {
	self.Lock()
	defer self.Unlock()
	
	return self.UserRobotImgMap[id]
}

func (self *RobotUserLogic) AddGroupImgUser(id int64, user string) {
	self.Lock()
	defer self.Unlock()
	
	userList := self.UserRobotImgMap[id]
	user = strings.Replace(user, "<span class=\"emoji", " ", -1)
	user = strings.Replace(user, "\"></span>", " ", -1)
	self.UserRobotImgMap[id] = append(userList, user)
	plog.Debugf("add group img user[%d]: %v", id, self.UserRobotImgMap[id])
}

func (self *RobotUserLogic) DelGroup(id int64) {
	self.Lock()
	defer self.Unlock()
	
	delete(self.UserRobotImgMap, id)
}

func (self *RobotUserLogic) ClearGroup(id int64) {
	self.Lock()
	defer self.Unlock()
	
	self.UserRobotImgMap[id] = nil
}
