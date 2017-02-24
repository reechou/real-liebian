package controller

import (
	"testing"
)

func TestUser(t *testing.T) {
	rul := NewRobotUserLogic()
	rul.AddGroupImgUser(0, `浅浅<span class="emoji emoji1f601"></span><span class="emoji emoji1f602"></span>`)
}
