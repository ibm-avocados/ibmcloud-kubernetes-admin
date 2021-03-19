package ibmcloud

import (
	"log"
	"time"
)

var endpoints *IdentityEndpoints

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

// GetIdentityEndpoints returns the list of endpoints for IBMCloud IAM
func GetIdentityEndpoints() (*IdentityEndpoints, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

// GetZones returns an array of Zones
func GetZones(showFlavors, location string) ([]Zone, error) {
	return getZones(showFlavors, location)
}

// GetVersions queries the IKS api for current kubernetes and openshift master versions
func GetVersions() (*ClusterVersion, error) {
	return getVersions()
}

// GetLocations returns an array of Locations
func GetLocations() ([]Location, error) {
	return getLocations()
}

// GetGeoLocations returns a list of location based on geo
func GetGeoLocations(geo string) ([]Location, error) {
	return getGeoLocations(geo)
}

// GetMachineType returns a array of flavors of machines in a location based on queries
func GetMachineType(datacenter, serverType, os string, cpuLimit, memoryLimit int) ([]MachineFlavor, error) {
	return getMachineTypes(datacenter, serverType, os, cpuLimit, memoryLimit)
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

// Authenticate uses the one time passcode to authenticate and return a Session
func Authenticate(otp string) (*Session, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		log.Println("error with cached data")
		return nil, err
	}
	token, err := getToken(endpoints.TokenEndpoint, otp)

	if err != nil {
		log.Println("error with token data")
		return nil, err
	}
	return &Session{Token: token}, nil
}

// GetAccounts get a list of accounts for the current session
func (s *Session) GetAccounts() (*Accounts, error) {
	return s.getAccountsWithEndpoint(nil)
}

// IsValid checks if session is expired or not
func (s *Session) IsValid() bool {
	now := time.Now().Unix()
	difference := int64(s.Token.Expiration) - now
	return difference > 100 // expires in 3600 second. keeping 100 second buffer
}

func (s *Session) getAccountsWithEndpoint(nextURL *string) (*Accounts, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, "")
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
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
func (s *Session) GetCluster(clusterID, resourceGroup string) (*Cluster, error) {
	return getCluster(s.Token.AccessToken, clusterID, resourceGroup)
}

// GetClusters returns array of Clusters for current session
func (s *Session) GetClusters(location string) ([]*Cluster, error) {
	return getClusters(s.Token.AccessToken, location)
}

// GetDatacenterVlan get an array of Vlan (public, private)
func (s *Session) GetDatacenterVlan(datacenter string) ([]Vlan, error) {
	return getDatacenterVlan(s.Token.AccessToken, s.Token.RefreshToken, datacenter)
}

// GetAccessGroups returns accessgroups list
func (s *Session) GetAccessGroups(accountID string) (*AccessGroups, error) {
	return getAccessGroups(s.Token.AccessToken, accountID)
}

// InviteUserToAccount returns UserInviteResponseList
func (s *Session) InviteUserToAccount(accountID, email string) (*UserInviteList, error) {
	return inviteUserToAccount(s.Token.AccessToken, accountID, email)
}

// AddMemberToAccessGroup returns MemberAddResponseList
func (s *Session) AddMemberToAccessGroup(accessGroupID, iamID, memberType string) (*MemberList, error) {
	return addMemberToAccessGroup(s.Token.AccessToken, accessGroupID, iamID, memberType)
}

// CreatePolicy returns Policy
func (s *Session) CreatePolicy(accountID, iamID, serviceName, serviceInstance, role string) (*PolicyResponse, error) {
	return createPolicy(s.Token.AccessToken, accountID, iamID, serviceName, serviceInstance, role)
}

// IsMemberOfAccessGroup checks to see if a user is a member of the specified access group
func (s *Session) IsMemberOfAccessGroup(accessGroupID, iamID string) error {
	return isMemberOfAccessGroup(s.Token.AccessToken, accessGroupID, iamID)
}

// GetAccountResources return AccountResources
func (s *Session) GetAccountResources(accountID string) (*AccountResources, error) {
	return getAccountResources(s.Token.AccessToken, accountID)
}

func (s *Session) GetUserInfo() (*UserInfo, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	if !s.IsValid() {
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, "")
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}

	return getUserInfo(endpoints.UserinfoEndpoint, s.Token.AccessToken)
}

func (s *Session) GetUserPreference(userID string) (*User, error) {
	if !s.IsValid() {
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, "")
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return getUserPreference(userID, s.Token.AccessToken)
}

// GetWorkers returns workers for a cluster
func (s *Session) GetWorkers(clusterID string) ([]Worker, error) {
	return getClusterWorkers(s.Token.AccessToken, clusterID)
}

