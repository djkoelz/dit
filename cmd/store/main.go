package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/djkoelz/dit/pkg/repo"
	docker "github.com/docker/docker/client"
	dockerF "github.com/fsouza/go-dockerclient"
	"github.com/urfave/cli"
	"io"
	"log"
	//"mime/multipart"
	"net/http"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add an image to the store",
			Action: func(c *cli.Context) error {
				imageName := c.Args().First()
				log.Print("Adding image ", imageName)
				addImage(imageName)
				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "Remove an image from the store",
			Action: func(c *cli.Context) error {
				imageName := c.Args().First()
				log.Print("Removing image ", imageName)
				removeImage(imageName)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all images in the store",
			Action: func(c *cli.Context) error {
				log.Print("Lising images")
				listImages()
				return nil
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"l"},
			Usage:   "Sync all images in the store with active ecr",
			Action: func(c *cli.Context) error {
				log.Print("Syncing images")
				syncImages()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createHttpRequest(uri string, image *repo.Image, imageName string) (*http.Request, error) {
	b := new(bytes.Buffer)
	decoder := json.NewEncoder(b)
	err := decoder.Encode(image)
	if err != nil {
		log.Print(err)
	}

	res, _ := http.Post(uri, "application/json; charset=utf-8", b)
	io.Copy(os.Stdout, res.Body)

	return nil, nil
}

func listImages() error {
	res, err := http.Get("http://localhost:6000/list")
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, res.Body)
	return nil
}

func syncImages() error {
	res, err := http.Get("http://localhost:6000/sync")
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, res.Body)
	return nil
}

func addImage(imageName string) error {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := client.ImageInspectWithRaw(context.Background(), imageName)

	// extract tar representation of the image
	var buf bytes.Buffer
	clientF, err := dockerF.NewClientFromEnv()
	if err != nil {
		return err
	}

	opts := dockerF.ExportImageOptions{Name: imageName, OutputStream: &buf}
	err = clientF.ExportImage(opts)
	if err != nil {
		return err
	}

	image := repo.NewImage(buf.Bytes(), img)

	log.Print(image.Meta)

	_, err = createHttpRequest("http://localhost:6000/add", image, imageName)

	return err
}

func removeImage(imageName string) error {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	img, buffer, err := client.ImageInspectWithRaw(context.Background(), imageName)
	image := repo.NewImage(buffer, img)

	log.Print(image.Meta)

	_, err = createHttpRequest("http://localhost:6000/remove", image, imageName)

	return err
}
