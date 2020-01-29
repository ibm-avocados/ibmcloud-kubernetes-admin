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

func GetIdentityEndpoints() (*IdentityEndpoints, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

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

func (s *Session) GetAccounts() (*Accounts, error) {
	return s.getAccountsWithEndpoint(nil)
}

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

func (s *Session) GetClusters(accountID, location string) ([]*Cluster, error) {
	return getClusters(s.Token.AccessToken, accountID, location)
}

func (s *Session) BindAccountToToken(accountID string) (*Session, error) {
	err := cacheIdentityEndpoints()
	if err != nil {
		return nil, err
	}
	token, err := upgradeToken(endpoints.TokenEndpoint, s.Token.RefreshToken, accountID)
	if err != nil {
		return nil, err
	}
	return &Session{Token: token}, nil
}

func (s *Session) DeleteCluster(id, resourceGroup, deleteResources string) error {
	return deleteCluster(s.Token.AccessToken, id, resourceGroup, deleteResources)
}

func (s *Session) SetTag(updateTag UpdateTag) (*TagResult, error) {
	return setTags(s.Token.AccessToken, updateTag)
}

func (s *Session) DeleteTag(updateTag UpdateTag) (*TagResult, error) {

	return deleteTags(s.Token.AccessToken, updateTag)
}
