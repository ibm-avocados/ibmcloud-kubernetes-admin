package infra

import "net/http"

type CloudSession interface {
	SetCookie(w http.ResponseWriter)
	DeleteCookie(w http.ResponseWriter)
	GetIdentityEndpoints() (*IdentityEndpoints, error)
	GetZones(showFlavors, location string) ([]Zone, error)
	GetVersions() (*ClusterVersion, error)
	GetLocations() ([]Location, error)
	GetGeoLocations(geo string) ([]Location, error)
	GetMachineType(datacenter, serverType, os string, cpuLimit, memoryLimit int) ([]MachineFlavor, error)
	GetAccounts() (*Accounts, error)
	GetCluster(clusterID, resourceGroup string) (*Cluster, error)
	GetClusters(location string) ([]*Cluster, error)
	GetDatacenterVlan(datacenter string) ([]Vlan, error)
	GetAccessGroups(accountID string) (*AccessGroups, error)
	InviteUserToAccount(accountID, email string) (*UserInviteList, error)
	AddMemberToAccessGroup(accessGroupID, iamID, memberType string) (*MemberList, error)
	CreatePolicy(accountID, iamID, serviceName, serviceInstance, role string) (*PolicyResponse, error)
	IsMemberOfAccessGroup(accessGroupID, iamID string) error
	GetAccountResources(accountID string) (*AccountResources, error)
	GetUserInfo() (*UserInfo, error)
	GetUserPreference(userID string) (*User, error)
	GetWorkers(clusterID string) ([]Worker, error)
	CheckToken(apikey string) (*ApiKeyDetails, error)
	BindAccountToToken(accountID string) (CloudSession, error)
	CreateCluster(request []byte) (*CreateClusterResponse, error)
	DeleteCluster(id, resourceGroup, deleteResources string) error
	SetTag(updateTag []byte) (*TagResult, error)
	SetClusterTag(tag, clusterID, resourceID string) (*TagResult, error)
	DeleteTag(updateTag []byte) (*TagResult, error)
	GetTags(clusterCRN string) (*Tags, error)
	GetBillingData(accountID, clusterID, clusterCRN string) (string, error)
}

type SessionProvider interface {
	GetSessionWithToken(token *AuthToken) (CloudSession, error)
	GetSessionWithCookie(r *http.Request) (CloudSession, error)
}
