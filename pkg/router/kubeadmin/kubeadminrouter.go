package kubeadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/infra"
)

type router struct {
	*echo.Echo
}

func NewRouter(provider infra.SessionProvider) (*router, error) {
	e := echo.New()
	e.Use(
		middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}),
		middleware.Secure(),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  "client/build",
			HTML5: true,
		}),
	)

	static := e.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=31536000")
			return next(c)
		}
	})

	static.Use(middleware.Static("client/build/static"))

	auth := e.Group("/auth")

	auth.GET("", AuthHandler)
	auth.GET("/callback", AuthDoneHandler(provider)) //url/auth/callback
	auth.POST("/logout", LogoutHandler(provider))

	api := e.Group("/api/v1")
	api.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, time=${latency_human}\n",
	}))

	api.POST("/auth/check", CheckApiKeyHandler(provider))

	api.GET("/clusters/versions", VersionEndpointHandler(provider))
	api.GET("/clusters/locations", LocationEndpointHandler(provider))
	api.GET("/clusters/:geo/locations", LocationGeoEndpointHandler(provider))
	api.GET("/clusters/zones", ZonesEndpointHandler(provider))
	api.GET("/clusters/:datacenter/machine-types", MachineTypeHandler(provider))

	api.GET("/clusters/locations/info", LocationGeoEndpointInfoHandler(provider))

	///v1/resource_groups?account_id=9b13b857a32341b7167255de717172f5
	api.GET("/identity-endpoints", TokenEndpointHandler(provider))
	api.POST("/authenticate/account", AuthenticationWithAccountHandler(provider))
	api.GET("/accounts", AccountListHandler(provider))
	api.GET("/login", LoginHandler(provider))
	api.GET("/clusters", ClusterListHandler(provider))
	api.POST("/clusters", ClusterCreateHandler(provider))
	api.DELETE("/clusters", ClusterDeleteHandler(provider))
	api.GET("/resourcegroups/:accountID", ResourceGroupHandler(provider))
	api.GET("/iam/accessGroups/:accountID", AccessGroupsHandler(provider))
	api.PUT("/iam/accessGroups/:accessGroupID/members", AddMemberHandler(provider))
	api.GET("/iam/accessGroups/:accessGroupID/members/:iamID", MembershipCheckHandler(provider))
	api.POST("/iam/policies", CreatePolicyHandler(provider))
	api.POST("/users/account/:accountID", InviteUserHandler(provider))
	api.GET("/user/info", UserInfoHandler(provider))

	api.GET("/clusters/:datacenter/vlans", VlanEndpointHandler(provider))
	api.GET("/clusters/:clusterID", ClusterHandler(provider))
	api.GET("/clusters/:clusterID/workers", ClusterWorkerListHandler(provider))

	api.POST("/clusters/settag", SetTagHandler(provider))
	api.POST("/clusters/:clusterID/settag", SetClusterTagHandler(provider))
	api.POST("/clusters/deletetag", DeleteTagHandler(provider))
	api.POST("/clusters/gettag", GetTagHandler(provider))
	api.POST("/billing", GetBillingHandler(provider))

	api.POST("/github/comment", GithubCommentHandler)

	// api.POST("/awx/cluster", CreateClusterWithAWX)
	api.GET("/awx/workflowjobtemplate", GetAWXWorkflowJobTemplates)
	api.GET("/awx/jobtemplate", GetAWXJobTemplates)
	api.POST("/awx/workflowjobtemplate/launch", LaunchAWXWorkflowJobTemplate)
	api.GET("/awx/grantclusterid", GetGrantClusterID)
	return &router{e}, nil
}

func (r *router) Serve(port string) error {
	return r.Start(port)
}
