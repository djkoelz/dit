package main

import (
	"archive/tar"
	"bytes"
	"context"
	"github.com/djkoelz/dit/pkg/repo"
	"github.com/djkoelz/dit/pkg/router"
	docker "github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

var buf = new(bytes.Buffer)
var tw = tar.NewWriter(buf)

func test(w http.ResponseWriter, r *http.Request) {

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Print(err)
			}
			body, err := ioutil.ReadAll(p)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Printf("Loading %s", p.FileName())
			hdr := &tar.Header{
				Name: p.FileName(),
				Mode: 0600,
				Size: int64(len(body)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatal(err)
			}
			if _, err := tw.Write(body); err != nil {
				log.Fatal(err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func sync(w http.ResponseWriter, r *http.Request) {
	log.Print("Syncing file")

	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	tw.Close()
	response, err := client.ImageLoad(context.Background(), buf, false)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Print(bodyString)

	// Open and iterate through the files in the archive.
	// tarballFilePath := "file.tar"
	// file, err := os.Create(tarballFilePath)
	// if err != nil {
	// 	log.Fatal(errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error())))
	// }
	// defer file.Close()

	// io.Copy(file, &buf)

	// log.Print("Done")

	// tr := tar.NewReader(&buf)
	// for {
	// 	hdr, err := tr.Next()
	// 	if err == io.EOF {
	// 		break // End of archive
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Contents of %s:\n", hdr.Name)
	// 	if _, err := io.Copy(file, tr); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println()
	// }
}

func main() {
	store := repo.NewStore("localhost:5000")
	service := repo.NewService(store)

	router := router.NewRouter()
	router.Register("/add", test)
	//router.Register("/add", sevice.AddImage)
	router.Register("/get", service.GetImage)
	router.Register("/remove", service.RemoveImage)
	router.Register("/list", service.ListImages)
	//router.Register("/sync", service.SyncImages)
	router.Register("/sync", sync)

	router.Start(6000)
}
