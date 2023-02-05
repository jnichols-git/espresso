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
