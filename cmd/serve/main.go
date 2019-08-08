package main

import (
	"bytes"
	"fmt"
	"github.com/djkoelz/dit/pkg/service"
	docker "github.com/fsouza/go-dockerclient"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	//"os"
	"strings"
)

func upload(w http.ResponseWriter, r *http.Request) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Part %q: %q\n", p.FileName(), slurp)

			// save the file submitted in the request
			if p.FileName() != "" {
				// f, err := os.Create(p.FileName())
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// defer f.Close()

				// n2, err := f.Write(slurp)
				// if err != nil {
				// 	log.Print(err)
				// }
				// fmt.Printf("wrote %d bytes\n", n2)

				// f.Sync()

				tar := bytes.NewBuffer(slurp)
				opts := docker.LoadImageOptions{InputStream: tar}

				err = client.LoadImage(opts)
				if nil != err {
					log.Fatal(err)
				}
			}
		}
	}
}

func main() {
	router := service.NewRouter()
	router.Register("/upload", upload)
	router.Start(5000)
}
