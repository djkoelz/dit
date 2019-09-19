package repo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	service := new(Service)
	service.store = store

	return service
}

func (this *Service) AddImage(w http.ResponseWriter, r *http.Request) {
	// var u User
	var image Image
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	this.store.AddImage(&image)

	w.WriteHeader(http.StatusOK)
}

func (this *Service) ListImages(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, this.store.CreateImagesString())
	w.WriteHeader(http.StatusOK)
}

func (this *Service) RemoveImage(w http.ResponseWriter, r *http.Request) {
	log.Print("In Remove image")
	// var u User
	var image Image
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	this.store.RemoveImage(image.Meta.ID)

	w.WriteHeader(http.StatusOK)
}

func (this *Service) SyncImages(w http.ResponseWriter, r *http.Request) {
	this.store.Sync()
}

func (this *Service) GetImage(w http.ResponseWriter, r *http.Request) {
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
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}

			imageName := string(slurp)
			log.Printf("Getting %s", imageName)

			w.WriteHeader(http.StatusOK)
			//w.Write(buf.Bytes())
		}
	}
}
