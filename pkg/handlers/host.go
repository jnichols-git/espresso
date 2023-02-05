package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jakenichols2719/simpleblog/pkg/events"
)

// Kick host requests back without a hardcoded API key
func HandleHostAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := os.Getenv("HOST_API_KEY")
		if req == "" {
			HandleError(w, errors.New("missing HOST_API_KEY env"))
		} else {
			req = "Bearer " + req
		}
		auth := r.Header.Get("Authorization")
		if auth != req {
			HandleError(w, errors.New("incorrect authorization key"))
		}
		next.ServeHTTP(w, r)
	})
}

func HandleCreateRequest(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		HandleError(w, err)
	}
	event := &events.Event{}
	err = json.Unmarshal(body, event)
	if err != nil {
		HandleError(w, err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
	}
	client := events.Connect(cfg)
	err = events.Schedule(event, client)
	if err != nil {
		HandleError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func HandleUpdateRequest(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		HandleError(w, err)
	}
	event := &events.Event{}
	err = json.Unmarshal(body, event)
	if err != nil {
		HandleError(w, err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
	}
	client := events.Connect(cfg)
	events.Schedule(event, client)
	if err != nil {
		HandleError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}
