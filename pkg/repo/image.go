package repo

import (
	"github.com/docker/docker/api/types"
)

type Image struct {
	Dockerfile string
	Data       []byte
	Meta       types.ImageInspect
}

func NewImage(data []byte, meta types.ImageInspect) *Image {
	image := new(Image)
	image.Dockerfile = ""
	image.Data = data
	image.Meta = meta

	return image
}
