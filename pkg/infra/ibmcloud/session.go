package ibmcloud

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ibm-avocados/ibmcloud-kubernetes-admin/pkg/infra"
)

const (
	accessTokenKey  = "access_token"
	refreshTokenKey = "refresh_token"
	expiration      = "expiration"
	cookiePath      = "/api/v1"
)

var endpoints *infra.IdentityEndpoints

type Session struct {
	Token *Token
}

// Token is used in every request that need authentication
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ImsUserID    int    `json:"ims_user_id"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Expiration   int    `json:"expiration"`
	Scope        string `json:"scope"`
}

func cacheIdentityEndpoints() error {
	if endpoints == nil {
		var err error
		endpoints, err = getIdentityEndpoints()
		if err != nil {
			return err
		}
	}
	return nil
}

// IAMAuthenticate uses the api key to authenticate and return a Session
func IAMAuthenticate(apikey string) (*Session, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		log.Println("error with cached data")
		return nil, err
	}
	token, err := getTokenFromIAM(endpoints.TokenEndpoint, apikey)

	if err != nil {
		log.Println("error with token data")
		return nil, err
	}
	return &Session{Token: token}, nil
}

func GetSession(r *http.Request) (*Session, error) {
	var accessToken string
	var refreshToken string
	var expirationTime int
	accessTokenVal, err := r.Cookie(accessTokenKey)
	if err != nil {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			return nil, err
		}
		parsedToken := strings.Split(bearerToken, " ")
		if len(parsedToken) != 2 {
			return nil, err
		}
		accessToken = parsedToken[1]
	} else {
		accessToken = accessTokenVal.Value
	}

	refreshTokenVal, err := r.Cookie(refreshTokenKey)
	if err != nil {
		refreshToken = r.Header.Get("X-Auth-Refresh-Token")
		if refreshToken == "" {
			return nil, err
		}
	} else {
		refreshToken = refreshTokenVal.Value
	}

	expirationValStr, err := r.Cookie(expiration)
	if err != nil {
		expiration := r.Header.Get("Expiration")
		if expiration == "" {
			expirationTime = 0
		}
	} else {
		expirationTime, err = strconv.Atoi(expirationValStr.Value)
		if err != nil {
			return nil, err
		}
	}

	session := &Session{
		Token: &Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Expiration:   expirationTime,
		},
	}

	return session.renewSession()
}

// GetIdentityEndpoints returns the list of endpoints for IBMCloud IAM
func (s *Session) GetIdentityEndpoints() (*infra.IdentityEndpoints, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

// GetZones returns an array of Zones
func (s *Session) GetZones(showFlavors, location string) ([]infra.Zone, error) {
	return getZones(showFlavors, location)
}

// GetVersions queries the IKS api for current kubernetes and openshift master versions
func (s *Session) GetVersions() (*infra.ClusterVersion, error) {
	return getVersions()
}

// GetLocations returns an array of Locations
func (s *Session) GetLocations() ([]infra.Location, error) {
	return getLocations()
}

// GetGeoLocations returns a list of location based on geo
func (s *Session) GetGeoLocations(geo string) ([]infra.Location, error) {
	return getGeoLocations(geo)
}

// GetMachineType returns a array of flavors of machines in a location based on queries
func (s *Session) GetMachineType(datacenter, serverType, os string, cpuLimit, memoryLimit int) ([]infra.MachineFlavor, error) {
	return getMachineTypes(datacenter, serverType, os, cpuLimit, memoryLimit)
}

func (s *Session) SetCookie(w http.ResponseWriter) {
	accessTokenCookie := &http.Cookie{Name: accessTokenKey, Value: s.Token.AccessToken, Path: cookiePath}
	http.SetCookie(w, accessTokenCookie)
	refreshTokenCookie := &http.Cookie{Name: refreshTokenKey, Value: s.Token.RefreshToken, Path: cookiePath}
	http.SetCookie(w, refreshTokenCookie)
	expirationStr := strconv.Itoa(s.Token.Expiration)

	expirationCookie := &http.Cookie{Name: expiration, Value: expirationStr, Path: cookiePath}
	http.SetCookie(w, expirationCookie)
}

func (s *Session) DeleteCookie(w http.ResponseWriter) {
	accessTokenCookie := &http.Cookie{Name: accessTokenKey, Value: "", Path: cookiePath}
	http.SetCookie(w, accessTokenCookie)

	refreshTokenCookie := &http.Cookie{Name: refreshTokenKey, Value: "", Path: cookiePath}
	http.SetCookie(w, refreshTokenCookie)

	expirationCookie := &http.Cookie{Name: expiration, Value: "", Path: cookiePath}
	http.SetCookie(w, expirationCookie)
}

// GetAccounts get a list of accounts for the current session
func (s *Session) GetAccounts() (*infra.Accounts, error) {
	return s.getAccountsWithEndpoint(nil)
}

// valid checks if session is expired or not
func (s *Session) valid() bool {
	now := time.Now().Unix()
	difference := int64(s.Token.Expiration) - now
	return difference > 600 // expires in 3600 second. keeping 10 min buffer
}

func (s *Session) getAccountsWithEndpoint(nextURL *string) (*infra.Accounts, error) {
	accounts, err := getAccounts(nextURL, s.Token.AccessToken)
	if err != nil {
		return nil, err
	}
	if accounts.NextURL != nil {
		nextAccounts, err := s.getAccountsWithEndpoint(accounts.NextURL)
		if err != nil {
			return nil, err
		}
		nextAccounts.Resources = append(nextAccounts.Resources, accounts.Resources...)
		return nextAccounts, nil
	}
	return accounts, nil
}

// GetCluster returns a Cluster with clusterID / name
func (s *Session) GetCluster(clusterID, resourceGroup string) (*infra.Cluster, error) {
	return getCluster(s.Token.AccessToken, clusterID, resourceGroup)
}

// GetClusters returns array of Clusters for current session
func (s *Session) GetClusters(location string) ([]*infra.Cluster, error) {
	return getClusters(s.Token.AccessToken, location)
}

// GetDatacenterVlan get an array of Vlan (public, private)
func (s *Session) GetDatacenterVlan(datacenter string) ([]infra.Vlan, error) {
	return getDatacenterVlan(s.Token.AccessToken, s.Token.RefreshToken, datacenter)
}

// GetAccessGroups returns accessgroups list
func (s *Session) GetAccessGroups(accountID string) (*infra.AccessGroups, error) {
	return getAccessGroups(s.Token.AccessToken, accountID)
}

// InviteUserToAccount returns UserInviteResponseList
func (s *Session) InviteUserToAccount(accountID, email string) (*infra.UserInviteList, error) {
	return inviteUserToAccount(s.Token.AccessToken, accountID, email)
}

// AddMemberToAccessGroup returns MemberAddResponseList
func (s *Session) AddMemberToAccessGroup(accessGroupID, iamID, memberType string) (*infra.MemberList, error) {
	return addMemberToAccessGroup(s.Token.AccessToken, accessGroupID, iamID, memberType)
}

// CreatePolicy returns Policy
func (s *Session) CreatePolicy(accountID, iamID, serviceName, serviceInstance, role string) (*infra.PolicyResponse, error) {
	return createPolicy(s.Token.AccessToken, accountID, iamID, serviceName, serviceInstance, role)
}

// IsMemberOfAccessGroup checks to see if a user is a member of the specified access group
func (s *Session) IsMemberOfAccessGroup(accessGroupID, iamID string) error {
	return isMemberOfAccessGroup(s.Token.AccessToken, accessGroupID, iamID)
}

// GetAccountResources return AccountResources
func (s *Session) GetAccountResources(accountID string) (*infra.AccountResources, error) {
	return getAccountResources(s.Token.AccessToken, accountID)
}

func (s *Session) GetUserInfo() (*infra.UserInfo, error) {
	return getUserInfo(endpoints.UserinfoEndpoint, s.Token.AccessToken)
}

func (s *Session) GetUserPreference(userID string) (*infra.User, error) {
	return getUserPreference(userID, s.Token.AccessToken)
}

// GetWorkers returns workers for a cluster
func (s *Session) GetWorkers(clusterID string) ([]infra.Worker, error) {
	return getClusterWorkers(s.Token.AccessToken, clusterID)
}

func (s *Session) CheckToken(apikey string) (*infra.ApiKeyDetails, error) {
	apiKeyDetails, err := checkToken(s.Token.AccessToken, apikey)
	if err != nil {
		return nil, err
	}
	return apiKeyDetails, nil
}

// BindAccountToToken upgrades session with account
func (s *Session) BindAccountToToken(accountID string) (infra.CloudSession, error) {
	session, err := bindAccountToToken(s.Token.RefreshToken, accountID)
	if err != nil {
		return nil, err
	}
	return session, err
}

// RenewSession renews session with refresh token
func (s *Session) renewSession() (*Session, error) {
	if s.valid() {
		return s, nil
	}
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, "")
	if err != nil {
		return nil, err
	}
	return &Session{Token: token}, nil
}

// DeleteCluster deletes cluster with a given id/name
func (s *Session) DeleteCluster(id, resourceGroup, deleteResources string) error {
	return deleteCluster(s.Token.AccessToken, id, resourceGroup, deleteResources)
}

// CreateCluster creates a cluster
func (s *Session) CreateCluster(request []byte) (*infra.CreateClusterResponse, error) {
	return createCluster(s.Token.AccessToken, request)
}

// SetTag sets the tag for a list of resources
func (s *Session) SetTag(updateTag []byte) (*infra.TagResult, error) {
	return setTags(s.Token.AccessToken, updateTag)
}

// SetClusterTag sets the tag for a cluster
func (s *Session) SetClusterTag(tag, clusterID, resourceID string) (*infra.TagResult, error) {
	return setClusterTags(s.Token.AccessToken, tag, clusterID, resourceID)
}

// DeleteTag deletes tag from list of resources
func (s *Session) DeleteTag(updateTag []byte) (*infra.TagResult, error) {
	return deleteTags(s.Token.AccessToken, updateTag)
}

// GetTags get tag infor for a cluster
func (s *Session) GetTags(clusterCRN string) (*infra.Tags, error) {
	return getTags(s.Token.AccessToken, clusterCRN)
}

// GetBillingData get billing info for the month for a specific cluster
func (s *Session) GetBillingData(accountID, clusterID, clusterCRN string) (string, error) {
	return getBillingData(s.Token.AccessToken, accountID, clusterID, clusterCRN)
}
