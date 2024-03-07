package entity

import (
	"fmt"
	"go-openstackclient/consts"
)

type Action string
type RuleType string
type CompareType string
type LBMethod string
type Protocol string
type TLSVersion string

const (
	ActionRedirectPrefix Action = "REDIRECT_PREFIX"
	ActionRedirectToPool Action = "REDIRECT_TO_POOL"
	ActionRedirectToURL  Action = "REDIRECT_TO_URL"
	ActionReject         Action = "REJECT"

	TypeCookie   RuleType = "COOKIE"
	TypeFileType RuleType = "FILE_TYPE"
	TypeHeader   RuleType = "HEADER"
	TypeHostName RuleType = "HOST_NAME"
	TypePath     RuleType = "PATH"

	CompareTypeContains  CompareType = "CONTAINS"
	CompareTypeEndWith   CompareType = "ENDS_WITH"
	CompareTypeEqual     CompareType = "EQUAL_TO"
	CompareTypeRegex     CompareType = "REGEX"
	CompareTypeStartWith CompareType = "STARTS_WITH"

	LBMethodRoundRobin       LBMethod = "ROUND_ROBIN"
	LBMethodLeastConnections LBMethod = "LEAST_CONNECTIONS"
	LBMethodSourceIp         LBMethod = "SOURCE_IP"
	LBMethodSourceIpPort     LBMethod = "SOURCE_IP_PORT"

	ProtocolPROXY Protocol = "PROXY"
	ProtocolHTTP  Protocol = "HTTP"
	ProtocolTCP   Protocol = "TCP"
	PING          string = "PING"
	ProtocolHTTPS Protocol = "HTTPS"
	// Protocol PROXYV2 requires octavia microversion 2.22
	ProtocolPROXYV2 Protocol = "PROXYV2"
	// Protocol SCTP requires octavia microversion 2.23
	ProtocolSCTP Protocol = "SCTP"

	TLSVersionSSLv3   TLSVersion = "SSLv3"
	TLSVersionTLSv1   TLSVersion = "TLSv1"
	TLSVersionTLSv1_1 TLSVersion = "TLSv1.1"
	TLSVersionTLSv1_2 TLSVersion = "TLSv1.2"
	TLSVersionTLSv1_3 TLSVersion = "TLSv1.3"
)

// CreateLoadbalancerOpts is the common options struct used in this package's Create
// operation.
type CreateLoadbalancerOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// Human-readable description for the Loadbalancer.
	Description string `json:"description,omitempty"`

	// Providing a neutron port ID for the vip_port_id tells Octavia to use this
	// port for the VIP. If the port has more than one subnet you must specify
	// either the vip_subnet_id or vip_address to clarify which address should
	// be used for the VIP.
	VipPortID string `json:"vip_port_id,omitempty"`

	// The subnet on which to allocate the Loadbalancer's address. A project can
	// only create Loadbalancers on networks authorized by policy (e.g. networks
	// that belong to them or networks that are shared).
	VipSubnetID string `json:"vip_subnet_id,omitempty"`

	// The network on which to allocate the Loadbalancer's address. A tenant can
	// only create Loadbalancers on networks authorized by policy (e.g. networks
	// that belong to them or networks that are shared).
	VipNetworkID string `json:"vip_network_id,omitempty"`

	// ProjectID is the UUID of the project who owns the Loadbalancer.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// The IP address of the Loadbalancer.
	VipAddress string `json:"vip_address,omitempty"`

	// The ID of the QoS Policy which will apply to the Virtual IP
	VipQosPolicyID string `json:"vip_qos_policy_id,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The UUID of a flavor.
	FlavorID string `json:"flavor_id,omitempty"`

	// The name of an Octavia availability zone.
	// Requires Octavia API version 2.14 or later.
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// The name of the provider.
	Provider string `json:"provider,omitempty"`

	// Listeners is a slice of listeners.CreateOpts which allows a set
	// of listeners to be created at the same time the Loadbalancer is created.
	//
	// This is only possible to use when creating a fully populated
	// load balancer.
	Listeners []CreateListenerOpts `json:"listeners,omitempty"`

	// Pools is a slice of pools.CreateOpts which allows a set of pools
	// to be created at the same time the Loadbalancer is created.
	//
	// This is only possible to use when creating a fully populated
	// load balancer.
	Pools []CreatePoolOpts `json:"pools,omitempty"`

	// Tags is a set of resource tags.
	Tags []string `json:"tags,omitempty"`
}

// UpdateLoadbalancerOpts is the common options struct used in this package's Update
// operation.
type UpdateLoadbalancerOpts struct {
	// Human-readable name for the Loadbalancer. Does not have to be unique.
	Name *string `json:"name,omitempty"`

	// Human-readable description for the Loadbalancer.
	Description *string `json:"description,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// The ID of the QoS Policy which will apply to the Virtual IP
	VipQosPolicyID *string `json:"vip_qos_policy_id,omitempty"`

	// Tags is a set of resource tags.
	Tags *[]string `json:"tags,omitempty"`
}

