package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

func pull(uri string, params map[string]string) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}

func main() {
	params := map[string]string{
		"image": "hello-world",
	}
	request, err := pull("http://localhost:5000/pull", params)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
	var bodyContent []byte
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header)
	resp.Body.Read(bodyContent)
	resp.Body.Close()
	fmt.Println(bodyContent)
}
