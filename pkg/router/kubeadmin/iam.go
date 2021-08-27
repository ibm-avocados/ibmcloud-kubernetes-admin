package kubeadmin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/labstack/echo/v4"
)

func AccessGroupsHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		accountID := c.Param("accountID")

		accessGroups, err := session.GetAccessGroups(accountID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, accessGroups)
	}
}

func InviteUserHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		accountID := c.Param("accountID")

		var body map[string]interface{}
		decoder := json.NewDecoder(c.Request().Body)
		if err := decoder.Decode(&body); err != nil {
			return err
		}

		_email, ok := body["email"]
		if !ok {
			return err
		}
		email := fmt.Sprintf("%v", _email)

		userInviteList, err := session.InviteUserToAccount(accountID, email)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, userInviteList)
	}
}

func AddMemberHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		accessGroupID := c.Param("accessGroupID")

		var body map[string]interface{}
		decoder := json.NewDecoder(c.Request().Body)
		if err := decoder.Decode(&body); err != nil {
			return err
		}

		_iamID, ok := body["iam_id"]
		if !ok {
			return err
		}
		iamID := fmt.Sprintf("%v", _iamID)

		_memberType, ok := body["type"]
		if !ok {
			return err
		}
		memberType := fmt.Sprintf("%v", _memberType)

		memberAddResponseList, err := session.AddMemberToAccessGroup(accessGroupID, iamID, memberType)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, memberAddResponseList)
	}
}

func CreatePolicyHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		var body map[string]interface{}
		decoder := json.NewDecoder(c.Request().Body)
		if err := decoder.Decode(&body); err != nil {
			return err
		}

		_accountID, ok := body["accountID"]
		if !ok {
			return err
		}
		accountID := fmt.Sprintf("%v", _accountID)

		_iamID, ok := body["iam_id"]
		if !ok {
			return err
		}
		iamID := fmt.Sprintf("%v", _iamID)

		_serviceName, ok := body["serviceName"]
		if !ok {
			return err
		}
		serviceName := fmt.Sprintf("%v", _serviceName)

		_serviceinstance, ok := body["serviceInstance"]
		if !ok {
			return err
		}
		serviceInstance := fmt.Sprintf("%v", _serviceinstance)

		_role, ok := body["role"]
		if !ok {
			return err
		}
		role := fmt.Sprintf("%v", _role)

		policy, err := session.CreatePolicy(accountID, iamID, serviceName, serviceInstance, role)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, policy)
	}
}

func MembershipCheckHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		accessGroupID := c.Param("accessGroupID")
		iamID := c.Param("iamID")

		err = session.IsMemberOfAccessGroup(accessGroupID, iamID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, StatusOK{Message: "User found in access group"})
	}
}
