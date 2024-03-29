package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"github.com/cheggaaa/pb/v3"
	"github.com/djkoelz/dit/pkg/repo"
	docker "github.com/docker/docker/client"
	dockerF "github.com/fsouza/go-dockerclient"
	"github.com/urfave/cli"
	"io"
	//"io/ioutil"
	"log"
	"mime/multipart"
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
		{
			Name:    "inspect",
			Aliases: []string{"a"},
			Usage:   "Add an image to the store",
			Action: func(c *cli.Context) error {
				imageName := c.Args().First()
				return inspectImage(imageName)
			},
		},
		{
			Name:    "save",
			Aliases: []string{"a"},
			Usage:   "Add an image to the store",
			Action: func(c *cli.Context) error {
				imageName := c.Args().First()
				return saveImage(imageName)
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

func inspectImage(imageName string) error {
	// get image meta data
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	img, _, err := client.ImageInspectWithRaw(context.Background(), imageName)

	log.Print(img.RootFS.Layers)

	return nil
}

type LayerFile struct {
	Name string
	Body []byte
}

func createTarFiles(layers map[string][]LayerFile) (map[string]*bytes.Buffer, error) {
	tars := make(map[string]*bytes.Buffer)

	// create tar file for each layer containing corresponding files
	for k, v := range layers {
		tars[k] = new(bytes.Buffer)
		tw := tar.NewWriter(tars[k])
		for _, file := range v {
			hdr := &tar.Header{
				Name: file.Name,
				Mode: 0600,
				Size: int64(len(file.Body)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return tars, err
			}
			if _, err := tw.Write(file.Body); err != nil {
				return tars, err
			}
		}
		if err := tw.Close(); err != nil {
			return tars, err
		}
	}

	return tars, nil
}

func transmitLayers(reader *tar.Reader, size int) map[string][]LayerFile {
	layers := make(map[string][]LayerFile)

	bar := pb.Start64(3000000)
	bar.Start()

	// for each layer transmit
	for {
		hdr, err := reader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}

		r, w := io.Pipe()
		m := multipart.NewWriter(w)

		go func() {
			defer w.Close()
			defer m.Close()

			part, err := m.CreateFormFile("file", hdr.Name)
			if err != nil {
				log.Fatal(err)
			}

			written, err := io.Copy(part, reader)
			if err != nil {
				log.Fatal(err)
			}

			bar.Add64(written)
		}()

		bar.Finish()
		request, err := http.NewRequest("POST", "http://localhost:6000/add", r)
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Content-Type", m.FormDataContentType())
		httpClient := &http.Client{}
		_, err = httpClient.Do(request)
		if err != nil {
			log.Fatal(err)
		}

	}

	return layers
}

func saveImage(imageName string) error {
	// get image meta data
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	images := []string{imageName}
	response, err := client.ImageSave(context.Background(), images)
	defer response.Close()

	// create layers from tar
	transmitLayers(tar.NewReader(response), 1000)

	// send sync command
	resp, err := http.Get("http://localhost:6000/sync")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// io.Copy(os.Stdout, resp.Body)

	return nil
}

func addImage(imageName string) error {
	// get image meta data
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return err
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

	image := repo.NewImage(img)

	_, err = createHttpRequest("http://localhost:6000/add", image, imageName)

	return err
}
