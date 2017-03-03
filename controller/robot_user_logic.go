package controller

import (
	"regexp"
	"strings"
	"sync"
	"strconv"

	//"github.com/kyokomi/emoji"
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

	//reg := regexp.MustCompile(`<span class=\"emoji (.*?)\"><\/span>`)
	//regList := reg.FindAllString(user, -1)
	//for _, v := range regList {
	//	newEmoji := strings.Replace(v, "<span class=\"emoji emoji", "", -1)
	//	newEmoji = strings.Replace(newEmoji, "\"></span>", "", -1)
	//	//fmt.Println(v, emoji.Sprintf(emojiMap[newEmoji]))
	//	ev, ok := emojiMap[newEmoji]
	//	if ok {
	//		user = strings.Replace(user, v, ev, -1)
	//	}
	//}
	//fmt.Println(emoji.Sprintf(user))
	//emojiUser := emoji.Sprintf(user)
	
	emojiUser := replaceEmoji(user)
	userList := self.UserRobotImgMap[id]
	for _, v := range userList {
		if v == emojiUser {
			return
		}
	}
	//user = strings.Replace(user, "<span class=\"emoji ", "", -1)
	//user = strings.Replace(user, "\"></span>", "", -1)
	self.UserRobotImgMap[id] = append(userList, emojiUser)
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

func replaceEmoji(oriStr string) string {
	
	newStr := oriStr
	
	if strings.Contains(oriStr, `<span class="emoji`) {
		reg, _ := regexp.Compile(`<span class="emoji emoji[a-f0-9]{5}"></span>`)
		newStr = reg.ReplaceAllStringFunc(oriStr, func(arg2 string) string {
			num := `'\U000` + arg2[len(arg2)-14:len(arg2)-9] + `'`
			emoji, err := strconv.Unquote(num)
			if err == nil {
				return emoji
			}
			return num
		})
	}
	
	return newStr
}
