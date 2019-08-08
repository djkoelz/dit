package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"log"
	"os"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatal(err)
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}

	// err = client.RemoveImage("hello-world")
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("Removed image")
	// }

	tar, err := os.Open("image.tar")
	if err != nil {
		log.Fatal(err)
	} else {
		defer tar.Close()
	}
	opts := docker.LoadImageOptions{InputStream: tar}
	err = client.LoadImage(opts)
	if nil != err {
		log.Fatal(err)
	}

}