type SessionPersistence struct {
	// The type of persistence mode.
	Type string `json:"type"`

	// Name of cookie if persistence mode is set appropriately.
	CookieName string `json:"cookie_name,omitempty"`
}

type CreateMemberOpts struct {
	// The IP address of the member to receive traffic from the load balancer.
	Address string `json:"address" required:"true"`

	// The port on which to listen for client traffic.
	ProtocolPort int `json:"protocol_port" required:"true"`

	// Name of the Member.
	Name string `json:"name,omitempty"`

	// ProjectID is the UUID of the project who owns the Member.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// A positive integer value that indicates the relative portion of traffic
	// that this member should receive from the pool. For example, a member with
	// a weight of 10 receives five times as much traffic as a member with a
	// weight of 2.
	Weight int `json:"weight,omitempty"`

	// If you omit this parameter, LBaaS uses the vip_subnet_id parameter value
	// for the subnet UUID.
	SubnetID string `json:"subnet_id,omitempty"`

	// The administrative state of the Pool. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// Is the member a backup? Backup members only receive traffic when all
	// non-backup members are down.
	// Requires microversion 2.1 or later.
	Backup *bool `json:"backup,omitempty"`

	// An alternate IP address used for health monitoring a backend member.
	MonitorAddress *string `json:"monitor_address,omitempty"`

	// An alternate protocol port used for health monitoring a backend member.
	MonitorPort *int `json:"monitor_port,omitempty"`

	// A list of simple strings assigned to the resource.
	// Requires microversion 2.5 or later.
	Tags []string `json:"tags,omitempty"`
}

type CreatePoolOpts struct {
	// The algorithm used to distribute load between the members of the pool. The
	// current specification supports LBMethodRoundRobin, LBMethodLeastConnections,
	// LBMethodSourceIp and LBMethodSourceIpPort as valid values for this attribute.
	LBMethod LBMethod `json:"lb_algorithm" required:"true"`

	// The protocol used by the pool members, you can use either
	// ProtocolTCP, ProtocolUDP, ProtocolPROXY, ProtocolHTTP, ProtocolHTTPS,
	// ProtocolSCTP or ProtocolPROXYV2.
	Protocol Protocol `json:"protocol" required:"true"`

	// The Loadbalancer on which the members of the pool will be associated with.
	// Note: one of LoadbalancerID or ListenerID must be provided.
	LoadbalancerID string `json:"loadbalancer_id,omitempty"`

	// The Listener on which the members of the pool will be associated with.
	// Note: one of LoadbalancerID or ListenerID must be provided.
	ListenerID string `json:"listener_id,omitempty"`

	// ProjectID is the UUID of the project who owns the Pool.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// Name of the pool.
	Name string `json:"name,omitempty"`

	// Human-readable description for the pool.
	Description string `json:"description,omitempty"`

	// Persistence is the session persistence of the pool.
	// Omit this field to prevent session persistence.
	Persistence *SessionPersistence `json:"session_persistence,omitempty"`

	// The administrative state of the Pool. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// Members is a slice of BatchUpdateMemberOpts which allows a set of
	// members to be created at the same time the pool is created.
	//
	// This is only possible to use when creating a fully populated
	// Loadbalancer.
	Members []CreateMemberOpts `json:"members,omitempty"`

	// Monitor is an instance of monitors.CreateOpts which allows a monitor
	// to be created at the same time the pool is created.
	//
	// This is only possible to use when creating a fully populated
	// Loadbalancer.
	Monitor *CreateHealthMonitorOpts `json:"healthmonitor,omitempty"`

	// Tags is a set of resource tags. New in version 2.5
	Tags []string `json:"tags,omitempty"`
}

