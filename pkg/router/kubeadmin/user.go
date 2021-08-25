package kubeadmin

import (
	"log"
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func UserInfoHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			log.Println(err)
			return err
		}

		userInfo, err := session.GetUserInfo()
		if err != nil {
			log.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, userInfo)
	}
}

func UserPreferenceHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		userID := c.Param("userID")

		if err != nil {
			log.Println(err)
			return err
		}

		userPreference, err := session.GetUserPreference(userID)
		if err != nil {
			log.Println(err)
			return err
		}
		return c.JSON(http.StatusOK, userPreference)
	}
}
