package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/jakenichols2719/simpleblog/pkg/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/blog/host/create", handlers.HandleCreateRequest)
	r.HandleFunc("/blog/host/update", handlers.HandleUpdateRequest)
	r.Use(handlers.HandleHostAuth)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", req.RequestURI)))
	})
	lambda.Start(gorillamux.NewV2(r).ProxyWithContext)
	//http.ListenAndServe(":2718", r)
}
