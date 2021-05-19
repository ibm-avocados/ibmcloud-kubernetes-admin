package kubeadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func GetBillingHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		var body map[string]interface{}
		decoder := json.NewDecoder(c.Request().Body)
		err = decoder.Decode(&body)
		if err != nil {
			return err
		}

		crn, ok := body["crn"]
		if !ok {
			return errors.New("No valid CRN")
		}

		clusterCRN := fmt.Sprintf("%v", crn)

		accnt, ok := body["accountID"]
		if !ok {
			return errors.New("No valid account ID")
		}

		accountID := fmt.Sprintf("%v", accnt)

		clustr, ok := body["clusterID"]
		if !ok {
			return errors.New("No valid Cluster ID")
		}

		clusterID := fmt.Sprintf("%v", clustr)

		billing, err := session.GetBillingData(accountID, clusterID, clusterCRN)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, Bill{Bill: billing})
	}
}
