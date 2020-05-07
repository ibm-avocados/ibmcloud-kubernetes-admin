package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	go timed()
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	///v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5
	api.HandleFunc("/identity-endpoints", tokenEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/authenticate/account", authenticationWithAccountHandler).Methods(http.MethodPost)
	api.HandleFunc("/authenticate", authenticationHandler).Methods(http.MethodPost)
	api.HandleFunc("/accounts", accountListHandler).Methods(http.MethodGet)
	api.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", clusterListHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", clusterCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters", clusterDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/resourcegroups/{accountID}", resourceGroupHandler).Methods(http.MethodGet)

	api.HandleFunc("/clusters/{datacenter}/vlans", vlanEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/{clusterID}/workers", clusterWorkerListHandler).Methods(http.MethodGet)

	api.HandleFunc("/clusters/settag", setTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/deletetag", deleteTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/gettag", getTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/billing", getBillingHandler).Methods(http.MethodPost)

	// public endpoints

	api.HandleFunc("/clusters/versions", versionEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/locations", locationEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/{geo}/locations", locationGeoEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/zones", zonesEndpointHandler).
		Queries("showFlavors", "{showFlavors}", "location", "{location}").
		Methods(http.MethodGet)
	api.HandleFunc("/clusters/{datacenter}/machine-types", machineTypeHandler).
		Queries("type", "{type}", "os", "{os}", "cpuLimit", "{cpuLimit}", "memoryLimit", "{memoryLimit}").
		Methods(http.MethodGet)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("client/build/"))).Methods("GET")
	r.HandleFunc("/", notFoundHandler)

	port := ":9000"

	log.Println("starting server on port ", port)

	log.Fatalln(http.ListenAndServe(port, r))
}
