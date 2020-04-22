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

	api.HandleFunc("/clusters/{clusterID}/workers", clusterWorkerListHandler).Methods(http.MethodGet)

	api.HandleFunc("/clusters/settag", setTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/deletetag", deleteTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/gettag", getTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/billing/{accountID}/{clusterID}/{clusterCRN}", getBillingHandler).Methods(http.MethodGet)

	api.HandleFunc("/clusters/locations", locationEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/zones", zonesEndpointHandler).
		Queries("showFlavors", "{showFlavors}", "location", "{location}").
		Methods(http.MethodGet)

	api.HandleFunc("/cluster/machineType", machineTypeHandler).
		Queries("machineType", "{machineType}").Methods(http.MethodGet)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("client/build/"))).Methods("GET")
	r.HandleFunc("/", notFoundHandler)

	port := "9000"

	log.Println("starting server on port ", port)

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
