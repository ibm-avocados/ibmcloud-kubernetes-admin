package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/cron"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/server"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func init() {
	ibmcloud.SetupCloudant()
}

func main() {
	cron.Start()

	server := server.NewServer()
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	///v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5
	api.HandleFunc("/identity-endpoints", server.TokenEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/authenticate/account", server.AuthenticationWithAccountHandler).Methods(http.MethodPost)
	api.HandleFunc("/authenticate", server.AuthenticationHandler).Methods(http.MethodPost)
	api.HandleFunc("/accounts", server.AccountListHandler).Methods(http.MethodGet)
	api.HandleFunc("/login", server.LoginHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", server.ClusterListHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters", server.ClusterCreateHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters", server.ClusterDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/resourcegroups/{accountID}", server.ResourceGroupHandler).Methods(http.MethodGet)

	// public endpoints

	api.HandleFunc("/clusters/versions", server.VersionEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/locations", server.LocationEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/{geo}/locations", server.LocationGeoEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/zones", server.ZonesEndpointHandler).
		Queries("showFlavors", "{showFlavors}", "location", "{location}").
		Methods(http.MethodGet)
	api.HandleFunc("/clusters/{datacenter}/machine-types", server.MachineTypeHandler).
		Queries("type", "{type}", "os", "{os}", "cpuLimit", "{cpuLimit}", "memoryLimit", "{memoryLimit}").
		Methods(http.MethodGet)

	api.HandleFunc("/clusters/{datacenter}/vlans", server.VlanEndpointHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/{clusterID}", server.ClusterHandler).Methods(http.MethodGet)
	api.HandleFunc("/clusters/{clusterID}/workers", server.ClusterWorkerListHandler).Methods(http.MethodGet)

	api.HandleFunc("/clusters/settag", server.SetTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/{clusterID}/settag", server.SetClusterTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/deletetag", server.DeleteTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/clusters/gettag", server.GetTagHandler).Methods(http.MethodPost)
	api.HandleFunc("/billing", server.GetBillingHandler).Methods(http.MethodPost)

	// scheduling
	api.HandleFunc("/schedule/api/create", server.SetAPITokenHandler).Methods(http.MethodPost)
	api.HandleFunc("/schedule/api", server.DeleteAPITokenHandler).Methods(http.MethodDelete)
	api.HandleFunc("/schedule/api", server.UpdateAPITokenHandler).Methods(http.MethodPut)
	api.HandleFunc("/schedule/api", server.CheckAPITokenHandler).Methods(http.MethodPost)
	api.HandleFunc("/schedule/{accountID}/create", server.SetScheduleHandler).Methods(http.MethodPost)
	api.HandleFunc("/schedule/{accountID}/all", server.GetAllScheduleHandler).Methods(http.MethodGet)
	api.HandleFunc("/schedule/{accountID}", server.GetScheduleHandler).Methods(http.MethodGet)
	api.HandleFunc("/schedule/{accountID}", server.UpdateScheduleHandler).Methods(http.MethodPut)
	api.HandleFunc("/schedule/{accountID}", server.DeleteScheduleHandler).Methods(http.MethodDelete)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("client/build/"))).Methods("GET")

	port := ":9000"

	log.Println("starting server on port ", port)

	log.Fatalln(http.ListenAndServe(port, r))
}
