package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/identity-endpoints", tokenEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/authenticate/account", authenticationWithAccountHandler).Methods(http.MethodPost)
	api.HandleFunc("/authenticate", authenticationHandler).Methods(http.MethodPost)
	api.HandleFunc("/accounts", accountListHandler).Methods(http.MethodGet)
	api.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", clusterListHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", clusterDeleteHandler).Methods(http.MethodDelete)

	r.HandleFunc("/", notFoundHandler)
	port := "9000"

	log.Println("starting server on port ", port)

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
