package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cloudretic/matcha/pkg/route"
	"github.com/cloudretic/matcha/pkg/router"
	"github.com/jakenichols2719/simpleblog/pkg/adapter"
	"github.com/jakenichols2719/simpleblog/pkg/handlers"
)

func main() {
	rt := router.Declare(
		router.Default(),
		router.WithRoute(route.Declare(http.MethodPost, "/blog/host/create"), http.HandlerFunc(handlers.HandleCreateRequest)),
		router.WithRoute(route.Declare(http.MethodPost, "/blog/host/update"), http.HandlerFunc(handlers.HandleUpdateRequest)),
		router.WithMiddleware(handlers.HandleHostAuth),
		router.WithNotFound(http.HandlerFunc(handlers.HandleNotFound)),
	)
	lambda.Start(adapter.LambdaAdapter(rt))
	/*
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
	*/
	//http.ListenAndServe(":2718", r)
}
