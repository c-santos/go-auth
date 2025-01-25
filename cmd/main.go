package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/c-santos/go-auth/internal/auth"
	"github.com/c-santos/go-auth/internal/config"
)

type Verify struct {
	AccessToken string `json:"access_token"`
}

func main() {
	port := config.LoadConfig().Port

	mux := &http.ServeMux{}

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, I am running on %q", port)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)

		if r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		var body map[string]string
		// Decode r.Body and put it in the memory addr of the body var
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Printf("Couldn't decode request body: %s", err)
			w.WriteHeader(400)
			return
		}

		token, err := auth.GenerateToken(body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		response := map[string]string{
			"access_token": token,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)

		if r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		var body Verify

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.AccessToken == "" {
			w.WriteHeader(400)
			return
		}

		claims, err := auth.VerifyToken(body.AccessToken)
		if err != nil {
			w.WriteHeader(401)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		response := map[string]interface{}{
			"exp": *exp,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	})

	log.Printf("Listening on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
