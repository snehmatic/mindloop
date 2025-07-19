package main

import (
	"github.com/gorilla/mux"
	v1 "github.com/snehmatic/mindloop/api/v1"
)

func CreateRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", v1.HandleHome)
	r.HandleFunc("/heatlz", v1.HandleHealthz)

	return r, nil
}