// UpdatePoolOpts is the common options struct used in this package's Update
// operation.
type UpdatePoolOpts struct {
	// Name of the pool.
	Name *string `json:"name,omitempty"`

	// Human-readable description for the pool.
	Description *string `json:"description,omitempty"`

	// The algorithm used to distribute load between the members of the pool. The
	// current specification supports LBMethodRoundRobin, LBMethodLeastConnections,
	// LBMethodSourceIp and LBMethodSourceIpPort as valid values for this attribute.
	LBMethod LBMethod `json:"lb_algorithm,omitempty"`

	// The administrative state of the Pool. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// Persistence is the session persistence of the pool.
	Persistence *SessionPersistence `json:"session_persistence,omitempty"`

	// Tags is a set of resource tags. New in version 2.5
	Tags *[]string `json:"tags,omitempty"`
}

type CreateListenerOpts struct {
	// The load balancer on which to provision this listener.
	LoadbalancerID string `json:"loadbalancer_id,omitempty"`

	// The protocol - can either be TCP, SCTP, HTTP, HTTPS or TERMINATED_HTTPS.
	Protocol Protocol `json:"protocol" required:"true"`

	// The port on which to listen for client traffic.
	ProtocolPort int `json:"protocol_port" required:"true"`

	// ProjectID is only required if the caller has an admin role and wants
	// to create a pool for another project.
	ProjectID string `json:"project_id,omitempty"`

	// Human-readable name for the Listener. Does not have to be unique.
	Name string `json:"name,omitempty"`

	// The ID of the default pool with which the Listener is associated.
	DefaultPoolID string `json:"default_pool_id,omitempty"`

	// DefaultPool an instance of pools.CreateOpts which allows a
	// (default) pool to be created at the same time the listener is created.
	//
	// This is only possible to use when creating a fully populated
	// load balancer.
	DefaultPool *CreatePoolOpts `json:"default_pool,omitempty"`

	// Human-readable description for the Listener.
	Description string `json:"description,omitempty"`

	// The maximum number of connections allowed for the Listener.
	ConnLimit *int `json:"connection_limit,omitempty"`

	// A reference to a Barbican container of TLS secrets.
	DefaultTlsContainerRef string `json:"default_tls_container_ref,omitempty"`

	// A list of references to TLS secrets.
	SniContainerRefs []string `json:"sni_container_refs,omitempty"`

	// The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// L7Policies is a slice of l7policies.CreateOpts which allows a set
	// of policies to be created at the same time the listener is created.
	//
	// This is only possible to use when creating a fully populated
	// Loadbalancer.
	L7Policies []CreateL7PoliciesOpts `json:"l7policies,omitempty"`

	// Frontend client inactivity timeout in milliseconds
	TimeoutClientData *int `json:"timeout_client_data,omitempty"`

	// Backend member inactivity timeout in milliseconds
	TimeoutMemberData *int `json:"timeout_member_data,omitempty"`

	// Backend member connection timeout in milliseconds
	TimeoutMemberConnect *int `json:"timeout_member_connect,omitempty"`

	// Time, in milliseconds, to wait for additional TCP packets for content inspection
	TimeoutTCPInspect *int `json:"timeout_tcp_inspect,omitempty"`

	// A dictionary of optional headers to insert into the request before it is sent to the backend member.
	InsertHeaders map[string]string `json:"insert_headers,omitempty"`

	// A list of IPv4, IPv6 or mix of both CIDRs
	AllowedCIDRs []string `json:"allowed_cidrs,omitempty"`

	// A list of TLS protocol versions. Available from microversion 2.17
	TLSVersions []TLSVersion `json:"tls_versions,omitempty"`

	// Tags is a set of resource tags. New in version 2.5
	Tags []string `json:"tags,omitempty"`
}

