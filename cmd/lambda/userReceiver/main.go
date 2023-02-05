package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/jakenichols2719/simpleblog/pkg/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/blog/user/listings", handlers.HandleGetListings)
	r.HandleFunc("/blog/user/listing", handlers.HandleGetListing)
	r.HandleFunc("/blog/user/post", handlers.HandleGetPost)
	lambda.Start(gorillamux.NewV2(r).ProxyWithContext)
	// http.ListenAndServe(":2717", http.DefaultServeMux)
}
