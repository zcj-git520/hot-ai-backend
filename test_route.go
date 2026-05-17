package main

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func main() {
	var c rest.RestConf
	c.Host = "0.0.0.0"
	c.Port = 8007
	
	server := rest.MustNewServer(c)
	defer server.Stop()
	
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path: "/test/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")
			fmt.Fprintf(w, "ID: %s", id)
		},
	})
	
	fmt.Println("Starting test server on 8007...")
	server.Start()
}