// UpdateListenerOpts represents options for updating a Listener.
type UpdateListenerOpts struct {
	// Human-readable name for the Listener. Does not have to be unique.
	Name *string `json:"name,omitempty"`

	// The ID of the default pool with which the Listener is associated.
	DefaultPoolID *string `json:"default_pool_id,omitempty"`

	// Human-readable description for the Listener.
	Description *string `json:"description,omitempty"`

	// The maximum number of connections allowed for the Listener.
	ConnLimit *int `json:"connection_limit,omitempty"`

	// A reference to a Barbican container of TLS secrets.
	DefaultTlsContainerRef *string `json:"default_tls_container_ref,omitempty"`

	// A list of references to TLS secrets.
	SniContainerRefs *[]string `json:"sni_container_refs,omitempty"`

	// The administrative state of the Listener. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// Frontend client inactivity timeout in milliseconds
	TimeoutClientData *int `json:"timeout_client_data,omitempty"`

	// Backend member inactivity timeout in milliseconds
	TimeoutMemberData *int `json:"timeout_member_data,omitempty"`

	// Backend member connection timeout in milliseconds
	TimeoutMemberConnect *int `json:"timeout_member_connect,omitempty"`

	// Time, in milliseconds, to wait for additional TCP packets for content inspection
	TimeoutTCPInspect *int `json:"timeout_tcp_inspect,omitempty"`

	// A dictionary of optional headers to insert into the request before it is sent to the backend member.
	InsertHeaders *map[string]string `json:"insert_headers,omitempty"`

	// A list of IPv4, IPv6 or mix of both CIDRs
	AllowedCIDRs *[]string `json:"allowed_cidrs,omitempty"`

	// A list of TLS protocol versions. Available from microversion 2.17
	TLSVersions *[]TLSVersion `json:"tls_versions,omitempty"`

	// Tags is a set of resource tags. New in version 2.5
	Tags *[]string `json:"tags,omitempty"`
}


