package handlers

import (
	"fmt"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func HandleBadRequest(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func HandleNotFound(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Not found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
