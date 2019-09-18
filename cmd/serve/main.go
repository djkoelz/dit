package main

import (
	"github.com/djkoelz/dit/pkg/repo"
	"github.com/djkoelz/dit/pkg/router"
)

func main() {
	store := repo.NewStore("localhost:5000")
	service := repo.NewService(store)

	router := router.NewRouter()
	router.Register("/add", service.AddImage)
	router.Register("/get", service.GetImage)
	router.Register("/remove", service.RemoveImage)
	router.Register("/list", service.ListImages)
	router.Register("/sync", service.SyncImages)

	router.Start(6000)
}