// CreateHealthMonitorOpts is the common options struct used in this package's Create
// operation.
type CreateHealthMonitorOpts struct {
	// The Pool to Monitor.
	PoolID string `json:"pool_id,omitempty"`

	// The type of probe, which is PING, TCP, HTTP, or HTTPS, that is
	// sent by the load balancer to verify the member state.
	Type string `json:"type" required:"true"`

	// The time, in seconds, between sending probes to members.
	Delay int `json:"delay" required:"true"`

	// Maximum number of seconds for a Monitor to wait for a ping reply
	// before it times out. The value must be less than the delay value.
	Timeout int `json:"timeout" required:"true"`

	// Number of permissible ping failures before changing the member's
	// status to INACTIVE. Must be a number between 1 and 10.
	MaxRetries int `json:"max_retries" required:"true"`

	// Number of permissible ping failures befor changing the member's
	// status to ERROR. Must be a number between 1 and 10.
	MaxRetriesDown int `json:"max_retries_down,omitempty"`

	// URI path that will be accessed if Monitor type is HTTP or HTTPS.
	URLPath string `json:"url_path,omitempty"`

	// The HTTP method used for requests by the Monitor. If this attribute
	// is not specified, it defaults to "GET". Required for HTTP(S) types.
	HTTPMethod string `json:"http_method,omitempty"`

	// Expected HTTP codes for a passing HTTP(S) Monitor. You can either specify
	// a single status like "200", a range like "200-202", or a combination like
	// "200-202, 401".
	ExpectedCodes string `json:"expected_codes,omitempty"`

	// TenantID is the UUID of the project who owns the Monitor.
	// Only administrative users can specify a project UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`

	// ProjectID is the UUID of the project who owns the Monitor.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// The Name of the Monitor.
	Name string `json:"name,omitempty"`

	// The administrative state of the Monitor. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

// UpdateHealthMonitorOpts is the common options struct used in this package's Update
// operation.
type UpdateHealthMonitorOpts struct {
	// The time, in seconds, between sending probes to members.
	Delay int `json:"delay,omitempty"`

	// Maximum number of seconds for a Monitor to wait for a ping reply
	// before it times out. The value must be less than the delay value.
	Timeout int `json:"timeout,omitempty"`

	// Number of permissible ping failures before changing the member's
	// status to INACTIVE. Must be a number between 1 and 10.
	MaxRetries int `json:"max_retries,omitempty"`

	// Number of permissible ping failures befor changing the member's
	// status to ERROR. Must be a number between 1 and 10.
	MaxRetriesDown int `json:"max_retries_down,omitempty"`

	// URI path that will be accessed if Monitor type is HTTP or HTTPS.
	// Required for HTTP(S) types.
	URLPath string `json:"url_path,omitempty"`

	// The HTTP method used for requests by the Monitor. If this attribute
	// is not specified, it defaults to "GET". Required for HTTP(S) types.
	HTTPMethod string `json:"http_method,omitempty"`

	// Expected HTTP codes for a passing HTTP(S) Monitor. You can either specify
	// a single status like "200", or a range like "200-202". Required for HTTP(S)
	// types.
	ExpectedCodes string `json:"expected_codes,omitempty"`

	// The Name of the Monitor.
	Name *string `json:"name,omitempty"`

	// The administrative state of the Monitor. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}


// CreateL7PoliciesOpts is the common options struct used in this package's Create
// operation.
type CreateL7PoliciesOpts struct {
	// Name of the L7 policy.
	Name string `json:"name,omitempty"`

	// The ID of the listener.
	ListenerID string `json:"listener_id,omitempty"`

	// The L7 policy action. One of REDIRECT_PREFIX, REDIRECT_TO_POOL, REDIRECT_TO_URL, or REJECT.
	Action Action `json:"action" required:"true"`

	// The position of this policy on the listener.
	Position int32 `json:"position,omitempty"`

	// A human-readable description for the resource.
	Description string `json:"description,omitempty"`

	// ProjectID is the UUID of the project who owns the L7 policy in octavia.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// Requests matching this policy will be redirected to this Prefix URL.
	// Only valid if action is REDIRECT_PREFIX.
	RedirectPrefix string `json:"redirect_prefix,omitempty"`

	// Requests matching this policy will be redirected to the pool with this ID.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolID string `json:"redirect_pool_id,omitempty"`

	// Requests matching this policy will be redirected to this URL.
	// Only valid if action is REDIRECT_TO_URL.
	RedirectURL string `json:"redirect_url,omitempty"`

	// Requests matching this policy will be redirected to the specified URL or Prefix URL
	// with the HTTP response code. Valid if action is REDIRECT_TO_URL or REDIRECT_PREFIX.
	// Valid options are: 301, 302, 303, 307, or 308. Default is 302. Requires version 2.9
	RedirectHttpCode int32 `json:"redirect_http_code,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`

	// Rules is a slice of CreateRuleOpts which allows a set of rules
	// to be created at the same time the policy is created.
	//
	// This is only possible to use when creating a fully populated
	// Loadbalancer.
	Rules []CreateRuleOpts `json:"rules,omitempty"`
}

// ListL7PoliciesOpts allows the filtering and sorting of paginated collections through
// the API.
type ListL7PoliciesOpts struct {
	Name           string `q:"name"`
	Description    string `q:"description"`
	ListenerID     string `q:"listener_id"`
	Action         string `q:"action"`
	ProjectID      string `q:"project_id"`
	RedirectPoolID string `q:"redirect_pool_id"`
	RedirectURL    string `q:"redirect_url"`
	Position       int32  `q:"position"`
	AdminStateUp   bool   `q:"admin_state_up"`
	ID             string `q:"id"`
	Limit          int    `q:"limit"`
	Marker         string `q:"marker"`
	SortKey        string `q:"sort_key"`
	SortDir        string `q:"sort_dir"`
}

// UpdateL7PoliciesOpts is the common options struct used in this package's Update
// operation.
type UpdateL7PoliciesOpts struct {
	// Name of the L7 policy, empty string is allowed.
	Name *string `json:"name,omitempty"`

	// The L7 policy action. One of REDIRECT_PREFIX, REDIRECT_TO_POOL, REDIRECT_TO_URL, or REJECT.
	Action Action `json:"action,omitempty"`

	// The position of this policy on the listener.
	Position int32 `json:"position,omitempty"`

	// A human-readable description for the resource, empty string is allowed.
	Description *string `json:"description,omitempty"`

	// Requests matching this policy will be redirected to this Prefix URL.
	// Only valid if action is REDIRECT_PREFIX.
	RedirectPrefix *string `json:"redirect_prefix,omitempty"`

	// Requests matching this policy will be redirected to the pool with this ID.
	// Only valid if action is REDIRECT_TO_POOL.
	RedirectPoolID *string `json:"redirect_pool_id,omitempty"`

	// Requests matching this policy will be redirected to this URL.
	// Only valid if action is REDIRECT_TO_URL.
	RedirectURL *string `json:"redirect_url,omitempty"`

	// Requests matching this policy will be redirected to the specified URL or Prefix URL
	// with the HTTP response code. Valid if action is REDIRECT_TO_URL or REDIRECT_PREFIX.
	// Valid options are: 301, 302, 303, 307, or 308. Default is 302. Requires version 2.9
	RedirectHttpCode int32 `json:"redirect_http_code,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

type CreateRuleOpts struct {
	// The L7 rule type. One of COOKIE, FILE_TYPE, HEADER, HOST_NAME, or PATH.
	RuleType RuleType `json:"type" required:"true"`

	// The comparison type for the L7 rule. One of CONTAINS, ENDS_WITH, EQUAL_TO, REGEX, or STARTS_WITH.
	CompareType CompareType `json:"compare_type" required:"true"`

	// The value to use for the comparison. For example, the file type to compare.
	Value string `json:"value" required:"true"`

	// ProjectID is the UUID of the project who owns the rule in octavia.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// The key to use for the comparison. For example, the name of the cookie to evaluate.
	Key string `json:"key,omitempty"`

	// When true the logic of the rule is inverted. For example, with invert true,
	// equal to would become not equal to. Default is false.
	Invert bool `json:"invert,omitempty"`

	// The administrative state of the Loadbalancer. A valid value is true (UP)
	// or false (DOWN).
	AdminStateUp *bool `json:"admin_state_up,omitempty"`
}

func (opts *CreateLoadbalancerOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.LOADBALANCER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateListenerOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.LISTENER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreatePoolOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.POOL)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateMemberOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.MEMBER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateHealthMonitorOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.HEALTHMONITOR)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateL7PoliciesOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.L7POLICY)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateRuleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.L7RULE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type Loadbalancer struct {
	Description        string `json:"description"`
	AdminStateUp       bool   `json:"admin_state_up"`
	ProjectId          string `json:"project_id"`
	ProvisioningStatus string `json:"provisioning_status"`
	FlavorId           string `json:"flavor_id"`
	VipSubnetId        string `json:"vip_subnet_id"`
	VipAddress         string `json:"vip_address"`
	VipNetworkId       string `json:"vip_network_id"`
	VipPortId          string `json:"vip_port_id"`
	AdditionalVips     []struct {
		SubnetId  string `json:"subnet_id"`
		IpAddress string `json:"ip_address"`
	} `json:"additional_vips"`
	Provider         string   `json:"provider"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	Id               string   `json:"id"`
	OperatingStatus  string   `json:"operating_status"`
	Name             string   `json:"name"`
	VipQosPolicyId   string   `json:"vip_qos_policy_id"`
	AvailabilityZone string   `json:"availability_zone"`
	Tags             []string `json:"tags"`
}

type LoadbalancerMap struct {
	Loadbalancer       `json:"loadbalancer"`
}

type Loadbalancers struct {
	LBs               []Loadbalancer          `json:"loadbalancers"`
}

type Listener struct {
	Description            string `json:"description"`
	AdminStateUp           bool   `json:"admin_state_up"`
	ProjectId              string `json:"project_id"`
	Protocol               string `json:"protocol"`
	ProtocolPort           int    `json:"protocol_port"`
	ProvisioningStatus     string `json:"provisioning_status"`
	DefaultTlsContainerRef string `json:"default_tls_container_ref"`
	Loadbalancers          []struct {
		Id string `json:"id"`
	} `json:"loadbalancers"`
	InsertHeaders struct {
		XForwardedPort string `json:"X-Forwarded-Port"`
		XForwardedFor  string `json:"X-Forwarded-For"`
	} `json:"insert_headers"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	Id               string   `json:"id"`
	OperatingStatus  string   `json:"operating_status"`
	DefaultPoolId    string   `json:"default_pool_id"`
	SniContainerRefs []string `json:"sni_container_refs"`
	L7Policies       []struct {
		Id string `json:"id"`
	} `json:"l7policies"`
	Name                    string   `json:"name"`
	TimeoutClientData       int      `json:"timeout_client_data"`
	TimeoutMemberConnect    int      `json:"timeout_member_connect"`
	TimeoutMemberData       int      `json:"timeout_member_data"`
	TimeoutTcpInspect       int      `json:"timeout_tcp_inspect"`
	Tags                    []string `json:"tags"`
	ClientCaTlsContainerRef string   `json:"client_ca_tls_container_ref"`
	ClientAuthentication    string   `json:"client_authentication"`
	ClientCrlContainerRef   string   `json:"client_crl_container_ref"`
	AllowedCidrs            []string `json:"allowed_cidrs"`
	TlsCiphers              string   `json:"tls_ciphers"`
	TlsVersions             []string `json:"tls_versions"`
	AlpnProtocols           []string `json:"alpn_protocols"`
}

type ListenerMap struct {
    Listener          	 `json:"listener"`
}

type Listeners struct {
	Liss            []Listener      `json:"listeners"`
}

type Pool struct {
	LbAlgorithm   string `json:"lb_algorithm"`
	Protocol      string `json:"protocol"`
	Description   string `json:"description"`
	AdminStateUp  bool   `json:"admin_state_up"`
	Loadbalancers []struct {
		Id string `json:"id"`
	} `json:"loadbalancers"`
	CreatedAt          string `json:"created_at"`
	ProvisioningStatus string `json:"provisioning_status"`
	UpdatedAt          string `json:"updated_at"`
	SessionPersistence struct {
		CookieName string `json:"cookie_name"`
		Type       string `json:"type"`
	} `json:"session_persistence"`
	Listeners []struct {
		Id string `json:"id"`
	} `json:"listeners"`
	Members           []Member      `json:"members"`
	HealthmonitorId   interface{}   `json:"healthmonitor_id"`
	ProjectId         string        `json:"project_id"`
	Id                string        `json:"id"`
	OperatingStatus   string        `json:"operating_status"`
	Name              string        `json:"name"`
	Tags              []string      `json:"tags"`
	TlsContainerRef   string        `json:"tls_container_ref"`
	CaTlsContainerRef string        `json:"ca_tls_container_ref"`
	CrlContainerRef   string        `json:"crl_container_ref"`
	TlsEnabled        bool          `json:"tls_enabled"`
	TlsCiphers        string        `json:"tls_ciphers"`
	TlsVersions       []string      `json:"tls_versions"`
	AlpnProtocols     []string      `json:"alpn_protocols"`
}

type PoolMap struct {
    Pool     	 `json:"pool"`
}

type Pools struct {
	Ps          []Pool        `json:"pools"`
}

type Member struct {
	MonitorPort        int         `json:"monitor_port"`
	ProjectId          string      `json:"project_id"`
	Name               string      `json:"name"`
	Weight             int         `json:"weight"`
	Backup             bool        `json:"backup"`
	AdminStateUp       bool        `json:"admin_state_up"`
	SubnetId           string      `json:"subnet_id"`
	CreatedAt          string      `json:"created_at"`
	ProvisioningStatus string      `json:"provisioning_status"`
	MonitorAddress     interface{} `json:"monitor_address"`
	UpdatedAt          string      `json:"updated_at"`
	Address            string      `json:"address"`
	ProtocolPort       int         `json:"protocol_port"`
	Id                 string      `json:"id"`
	OperatingStatus    string      `json:"operating_status"`
	Tags               []string    `json:"tags"`
}

type MemberMap struct {
    Member     	 `json:"member"`
}

type Members struct {
	Ms                []Member         `json:"members"`
}

type Healthmonitor struct {
	ProjectId    string `json:"project_id"`
	Name         string `json:"name"`
	AdminStateUp bool   `json:"admin_state_up"`
	Pools        []struct {
		Id string `json:"id"`
	} `json:"pools"`
	CreatedAt          string   `json:"created_at"`
	ProvisioningStatus string   `json:"provisioning_status"`
	UpdatedAt          string   `json:"updated_at"`
	Delay              int      `json:"delay"`
	ExpectedCodes      string   `json:"expected_codes"`
	MaxRetries         int      `json:"max_retries"`
	HttpMethod         string   `json:"http_method"`
	Timeout            int      `json:"timeout"`
	MaxRetriesDown     int      `json:"max_retries_down"`
	UrlPath            string   `json:"url_path"`
	Type               string   `json:"type"`
	Id                 string   `json:"id"`
	OperatingStatus    string   `json:"operating_status"`
	Tags               []string `json:"tags"`
	HttpVersion        float64  `json:"http_version"`
	DomainName         string   `json:"domain_name"`
}

type HealthMonitorMap struct {
    Healthmonitor       	 `json:"healthmonitor"`
}

type HealthMonitors struct {
	HMs           []Healthmonitor            `json:"healthmonitors"`
}

type L7Rule struct {
    Id     string `json:"id"`
}

type L7Policy struct {
	ListenerId   string `json:"listener_id"`
	Description  string `json:"description"`
	AdminStateUp bool   `json:"admin_state_up"`
	Rules       []L7Rule           `json:"rules"`
	CreatedAt          string      `json:"created_at"`
	ProvisioningStatus string      `json:"provisioning_status"`
	UpdatedAt          string      `json:"updated_at"`
	RedirectHttpCode   int         `json:"redirect_http_code"`
	RedirectPoolId     interface{} `json:"redirect_pool_id"`
	RedirectPrefix     interface{} `json:"redirect_prefix"`
	RedirectUrl        string      `json:"redirect_url"`
	Action             string      `json:"action"`
	Position           int         `json:"position"`
	ProjectId          string      `json:"project_id"`
	Id                 string      `json:"id"`
	OperatingStatus    string      `json:"operating_status"`
	Name               string      `json:"name"`
	Tags               []string    `json:"tags"`
}

type L7PolicyMap struct {
    L7Policy      	 `json:"l7policy"`
}

type L7Policies struct {
	L7Ps           []L7Policy   `json:"l7policies"`
}

type L7RuleMap struct {
	Rule struct {
		CreatedAt          string      `json:"created_at"`
		CompareType        string      `json:"compare_type"`
		ProvisioningStatus string      `json:"provisioning_status"`
		Invert             bool        `json:"invert"`
		AdminStateUp       bool        `json:"admin_state_up"`
		UpdatedAt          string      `json:"updated_at"`
		Value              string      `json:"value"`
		Key                interface{} `json:"key"`
		ProjectId          string      `json:"project_id"`
		Type               string      `json:"type"`
		Id                 string      `json:"id"`
		OperatingStatus    string      `json:"operating_status"`
		Tags               []string    `json:"tags"`
	} `json:"rule"`
}
