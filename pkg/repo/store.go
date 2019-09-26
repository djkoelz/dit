package repo

import (
	//"bytes"
	"context"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	//"io"
	"log"
	//"os"
)

type Store struct {
	registry string            // registry uri
	images   map[string]*Image // image map
}

func NewStore(registry string) *Store {
	store := new(Store)
	store.registry = registry
	store.images = make(map[string]*Image)

	return store
}

func (this *Store) AddImage(image *Image) {
	this.images[image.Meta.ID] = image

	// client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	// if err != nil {
	// 	log.Print(err)
	// }

	// r := bytes.NewReader(image.Data)
	// resp, err := client.ImageLoad(context.Background(), r, true)
	// defer resp.Body.Close()
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	io.Copy(os.Stdout, resp.Body)
	// }
}

func (this *Store) RemoveImage(id string) {
	log.Print("Removing: ", id)
	delete(this.images, id)
}

func (this *Store) Sync() {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	for _, img := range this.images {
		ctx := context.Background()

		// tag the image to include the ecr registry address
		imageName := img.Meta.RepoTags[0]
		registryImageName := this.registry + "/" + imageName
		err = client.ImageTag(ctx, imageName, registryImageName)
		if err != nil {
			log.Print(err)
			continue
		}

		// push the registry
		log.Printf("Pushing %s to %s", registryImageName, this.registry)
		_, err := client.ImagePush(ctx, registryImageName, types.ImagePushOptions{All: true, RegistryAuth: "0"})
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("%s Pushed", imageName)
		}
	}
}

func (this *Store) CreateImagesString() string {
	msg := "{ "
	for _, img := range this.images {
		msg += "[ "
		for _, tag := range img.Meta.RepoTags {
			msg += tag + " "
		}
		msg += "] "
	}
	msg += "}\n"
	return msg
}

func (this *Store) PrintImages() {
	for _, img := range this.images {
		log.Print(img.Meta.RepoTags)
	}
}
