package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/identity-endpoints", tokenEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/authenticate", authenticationHandler).Methods(http.MethodPost)
	api.HandleFunc("/accounts", accountListHandler).Methods(http.MethodGet)
	api.HandleFunc("/login", loginHandler).Methods(http.MethodGet)

	r.HandleFunc("/", notFoundHandler)

	http.ListenAndServe(":9000", r)
}
