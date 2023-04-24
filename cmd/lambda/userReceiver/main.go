package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/route"
	"github.com/cloudretic/matcha/pkg/router"
	"github.com/jakenichols2719/simpleblog/pkg/adapter"
	"github.com/jakenichols2719/simpleblog/pkg/handlers"
)

func main() {
	/*
		r := mux.NewRouter()
		r.HandleFunc("/blog/user/listings", handlers.HandleGetListings)
		r.HandleFunc("/blog/user/listing", handlers.HandleGetListing)
		r.HandleFunc("/blog/user/post", handlers.HandleGetPost)
	*/
	rt := router.Declare(
		router.Default(),
		router.WithRoute(route.Declare(http.MethodGet, "/blog/user/listings"), http.HandlerFunc(handlers.HandleGetListings)),
		router.WithRoute(route.Declare(http.MethodGet, "/blog/user/listing"), http.HandlerFunc(handlers.HandleGetListing)),
		router.WithRoute(route.Declare(http.MethodGet, "/blog/user/post"), http.HandlerFunc(handlers.HandleGetPost)),
		router.WithMiddleware(middleware.LogRequests(os.Stdout)),
		router.WithNotFound(http.HandlerFunc(handlers.HandleNotFound)),
	)
	lambda.Start(adapter.LambdaAdapter(rt))
	//http.ListenAndServe(":2717", rt)
}