func bindAccountToToken(refreshToken, accountID string) (*Session, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	token, err := upgradeToken(endpoints.TokenEndpoint, refreshToken, accountID)
	if err != nil {
		return nil, err
	}
	return &Session{Token: token}, nil
}

func (s *Session) CheckToken(apikey string) (*ApiKeyDetails, error) {
	apiKeyDetails, err := checkToken(s.Token.AccessToken, apikey)
	if err != nil {
		return nil, err
	}
	return apiKeyDetails, nil
}

// BindAccountToToken upgrades session with account
func (s *Session) BindAccountToToken(accountID string) (*Session, error) {
	session, err := bindAccountToToken(s.Token.RefreshToken, accountID)
	if err != nil {
		return nil, err
	}
	return session, err
}

// RenewSession renews session with refresh token
func (s *Session) RenewSession() (*Session, error) {
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
func (s *Session) CreateCluster(request CreateClusterRequest) (*CreateClusterResponse, error) {
	return createCluster(s.Token.AccessToken, request)
}

// SetTag sets the tag for a list of resources
func (s *Session) SetTag(updateTag UpdateTag) (*TagResult, error) {
	return setTags(s.Token.AccessToken, updateTag)
}

// SetClusterTag sets the tag for a cluster
func (s *Session) SetClusterTag(tag, clusterID, resourceID string) (*TagResult, error) {
	return setClusterTags(s.Token.AccessToken, tag, clusterID, resourceID)
}

// DeleteTag deletes tag from list of resources
func (s *Session) DeleteTag(updateTag UpdateTag) (*TagResult, error) {
	return deleteTags(s.Token.AccessToken, updateTag)
}

// GetTags get tag infor for a cluster
func (s *Session) GetTags(clusterCRN string) (*Tags, error) {
	return getTags(s.Token.AccessToken, clusterCRN)
}

// GetBillingData get billing info for the month for a specific cluster
func (s *Session) GetBillingData(accountID, clusterID, clusterCRN string) (string, error) {
	return getBillingData(s.Token.AccessToken, accountID, clusterID, clusterCRN)
}

// CLOUDANT RELATED METHODS

// SetAPIKey sets the api key for an account
func (s *Session) SetAPIKey(apiKey, accountID string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return SetAPIKey(apiKey, accountID)
}

// CheckAPIKey checks if the APIKey is valid
func (s *Session) CheckAPIKey(accountID string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return CheckAPIKey(accountID)
}

// GetAPIKey returns the APIKey
func (s *Session) GetAPIKey(accountID string) (string, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return "", err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return GetAPIKey(accountID)
}

// UpdateAPIKey update the existing api key
func (s *Session) UpdateAPIKey(apiKey, accountID string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return UpdateAPIKey(apiKey, accountID)
}

// DeleteAPIKey delete the api key
func (s *Session) DeleteAPIKey(accountID string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return DeleteAPIKey(accountID)
}

func (s *Session) GetAccountMetaData(accountID string) (*AccountMetaData, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return GetAccountMetaData(accountID)
}

func (s *Session) CreateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return CreateAccountMetadata(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken)
}

func (s *Session) UpdateAccountMetaData(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return UpdateAccountMetadata(accountID, org, space, region, accessGroup, issueRepo, grantClusterRepo, githubUser, githubToken)
}

func (s *Session) GetDocument(accountID string) ([]Schedule, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return GetDocument(accountID)
}

func (s *Session) GetAllDocument(accountID string) ([]interface{}, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return GetAllDocument(accountID)
}

func (s *Session) CreateDocument(accountID string, data interface{}) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return CreateDocument(accountID, data)
}

func (s *Session) DeleteDocument(accountID, id, rev string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}

	return DeleteDocument(accountID, id, rev)
}

func (s *Session) UpdateDocument(accountID, id, rev string, data interface{}) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return UpdateDocument(accountID, id, rev, data)
}

func (s *Session) CreateAdminEmails(accountID string, emails ...string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return CreateAdminEmails(accountID, emails...)
}

func (s *Session) RemoveAdminEmails(accountID string, emails ...string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return RemoveAdminEmails(accountID, emails...)
}

func (s *Session) AddAdminEmails(accountID string, emails ...string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return AddAdminEmails(accountID, emails...)
}

func (s *Session) DeleteAdminEmails(accountID string) error {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return DeleteAdminEmails(accountID)
}

func (s *Session) GetAccountAdminEmails(accountID string) ([]string, error) {
	if !s.IsValid() {
		log.Println("Access token expired.")
		token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
		if err != nil {
			return nil, err
		}
		log.Println("Token Refreshed.")
		s.Token = token
	}
	return GetAccountAdminEmails(accountID)
}
