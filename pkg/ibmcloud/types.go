package ibmcloud

import "time"

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

type IdentityEndpoints struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	PasscodeEndpoint                  string   `json:"passcode_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
}

//ACCOUNT RELATED TYPES
type Accounts struct {
	NextURL      *string   `json:"next_url"`
	TotalResults int       `json:"total_results"`
	Resources    []Account `json:"resources"`
}
type Metadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type TermsAndConditions struct {
	Required  bool      `json:"required"`
	Accepted  bool      `json:"accepted"`
	Timestamp time.Time `json:"timestamp"`
}
type OrganizationsRegion struct {
	GUID   string `json:"guid"`
	Region string `json:"region"`
}
type Linkages struct {
	Origin string `json:"origin"`
	State  string `json:"state"`
}
type PaymentMethod struct {
	Type           string      `json:"type"`
	Started        time.Time   `json:"started"`
	Ended          string      `json:"ended"`
	CurrencyCode   string      `json:"currencyCode"`
	AnniversaryDay interface{} `json:"anniversaryDay"`
}
type History struct {
	Type               string    `json:"type"`
	State              string    `json:"state"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	CurrencyCode       string    `json:"currencyCode"`
	CountryCode        string    `json:"countryCode"`
	BillingCountryCode string    `json:"billingCountryCode"`
	BillingSystem      string    `json:"billingSystem"`
}
type BluemixSubscriptions struct {
	Type                  string        `json:"type"`
	State                 string        `json:"state"`
	PaymentMethod         PaymentMethod `json:"payment_method"`
	SubscriptionID        string        `json:"subscription_id"`
	PartNumber            string        `json:"part_number"`
	SubscriptionTags      []interface{} `json:"subscriptionTags"`
	PaygPendingTimestamp  time.Time     `json:"payg_pending_timestamp"`
	History               []History     `json:"history"`
	CurrentStateTimestamp time.Time     `json:"current_state_timestamp"`
	SoftlayerAccountID    string        `json:"softlayer_account_id"`
	BillingSystem         string        `json:"billing_system"`
}
type Entity struct {
	Name                 string                 `json:"name"`
	Type                 string                 `json:"type"`
	State                string                 `json:"state"`
	Owner                string                 `json:"owner"`
	OwnerUserid          string                 `json:"owner_userid"`
	OwnerUniqueID        string                 `json:"owner_unique_id"`
	OwnerIamID           string                 `json:"owner_iam_id"`
	CustomerID           string                 `json:"customer_id"`
	CountryCode          string                 `json:"country_code"`
	CurrencyCode         string                 `json:"currency_code"`
	BillingCountryCode   string                 `json:"billing_country_code"`
	IsIBMer              bool                   `json:"isIBMer"`
	Promotions           []string               `json:"promotions"`
	Quota                string                 `json:"quota"`
	TermsAndConditions   TermsAndConditions     `json:"terms_and_conditions"`
	Tags                 []interface{}          `json:"tags"`
	TeamDirectoryEnabled bool                   `json:"team_directory_enabled"`
	OrganizationsRegion  []OrganizationsRegion  `json:"organizations_region"`
	Linkages             []Linkages             `json:"linkages"`
	BluemixSubscriptions []BluemixSubscriptions `json:"bluemix_subscriptions"`
	SubscriptionID       string                 `json:"subscription_id"`
	ConfigurationID      string                 `json:"configuration_id"`
	Onboarded            int                    `json:"onboarded"`
}
type Account struct {
	Metadata Metadata `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

// CLUSTER RELATED TYPES
type Cluster struct {
	Location                      string        `json:"location"`
	DataCenter                    string        `json:"dataCenter"`
	MultiAzCapable                bool          `json:"multiAzCapable"`
	Vlans                         []string      `json:"vlans"`
	WorkerVlans                   []string      `json:"worker_vlans"`
	WorkerZones                   []string      `json:"workerZones"`
	ID                            string        `json:"id"`
	Name                          string        `json:"name"`
	Region                        string        `json:"region"`
	ResourceGroup                 string        `json:"resourceGroup"`
	ResourceGroupName             string        `json:"resourceGroupName"`
	ServerURL                     string        `json:"serverURL"`
	State                         string        `json:"state"`
	CreatedDate                   string        `json:"createdDate"`
	ModifiedDate                  string        `json:"modifiedDate"`
	WorkerCount                   int           `json:"workerCount"`
	IsPaid                        bool          `json:"isPaid"`
	MasterKubeVersion             string        `json:"masterKubeVersion"`
	TargetVersion                 string        `json:"targetVersion"`
	IngressHostname               string        `json:"ingressHostname"`
	IngressSecretName             string        `json:"ingressSecretName"`
	OwnerEmail                    string        `json:"ownerEmail"`
	LogOrg                        string        `json:"logOrg"`
	LogOrgName                    string        `json:"logOrgName"`
	LogSpace                      string        `json:"logSpace"`
	LogSpaceName                  string        `json:"logSpaceName"`
	APIUser                       string        `json:"apiUser"`
	MonitoringURL                 string        `json:"monitoringURL"`
	Addons                        []interface{} `json:"addons"`
	VersionEOS                    string        `json:"versionEOS"`
	DisableAutoUpdate             bool          `json:"disableAutoUpdate"`
	EtcdPort                      string        `json:"etcdPort"`
	MasterStatus                  string        `json:"masterStatus"`
	MasterStatusModifiedDate      string        `json:"masterStatusModifiedDate"`
	MasterHealth                  string        `json:"masterHealth"`
	MasterState                   string        `json:"masterState"`
	KeyProtectEnabled             bool          `json:"keyProtectEnabled"`
	PullSecretApplied             bool          `json:"pullSecretApplied"`
	Crn                           string        `json:"crn"`
	PrivateServiceEndpointEnabled bool          `json:"privateServiceEndpointEnabled"`
	PrivateServiceEndpointURL     string        `json:"privateServiceEndpointURL"`
	PublicServiceEndpointEnabled  bool          `json:"publicServiceEndpointEnabled"`
	PublicServiceEndpointURL      string        `json:"publicServiceEndpointURL"`
	PodSubnet                     string        `json:"podSubnet"`
	ServiceSubnet                 string        `json:"serviceSubnet"`
	Type                          string        `json:"type"`
	Tags                          []string      `json:"tags"`
	Workers                       []Worker      `json:"workers"`
	Cost                          string        `json:"cost"`
}

type Worker struct {
	PrivateVlan      string `json:"privateVlan"`
	PublicVlan       string `json:"publicVlan"`
	PrivateIP        string `json:"privateIP"`
	PublicIP         string `json:"publicIP"`
	MachineType      string `json:"machineType"`
	Location         string `json:"location"`
	ID               string `json:"id"`
	State            string `json:"state"`
	Status           string `json:"status"`
	StatusDate       string `json:"statusDate"`
	StatusDetails    string `json:"statusDetails"`
	ErrorMessage     string `json:"errorMessage"`
	ErrorMessageDate string `json:"errorMessageDate"`
	Isolation        string `json:"isolation"`
	KubeVersion      string `json:"kubeVersion"`
	TargetVersion    string `json:"targetVersion"`
	ReasonForDelete  string `json:"reasonForDelete"`
	VersionEOS       string `json:"versionEOS"`
	MasterVersionEOS string `json:"masterVersionEOS"`
	TrustedStatus    string `json:"trustedStatus"`
	Poolid           string `json:"poolid"`
	PoolName         string `json:"poolName"`
	PendingOperation string `json:"pendingOperation"`
}

type UpdateTag struct {
	TagName   string     `json:"tag_name"`
	Resources []Resource `json:"resources"`
}
type Resource struct {
	ResourceID string `json:"resource_id"`
}

type TagResult struct {
	Results []Results `json:"results"`
}
type Results struct {
	ResourceID string      `json:"resource_id"`
	IsError    string      `json:"isError"`
	ISError    interface{} `json:"is_error"`
}

type Tags struct {
	TotalCount int   `json:"total_count"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
	Items      []Tag `json:"items"`
}
type Tag struct {
	Name string `json:"name"`
}

