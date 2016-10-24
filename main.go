package main

import (
	"github.com/reechou/real-liebian/config"
	"github.com/reechou/real-liebian/controller"
)

func main() {
	controller.NewControllerLogic(config.NewConfig()).Start()
}
