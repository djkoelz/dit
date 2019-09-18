package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"log"
)

func main() {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	imgs, err := client.ImageList(context.Background(), types.ImageListOptions{All: false})
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

	//images := []string{"hello-world"}
	//reader, err := client.ImageSave(context.Background(), images)
	image, _, err := client.ImageInspectWithRaw(context.Background(), "hello-world")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(image)
	//var buf bytes.Buffer
	// buf.ReadFrom(reader)

	// fmt.Println(buf)

}