type Location struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	Metro          string `json:"metro,omitempty"`
	MultizoneMetro string `json:"multizone_metro,omitempty"`
	Country        string `json:"country,omitempty"`
	Geography      string `json:"geography,omitempty"`
	DisplayName    string `json:"display_name"`
}

type LocationWithClusterCount struct {
	ID             string `json:"id"`
	ClusterCount   int    `json:"clusterCount"`
	Name           string `json:"name"`
	Kind           string `json:"kind"`
	Metro          string `json:"metro,omitempty"`
	MultizoneMetro string `json:"multizone_metro,omitempty"`
	Country        string `json:"country,omitempty"`
	Geography      string `json:"geography,omitempty"`
	DisplayName    string `json:"display_name"`
}

type Zone struct {
	ID      string    `json:"id"`
	Metro   string    `json:"metro"`
	Flavors []Flavors `json:"flavors"`
}
type Flavors struct {
	Name                      string `json:"name"`
	Provider                  string `json:"provider"`
	Memory                    string `json:"memory"`
	NetworkSpeed              string `json:"networkSpeed"`
	Cores                     string `json:"cores"`
	Os                        string `json:"os"`
	ServerType                string `json:"serverType"`
	Storage                   string `json:"storage"`
	SecondaryStorage          string `json:"secondaryStorage"`
	SecondaryStorageEncrypted bool   `json:"secondaryStorageEncrypted"`
	Deprecated                bool   `json:"deprecated"`
	CorrespondingMachineType  string `json:"correspondingMachineType"`
	IsTrusted                 bool   `json:"isTrusted"`
	Gpus                      string `json:"gpus"`
}

