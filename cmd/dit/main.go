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
				return addImage(imageName)
			},
		},
		// {
		// 	Name:    "remove",
		// 	Aliases: []string{"r"},
		// 	Usage:   "Remove an image from the store",
		// 	Action: func(c *cli.Context) error {
		// 		imageName := c.Args().First()
		// 		removeImage(imageName)
		// 		return nil
		// 	},
		// },
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all images in the store",
			Action: func(c *cli.Context) error {
				return listImages()
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"l"},
			Usage:   "Sync all images in the store with active ecr",
			Action: func(c *cli.Context) error {
				return syncImages()
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listImages() error {
	resp, err := http.Get("http://localhost:6000/list")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func syncImages() error {
	resp, err := http.Get("http://localhost:6000/sync")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return nil
}

func createHttpRequest(uri string, image *repo.Image, imageName string) (*http.Request, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(image)
	if err != nil {
		return nil, err
	}

	log.Print(b)
	res, err := http.Post(uri, "application/json; charset=utf-8", b)
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, res.Body)

	return nil, nil
}

func addImage(imageName string) error {
	// get image meta data
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

	_, err = createHttpRequest("http://localhost:6000/add", image, imageName)

	return err
}

// func removeImage(imageName string) error {
// 	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	img, buffer, err := client.ImageInspectWithRaw(context.Background(), imageName)
// 	image := repo.NewImage(buffer, img)

// 	log.Print(image.Meta)

// 	_, err = createHttpRequest("http://localhost:6000/remove", image, imageName)

// 	return err
// }
