package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/cron"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/server"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func main() {
	ibmcloud.SetupCloudant()

	go cron.Start()

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

	api.HandleFunc("/notification/{accountID}/email", server.GetAdminEmails).Methods(http.MethodGet)
	api.HandleFunc("/notification/email/create", server.CreateAdminEmails).Methods(http.MethodPost)
	api.HandleFunc("/notification/email/add", server.AddAdminEmails).Methods(http.MethodPut)
	api.HandleFunc("/notification/email/remove", server.RemoveAdminEmails).Methods(http.MethodPut)
	api.HandleFunc("/notification/email", server.DeleteAdminEmails).Methods(http.MethodDelete)

	spa := spaHandler{staticPath: "client/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	port := ":9000"

	log.Println("starting server on port serving index", port)

	log.Fatalln(http.ListenAndServe(port, r))
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
