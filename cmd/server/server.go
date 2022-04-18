package main

import (
	"net/http"

	"github.com/Schnides123/wordle-vs/pkg/endpoints"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	endpoints.SetupRoutes(r)
	http.ListenAndServe(":42069", r)
}
