package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jakenichols2719/simpleblog/pkg/events"
)

// Kick host requests back without a hardcoded API key
func HandleHostAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Authenticating request\n")
		req := os.Getenv("HOST_API_KEY")
		if req == "" {
			HandleError(w, errors.New("missing HOST_API_KEY env"))
			return
		} else {
			req = "Bearer " + req
		}
		auth := r.Header.Get("Authorization")
		if auth != req {
			HandleError(w, errors.New("incorrect authorization key"))
			return
		}
		fmt.Printf("Request authenticated, forwarding to handler\n")
		next.ServeHTTP(w, r)
	})
}

func HandleCreateRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Handling post creation request\n")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		HandleError(w, err)
		return
	}
	event := &events.Event{}
	err = json.Unmarshal(body, event)
	if err != nil {
		HandleError(w, err)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
		return
	}
	client := events.Connect(cfg)
	err = events.Schedule(event, client)
	if err != nil {
		HandleError(w, err)
		return
	}
	fmt.Printf("Post creation event created successfully\n")
	w.WriteHeader(http.StatusOK)
}

func HandleUpdateRequest(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Handling post update request\n")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		HandleError(w, err)
		return
	}
	event := &events.Event{}
	err = json.Unmarshal(body, event)
	if err != nil {
		HandleError(w, err)
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
		return
	}
	client := events.Connect(cfg)
	events.Schedule(event, client)
	if err != nil {
		HandleError(w, err)
		return
	}
	fmt.Printf("Post update event created successfully\n")
	w.WriteHeader(http.StatusOK)
}
