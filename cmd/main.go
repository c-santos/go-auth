package main

import (
	"fmt"
	"github.com/c-santos/go-auth/internal/config"
	"log"
	"net/http"
)

func main() {

	port := config.LoadConfig().Port

	mux := &http.ServeMux{}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, I am running on %q", port)
	})

	log.Printf("Listening on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
