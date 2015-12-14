package main

import (
	"log"
	"net/http"

	"github.com/antonholmquist/jason"
)

func getExpvarData(url string) (*ExpvarData, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get(%v) error. err: %v", url, err)
		return nil, err
	}
	defer res.Body.Close()

	object, err := jason.NewObjectFromReader(res.Body)
	if err != nil {
		log.Printf("json decode error!! err: %v", err)
		return nil, err
	}

	return &ExpvarData{object}, nil
}
