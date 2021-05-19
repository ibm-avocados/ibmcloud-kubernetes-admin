package kubeadmin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/infra"

	"github.com/moficodes/ibmcloud-kubernetes-admin/internals/token"

	"github.com/labstack/echo/v4"
)

var oauthSettings = map[string]OauthSettings{
	"ibm": {
		authURL:  "https://iam.cloud.ibm.com/identity/authorize",
		tokenURL: "https://iam.cloud.ibm.com/identity/token",
		// authURL:      "https://identity-2.us-south.iam.cloud.ibm.com/identity/authorize",
		// tokenURL:     "https://identity-2.us-south.iam.cloud.ibm.com/identity/token",
		clientID:     os.Getenv("IBM_LOGIN_CLIENT_ID"),
		clientSecret: os.Getenv("IBM_LOGIN_CLIENT_SECRET"),
		redirectURI:  os.Getenv("IBM_REDIRECT_URI"), // "https://appstatic.dev/auth/callback",
	},
}

func buildRedirect(provider string, login bool, extraData string, settings OauthSettings) (string, error) {
	redirectURL, err := url.Parse(settings.authURL)
	if err != nil {
		return "", err
	}

	token, err := token.New(token.Claims{Provider: provider, Login: login, ExtraData: extraData})
	if err != nil {
		return "", err
	}

	query := redirectURL.Query()
	query.Set("client_id", settings.clientID)
	query.Set("redirect_uri", settings.redirectURI)
	query.Set("state", token)

	redirectURL.RawQuery = query.Encode()

	return redirectURL.String(), nil
}

func AuthHandler(c echo.Context) error {
	fmt.Printf("%+v", oauthSettings)
	provider := c.QueryParam("provider")
	// We don't need to crash for bad bool parsing.
	login, _ := strconv.ParseBool(c.QueryParam("login"))
	query := c.QueryString()
	queries := strings.Split(query, "&")
	extraData := strings.Join(queries[2:], "&")

	account := c.QueryParam("account")
	log.Println(account)
	settings, ok := oauthSettings[provider]
	if !ok {
		return errors.New("invalid provider")
	}

	redirectURL, err := buildRedirect(provider, login, extraData, settings)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

func AuthDoneHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		state := c.QueryParam("state")
		claims, err := token.Verify(state)

		if err != nil {
			return err
		}

		settings, ok := oauthSettings[claims.Provider]
		if !ok {
			return errors.New("invalid provider")
		}
		form := url.Values{}
		form.Set("client_id", settings.clientID)
		form.Set("client_secret", settings.clientSecret)
		form.Set("grant_type", "authorization_code")
		form.Set("redirect_uri", settings.redirectURI)
		form.Set("state", state)
		form.Set("code", code)

		client := &http.Client{}
		req, err := http.NewRequest("POST", settings.tokenURL, bytes.NewBufferString(form.Encode()))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		authToken := infra.AuthToken{}
		err = json.Unmarshal(data, &authToken)
		if err != nil {
			return err
		}

		session, err := provider.GetSessionWithToken(&authToken)
		if err != nil {
			return err
		}
		session.SetCookie(c.Response())
		return c.Redirect(http.StatusFound, "/?"+claims.ExtraData)
	}
}

func AuthenticationWithAccountHandler(provider infra.SessionProvider) echo.HandlerFunc {
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

		accountID := fmt.Sprintf("%v", body["id"])

		log.Println("Account id", accountID)

		accountSession, err := session.BindAccountToToken(accountID)
		if err != nil {
			return err
		}

		accountSession.SetCookie(c.Response())
		return c.JSON(http.StatusOK, StatusOK{Message: "success"})
	}
}

func LoginHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, StatusOK{Message: "success"})
	}
}

func LogoutHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		session.DeleteCookie(c.Response())

		return c.JSON(http.StatusOK, StatusOK{Message: "success"})
	}
}

func TokenEndpointHandler(provider infra.SessionProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := provider.GetSessionWithCookie(c.Request())

		if err != nil {
			return err
		}

		endpoints, err := session.GetIdentityEndpoints()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, endpoints)
	}
}
