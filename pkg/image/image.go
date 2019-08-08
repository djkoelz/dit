package image

import (
	"bytes"
	docker "github.com/fsouza/go-dockerclient"
)

func Push(image []byte) error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	stream := bytes.NewBuffer(image)
	opts := docker.LoadImageOptions{InputStream: stream}
	return client.LoadImage(opts)
}

func Remove(imageName string) error {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	return client.RemoveImage(imageName)
}

func Pull(imageName string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return buf, err
	}

	opts := docker.ExportImageOptions{Name: imageName, OutputStream: &buf}
	err = client.ExportImage(opts)

	return buf, err
}
