package main

import (
	"errors"
	"io"
	"net/http"
)

// Gets a reader for the resource located at the url
func getUrlReader(imageUrl string) (io.Reader, error) {
	response, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("Non 200 response")
	}
	return response.Body, nil
}
