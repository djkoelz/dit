package image

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func Push(image []byte) error {
	client, err := docker.NewClientWithOpts(docker.WithVersion("1.39"))
	if err != nil {
		return err
	}

	stream := bytes.NewBuffer(image)
	_, err = client.ImageLoad(context.Background(), stream, false)
	return err
}

func Pull(imageName string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return buf, err
	}

	reader, err := client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return buf, err
	}

	buf.ReadFrom(reader)

	return buf, err
}

func Remove(imageName string) error {
	client, err := docker.NewClientWithOpts(docker.WithVersion("1.39"))
	if err != nil {
		return err
	}

	opts := types.ImageRemoveOptions{Force: true}
	_, err = client.ImageRemove(context.Background(), imageName, opts)

	return err
}
