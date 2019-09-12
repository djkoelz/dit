package main

import (
	"fmt"
	"github.com/djkoelz/dit/pkg/router"
	"github.com/djkoelz/dit/pkg/service"
)

func main() {
	fmt.Println("chk")
	router := router.NewRouter()
	router.Register("/push", service.PushImage)
	router.Register("/pull", service.PullImage)
	router.Register("/remove", service.RemoveImage)

	router.Start(6000)
}
