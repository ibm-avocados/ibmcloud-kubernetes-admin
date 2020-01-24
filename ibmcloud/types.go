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
	Tags                          []Tag         `json:"tags"`
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

type Zones struct {
	Zones []Zone
}

type Zone struct {
	ID    string `json:"id"`
	Metro string `json:"metro"`
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

type Locations struct {
	Locations []Location
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
