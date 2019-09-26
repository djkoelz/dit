package repo

import (
	"github.com/docker/docker/api/types"
)

type Image struct {
	Dockerfile string
	Meta       types.ImageInspect
}

func NewImage(meta types.ImageInspect) *Image {
	image := new(Image)
	image.Dockerfile = ""
	image.Meta = meta

	return image
}

// func (this *Image) GetImageLayerFiles() []string {
// 	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	inspection, _, err := client.ImageInspectWithRaw(context.Background(), )
// }