type ClusterVersion struct {
	Kubernetes []Kubernetes `json:"kubernetes"`
	Openshift  []Openshift  `json:"openshift"`
}
type Kubernetes struct {
	Major        int    `json:"major"`
	Minor        int    `json:"minor"`
	Patch        int    `json:"patch"`
	Default      bool   `json:"default"`
	EndOfService string `json:"end_of_service"`
}
type Openshift struct {
	Major        int    `json:"major"`
	Minor        int    `json:"minor"`
	Patch        int    `json:"patch"`
	Default      bool   `json:"default"`
	EndOfService string `json:"end_of_service"`
}

type ResourceGroups struct {
	Resources []Group `json:"resources"`
}
type Group struct {
	ID                string        `json:"id"`
	Crn               string        `json:"crn"`
	AccountID         string        `json:"account_id"`
	Name              string        `json:"name"`
	State             string        `json:"state"`
	Default           bool          `json:"default"`
	EnableReclamation bool          `json:"enable_reclamation"`
	QuotaID           string        `json:"quota_id"`
	QuotaURL          string        `json:"quota_url"`
	PaymentMethodsURL string        `json:"payment_methods_url"`
	ResourceLinkages  []interface{} `json:"resource_linkages"`
	TeamsURL          string        `json:"teams_url"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// USER RELATED TYPES
type Users struct {
	TotalResults int    `json:"total_results"`
	Limit        int    `json:"limit"`
	FirstURL     string `json:"first_url"`
	NextURL      string `json:"next_url"`
	Resources    []User `json:"resources"`
}
type User struct {
	ID             string `json:"id"`
	IamID          string `json:"iam_id"`
	Realm          string `json:"realm"`
	UserID         string `json:"user_id"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	State          string `json:"state"`
	Email          string `json:"email"`
	Phonenumber    string `json:"phonenumber"`
	Altphonenumber string `json:"altphonenumber"`
	Photo          string `json:"photo"`
	AccountID      string `json:"account_id"`
}

