package main

import (
	"bytes"
	"fmt"
	//"io/ioutil"
	"github.com/djkoelz/dit/pkg/image"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Creates a new file push http request with optional extra params
func push(uri string, params map[string]string, image bytes.Buffer) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", params["title"])
	if err != nil {
		return nil, err
	}
	part.Write(image.Bytes())

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}

func main() {
	args := os.Args

	repoLocation := args[1]
	imageName := args[2]
	url := fmt.Sprintf("http://%s:5000/push", repoLocation)

	// get the image data
	buf, err := image.Pull(imageName)
	if err != nil {
		log.Fatal(err)
	}

	params := map[string]string{
		"title": imageName,
	}

	request, err := push(url, params, buf)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)
	}
}
