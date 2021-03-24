package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/server"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/ibmcloud"
)

func main() {
	ibmcloud.SetupCloudant()

	e := echo.New()
	e.Use(
		middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}),
	)

	e.Use(middleware.Secure())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "client/build",
		HTML5: true,
	}))

	// Set cache control to a year on static directory. Filenames are hashed so
	// they can safely be aggressively cached.
	static := e.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=31536000")
			return next(c)
		}
	})
	static.Use(middleware.Static("client/build/static"))

	auth := e.Group("/auth")

	auth.GET("", server.AuthHandler)
	auth.GET("/callback", server.AuthDoneHandler) //url/auth/callback
	auth.POST("/logout", server.LogoutHandler)

	api := e.Group("/api/v1")
	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${latency_human}\n",
	}))

	api.POST("/auth/check", server.CheckApiKeyHandler)

	api.GET("/clusters/versions", server.VersionEndpointHandler)
	api.GET("/clusters/locations", server.LocationEndpointHandler)
	api.GET("/clusters/:geo/locations", server.LocationGeoEndpointHandler)
	api.GET("/clusters/zones", server.ZonesEndpointHandler)
	api.GET("/clusters/:datacenter/machine-types", server.MachineTypeHandler)

	api.GET("/clusters/locations/info", server.LocationGeoEndpointInfoHandler)

	///v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5
	api.GET("/identity-endpoints", server.TokenEndpointHandler)
	api.POST("/authenticate/account", server.AuthenticationWithAccountHandler)
	api.POST("/authenticate", server.AuthenticationHandler)
	api.GET("/accounts", server.AccountListHandler)
	api.GET("/login", server.LoginHandler)
	api.GET("/clusters", server.ClusterListHandler)
	api.POST("/clusters", server.ClusterCreateHandler)
	api.DELETE("/clusters", server.ClusterDeleteHandler)
	api.GET("/resourcegroups/:accountID", server.ResourceGroupHandler)
	api.GET("/iam/accessGroups/:accountID", server.AccessGroupsHandler)
	api.PUT("/iam/accessGroups/:accessGroupID/members", server.AddMemberHandler)
	api.GET("/iam/accessGroups/:accessGroupID/members/:iamID", server.MembershipCheckHandler)
	api.POST("/iam/policies", server.CreatePolicyHandler)
	api.POST("/users/account/:accountID", server.InviteUserHandler)
	api.GET("/user/info", server.UserInfoHandler)

	api.GET("/clusters/:datacenter/vlans", server.VlanEndpointHandler)
	api.GET("/clusters/:clusterID", server.ClusterHandler)
	api.GET("/clusters/:clusterID/workers", server.ClusterWorkerListHandler)

	api.POST("/clusters/settag", server.SetTagHandler)
	api.POST("/clusters/:clusterID/settag", server.SetClusterTagHandler)
	api.POST("/clusters/deletetag", server.DeleteTagHandler)
	api.POST("/clusters/gettag", server.GetTagHandler)
	api.POST("/billing", server.GetBillingHandler)

	// scheduling
	api.POST("/schedule/api/create", server.SetAPITokenHandler)
	api.DELETE("/schedule/api", server.DeleteAPITokenHandler)
	api.PUT("/schedule/api", server.UpdateAPITokenHandler)
	api.POST("/schedule/api", server.CheckAPITokenHandler)
	api.POST("/schedule/:accountID/create", server.SetScheduleHandler)
	api.GET("/schedule/:accountID/all", server.GetAllScheduleHandler)

	// api.HandleFunc("/schedule/{accountID}", server.GetScheduleHandler)
	api.PUT("/schedule/:accountID", server.UpdateScheduleHandler)
	api.DELETE("/schedule/:accountID", server.DeleteScheduleHandler)

	api.POST("/workshop/:accountID/metadata", server.CreateMetaDataHandler)
	api.PUT("/workshop/:accountID/metadata", server.UpdateMetaDataHandler)
	api.GET("/workshop/accountID/metadata", server.GetMetaDataHandler)

	api.POST("/github/comment", server.GithubCommentHandler)

	api.GET("/notification/:accountID/email", server.GetAdminEmails)
	api.POST("/notification/email/create", server.CreateAdminEmails)
	api.PUT("/notification/email/add", server.AddAdminEmails)
	api.PUT("/notification/email/remove", server.RemoveAdminEmails)
	api.DELETE("/notification/email", server.DeleteAdminEmails)

	// api.POST("/awx/cluster", server.CreateClusterWithAWX)
	api.GET("/awx/workflowjobtemplate", server.GetAWXWorkflowJobTemplates)
	api.GET("/awx/jobtemplate", server.GetAWXJobTemplates)
	api.POST("/awx/workflowjobtemplate/launch", server.LaunchAWXWorkflowJobTemplate)

	// spa := spaHandler{staticPath: "client/build", indexPath: "index.html"}
	// r.PathPrefix("/").Handler(spa)

	port := ":9000"

	log.Println("starting server on port serving index", port)

	e.Logger.Fatal(e.Start(port))
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