type UserInfo struct {
	Active     bool     `json:"active"`
	RealmID    string   `json:"realmId"`
	Identifier string   `json:"identifier"`
	IamID      string   `json:"iam_id"`
	GivenName  string   `json:"given_name"`
	FamilyName string   `json:"family_name"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Sub        string   `json:"sub"`
	Account    Account  `json:"account"`
	Iat        int      `json:"iat"`
	Exp        int      `json:"exp"`
	Iss        string   `json:"iss"`
	GrantType  string   `json:"grant_type"`
	ClientID   string   `json:"client_id"`
	Scope      string   `json:"scope"`
	Acr        int      `json:"acr"`
	Amr        []string `json:"amr"`
}

type ErrorMessage struct {
	ErrorDescription string  `json:"error_description"`
	Trace            string  `json:"trace"`
	Error            []Error `json:"error"`
	Errors           []Error `json:"errors"`
}

type Target struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
type Error struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Target   Target `json:"target"`
}

// BILLING TYPES
type ResourceUsage struct {
	Limit     int         `json:"limit"`
	Count     int         `json:"count"`
	First     URL         `json:"first"`
	Resources []Resources `json:"resources"`
	Next      URL         `json:"next"`
}

type Price struct {
	UnitQuantity string      `json:"unitQuantity"`
	TierModel    string      `json:"tier_model"`
	Price        interface{} `json:"price"`
	QuantityTier string      `json:"quantity_tier"`
}
type Discounts struct {
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Ref         string      `json:"ref"`
	Discount    interface{} `json:"discount"`
}
type Usage struct {
	Metric           string      `json:"metric"`
	Unit             string      `json:"unit"`
	Price            []Price     `json:"price"`
	Quantity         float64     `json:"quantity"`
	Cost             float64     `json:"cost"`
	RatedCost        float64     `json:"rated_cost"`
	RateableQuantity float64     `json:"rateable_quantity"`
	Discounts        []Discounts `json:"discounts"`
}
type Resources struct {
	AccountID          string  `json:"account_id"`
	ResourceInstanceID string  `json:"resource_instance_id"`
	ResourceGroupID    string  `json:"resource_group_id"`
	Month              string  `json:"month"`
	PricingCountry     string  `json:"pricing_country"`
	BillingCountry     string  `json:"billing_country"`
	CurrencyCode       string  `json:"currency_code"`
	PlanID             string  `json:"plan_id"`
	ResourceID         string  `json:"resource_id"`
	Billable           bool    `json:"billable"`
	PricingPlanID      string  `json:"pricing_plan_id"`
	Region             string  `json:"region"`
	Usage              []Usage `json:"usage"`
}
type URL struct {
	Href   string `json:"href"`
	Offset string `json:"offset"`
}

type ApiKey struct {
	ID     string `json:"_id"`
	Rev    string `json:"_rev"`
	APIKey string `json:"apiKey"`
}

type AccountAdminEmails struct {
	ID     string   `json:"_id"`
	Rev    string   `json:"_rev"`
	Emails []string `json:"emails"`
}

type AccountMetaData struct {
	ID               string `json:"_id"`
	Rev              string `json:"_rev"`
	Org              string `json:"org"`
	Space            string `json:"space"`
	Region           string `json:"region"`
	AccessGroup      string `json:"accessGroup"`
	IssueRepo        string `json:"issueRepo"`
	GrantClusterRepo string `json:"grantClusterRepo"`
	GithubUser       string `json:"githubUser"`
	GithubToken      string `json:"githubToken"`
}

// type Schedule struct {
// 	ID              string                 `json:"_id" mapstructure:"_id"`
// 	Rev             string                 `json:"_rev" mapstructure:"_rev"`
// 	CreateAt        int                    `json:"createAt"`
// 	DestroyAt       int                    `json:"destroyAt"`
// 	Status          string                 `json:"status"`
// 	Tags            string                 `json:"tags"`
// 	Count           int                    `json:"count"`
// 	ClusterRequests []CreateClusterRequest `json:"ClusterRequests"`
// }

type CreateClusterRequest struct {
	ClusterRequest ClusterRequest `json:"clusterRequest"`
	ResourceGroup  string         `json:"resourceGroup"`
}

type ScheduleClusterRequest struct {
	ScheduleRequest ScheduleRequest `json:"clusterRequest"`
	ResourceGroup   string          `json:"resourceGroup"`
}

type Schedule struct {
	ID                string                 `json:"_id" mapstructure:"_id"`
	Rev               string                 `json:"_rev" mapstructure:"_rev"`
	CreateAt          int                    `json:"createAt"`
	DestroyAt         int                    `json:"destroyAt"`
	Status            string                 `json:"status"`
	Tags              string                 `json:"tags"`
	Count             string                 `json:"count"`
	UserCount         string                 `json:"userCount"`
	ScheduleRequest   ScheduleClusterRequest `json:"scheduleRequest"`
	Clusters          []string               `json:"clusters"`
	NotifyEmails      []string               `json:"notifyEmails"`
	EventName         string                 `json:"eventName"`
	Password          string                 `json:"password"`
	ResourceGroupName string                 `json:"resourceGroupName"`
	GithubIssueNumber string                 `json:"githubIssueNumber"`
	IsWorkshop        bool                   `json:"isWorkshop"`
}

type Vlan struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Properties VlanProperties `json:"properties"`
}
type VlanProperties struct {
	Name                       string `json:"name"`
	Note                       string `json:"note"`
	PrimaryRouter              string `json:"primary_router"`
	VlanNumber                 string `json:"vlan_number"`
	VlanType                   string `json:"vlan_type"`
	Location                   string `json:"location"`
	LocalDiskStorageCapability string `json:"local_disk_storage_capability"`
	SanStorageCapability       string `json:"san_storage_capability"`
}

type MachineFlavor struct {
	Name                      string `json:"name"`
	Provider                  string `json:"provider"`
	Memory                    string `json:"memory"`
	NetworkSpeed              string `json:"networkSpeed"`
	Cores                     string `json:"cores"`
	Os                        string `json:"os"`
	ServerType                string `json:"serverType"`
	Storage                   string `json:"storage"`
	SecondaryStorage          string `json:"secondaryStorage"`
	SecondaryStorageEncrypted bool   `json:"secondaryStorageEncrypted"`
	Deprecated                bool   `json:"deprecated"`
	CorrespondingMachineType  string `json:"correspondingMachineType"`
	IsTrusted                 bool   `json:"isTrusted"`
	Gpus                      string `json:"gpus"`
	OcpUnsupported            bool   `json:"ocp_unsupported"`
}

type AccountResources struct {
	Resources []AccountResource `json:"resources"`
}
type AccountResource struct {
	ID                string        `json:"id"`
	Crn               string        `json:"crn"`
	AccountID         string        `json:"account_id"`
	Name              string        `json:"name"`
	State             string        `json:"state"`
	Default           bool          `json:"default"`
	EnableReclamation bool          `json:"enable_reclamation"`
	QuotaID           string        `json:"quota_id"`
	QuotaURL          string        `json:"quota_url"`
	PaymentMethodsURL string        `json:"payment_methods_url"`
	ResourceLinkages  []interface{} `json:"resource_linkages"`
	TeamsURL          string        `json:"teams_url"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

type CreateClusterResponse struct {
	ID                string            `json:"id"`
	NonCriticalErrors NonCriticalErrors `json:"non_critical_errors"`
}
type ErrorItems struct {
	Code             string `json:"code"`
	Description      string `json:"description"`
	RecoveryCLI      string `json:"recoveryCLI"`
	RecoveryUI       string `json:"recoveryUI"`
	TerseDescription string `json:"terseDescription"`
	Type             string `json:"type"`
}
type NonCriticalErrors struct {
	IncidentID string       `json:"incidentID"`
	Items      []ErrorItems `json:"items"`
}

type ScheduleRequest struct {
	DataCenters                  []string `json:"dataCenters"`
	DefaultWorkerPoolEntitlement string   `json:"defaultWorkerPoolEntitlement"`
	DefaultWorkerPoolName        string   `json:"defaultWorkerPoolName"`
	DisableAutoUpdate            bool     `json:"disableAutoUpdate"`
	DiskEncryption               bool     `json:"diskEncryption"`
	GatewayEnabled               bool     `json:"gatewayEnabled"`
	Isolation                    string   `json:"isolation"`
	MachineType                  string   `json:"machineType"`
	MasterVersion                string   `json:"masterVersion"`
	Name                         string   `json:"name"`
	NoSubnet                     bool     `json:"noSubnet"`
	PodSubnet                    string   `json:"podSubnet"`
	Prefix                       string   `json:"prefix"`
	PrivateSeviceEndpoint        bool     `json:"privateSeviceEndpoint"`
	PublicServiceEndpoint        bool     `json:"publicServiceEndpoint"`
	ServiceSubnet                string   `json:"serviceSubnet"`
	SkipPermPrecheck             bool     `json:"skipPermPrecheck"`
	WorkerNum                    int      `json:"workerNum"`
}

type ClusterRequest struct {
	DataCenter                   string `json:"dataCenter"`
	DefaultWorkerPoolEntitlement string `json:"defaultWorkerPoolEntitlement"`
	DefaultWorkerPoolName        string `json:"defaultWorkerPoolName"`
	DisableAutoUpdate            bool   `json:"disableAutoUpdate"`
	DiskEncryption               bool   `json:"diskEncryption"`
	GatewayEnabled               bool   `json:"gatewayEnabled"`
	Isolation                    string `json:"isolation"`
	MachineType                  string `json:"machineType"`
	MasterVersion                string `json:"masterVersion"`
	Name                         string `json:"name"`
	NoSubnet                     bool   `json:"noSubnet"`
	PodSubnet                    string `json:"podSubnet"`
	Prefix                       string `json:"prefix"`
	PrivateSeviceEndpoint        bool   `json:"privateSeviceEndpoint"`
	PrivateVlan                  string `json:"privateVlan"`
	PublicServiceEndpoint        bool   `json:"publicServiceEndpoint"`
	PublicVlan                   string `json:"publicVlan"`
	ServiceSubnet                string `json:"serviceSubnet"`
	SkipPermPrecheck             bool   `json:"skipPermPrecheck"`
	WorkerNum                    int    `json:"workerNum"`
}

type ApiKeyDetails struct {
	ID          string `json:"id"`
	EntityTag   string `json:"entity_tag"`
	Crn         string `json:"crn"`
	Locked      bool   `json:"locked"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	ModifiedAt  string `json:"modified_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IamID       string `json:"iam_id"`
	AccountID   string `json:"account_id"`
}

type AccessGroups struct {
	Limit      int      `json:"limit"`
	Offset     int      `json:"offset"`
	TotalCount int      `json:"total_count"`
	First      First    `json:"first"`
	Last       Last     `json:"last"`
	Groups     []Groups `json:"groups"`
}
type First struct {
	Href string `json:"href"`
}
type Last struct {
	Href string `json:"href"`
}
type Groups struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedByID      string    `json:"createdById"`
	LastModifiedAt   time.Time `json:"lastModifiedAt"`
	LastModifiedByID string    `json:"lastModifiedById"`
	Href             string    `json:"href"`
}

