package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"github.com/djkoelz/dit/pkg/repo"
	"github.com/djkoelz/dit/pkg/router"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var buf bytes.Buffer
var tw = tar.NewWriter(&buf)

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
	// Open and iterate through the files in the archive.
	tw.Close()
	tarballFilePath := "file.tar"
	file, err := os.Create(tarballFilePath)
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarballFilePath, err.Error())))
	}
	defer file.Close()

	io.Copy(file, &buf)

	log.Print("Done")
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
	router.Register("/get", service.GetImage)
	router.Register("/remove", service.RemoveImage)
	router.Register("/list", service.ListImages)
	//router.Register("/sync", service.SyncImages)
	router.Register("/sync", sync)

	router.Start(6000)
}