type GithubIssueComment struct {
	IssueNumber    string                    `json:"issueNumber"`
	EventName      string                    `json:"eventName"`
	Password       string                    `json:"password"`
	AccountID      string                    `json:"accountID"`
	GithubUser     string                    `json:"githubUser"`
	GithubToken    string                    `json:"githubToken"`
	ClusterRequest GithubIssueClusterRequest `json:"clusterRequest"`
}

type GithubIssueClusterRequest struct {
	Count      int    `json:"count"`
	Type       string `json:"type"`
	ErrorCount int    `json:"errorCount"`
	Regions    string `json:"regions"`
}

type UserInviteList struct {
	Users     []UserInvite `json:"users,omitempty"`
	Resources []UserInvite `json:"resources,omitempty"`
}

// type UserInvite struct {
// 	Email       string `json:"email"`
// 	AccountRole string `json:"account_role,omitempty"`
// }

// type UserInviteResponseList struct {
// 	Resources []UserInviteResponse `json:"resources"`
// }

type UserInvite struct {
	ID          string `json:"id,omitempty"`
	AccountRole string `json:"account_role,omitempty"`
	Email       string `json:"email"`
	State       string `json:"state,omitempty"`
}

type MemberList struct {
	Members []Member `json:"members"`
}

type Member struct {
	IamID       string     `json:"iam_id"`
	Type        string     `json:"type"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	CreatedByID string     `json:"created_by_id,omitempty"`
	StatusCode  int        `json:"status_code,omitempty"`
}

type PolicyResponse struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Subjects    []Subjects  `json:"subjects"`
	Roles       []Roles     `json:"roles"`
	Resources   []PolicyResources `json:"resources"`
	Href             string          `json:"href",omitempty`
	CreatedAt        time.Time       `json:"created_at,omitempty"`
	CreatedByID      string          `json:"created_by_id,omitempty"`
	LastModifiedAt   time.Time       `json:"last_modified_at,omitempty"`
	LastModifiedByID string          `json:"last_modified_by_id,omitempty"`
}

type Policy struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Subjects    []Subjects  `json:"subjects"`
	Roles       []Roles     `json:"roles"`
	Resources   []PolicyResources `json:"resources"`
}
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Subjects struct {
	Attributes []Attribute `json:"attributes"`
}
type Roles struct {
	RoleID string `json:"role_id"`
}
type PolicyResources struct {
	Attributes []Attribute `json:"attributes"`
}
