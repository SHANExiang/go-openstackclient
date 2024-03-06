package entity

import (
	"fmt"
	"math/rand"
	"reflect"
	"go-openstackclient/consts"
	"strings"
	"time"
)

type Subnet struct {
	Name              string        `json:"name"`
	EnableDhcp        bool          `json:"enable_dhcp"`
	NetworkId         string        `json:"network_id"`
	SegmentId         interface{}   `json:"segment_id"`
	ProjectId         string        `json:"project_id"`
	TenantId          string        `json:"tenant_id"`
	CreatedAt         string        `json:"created_at"`
	DnsNameservers    []interface{} `json:"dns_nameservers"`
	DnsPublishFixedIp bool          `json:"dns_publish_fixed_ip"`
	AllocationPools   []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"allocation_pools"`
	HostRoutes      []interface{} `json:"host_routes"`
	IpVersion       int           `json:"ip_version"`
	GatewayIp       string        `json:"gateway_ip"`
	Cidr            string        `json:"cidr"`
	UpdatedAt       string        `json:"updated_at"`
	Id              string        `json:"id"`
	Description     string        `json:"description"`
	Ipv6AddressMode interface{}   `json:"ipv6_address_mode"`
	Ipv6RaMode      interface{}   `json:"ipv6_ra_mode"`
	RevisionNumber  int           `json:"revision_number"`
	ServiceTypes    []interface{} `json:"service_types"`
	SubnetpoolId    interface{}   `json:"subnetpool_id"`
	Tags            []string      `json:"tags"`
}

type SubnetMap struct {
	 Subnet          `json:"subnet"`
}

type Subnets struct {
	Ss            []Subnet           `json:"subnets"`
}

type HostRoute struct {
	DestinationCIDR string `json:"destination"`
	NextHop         string `json:"nexthop"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type CreateSubnetOpts struct {
	// NetworkID is the UUID of the network the subnet will be associated with.
	NetworkID string `json:"network_id" required:"true"`

	// CIDR is the address CIDR of the subnet.
	CIDR string `json:"cidr,omitempty"`

	// Name is a human-readable name of the subnet.
	Name string `json:"name,omitempty"`

	// Description of the subnet.
	Description string `json:"description,omitempty"`

	// The UUID of the project who owns the Subnet. Only administrative users
	// can specify a project UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`

	// The UUID of the project who owns the Subnet. Only administrative users
	// can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// AllocationPools are IP Address pools that will be available for DHCP.
	AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`

	// GatewayIP sets gateway information for the subnet. Setting to nil will
	// cause a default gateway to automatically be created. Setting to an empty
	// string will cause the subnet to be created with no gateway. Setting to
	// an explicit address will set that address as the gateway.
	GatewayIP *string `json:"gateway_ip,omitempty"`

	// IPVersion is the IP version for the subnet.
	IPVersion int `json:"ip_version,omitempty"`

	// EnableDHCP will either enable to disable the DHCP service.
	EnableDHCP *bool `json:"enable_dhcp,omitempty"`

	// DNSNameservers are the nameservers to be set via DHCP.
	DNSNameservers []string `json:"dns_nameservers,omitempty"`

	// ServiceTypes are the service types associated with the subnet.
	ServiceTypes []string `json:"service_types,omitempty"`

	// HostRoutes are any static host routes to be set via DHCP.
	HostRoutes []HostRoute `json:"host_routes,omitempty"`

	// The IPv6 address modes specifies mechanisms for assigning IPv6 IP addresses.
	IPv6AddressMode string `json:"ipv6_address_mode,omitempty"`

	// The IPv6 router advertisement specifies whether the networking service
	// should transmit ICMPv6 packets.
	IPv6RAMode string `json:"ipv6_ra_mode,omitempty"`

	// SubnetPoolID is the id of the subnet pool that subnet should be associated to.
	SubnetPoolID string `json:"subnetpool_id,omitempty"`

	// Prefixlen is used when user creates a subnet from the subnetpool. It will
	// overwrite the "default_prefixlen" value of the referenced subnetpool.
	Prefixlen int `json:"prefixlen,omitempty"`
}

func (opts *CreateSubnetOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.SUBNET)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateSubnetOpts) RandomCidr() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(200)
	return fmt.Sprintf("192.%d.%d.0/24", randomNum, randomNum)
}

// AssignProps second output parameter is dependent resources slice
func (opts *CreateSubnetOpts) AssignProps(props map[string]interface{}) (*CreateSubnetOpts, map[string]string) {
	var deps = make(map[string]string)
	typ := reflect.TypeOf(opts)
	val := reflect.ValueOf(opts).Elem()
	for i := 0;i < typ.Elem().NumField();i++ {
		field := typ.Elem().Field(i)
		value := val.Field(i)
		fieldType := field.Type
		tag := field.Tag.Get("json")
		compareFieldName := strings.Split(tag, ",")[0]
		if v, ok := props[compareFieldName]; ok {
			switch fieldType.Kind() {
			case reflect.String:
				if compareFieldName == "cidr" && v.(string) == "random" {
					value.SetString(opts.RandomCidr())
				} else if compareFieldName == "network_id" && !IsUUID(v.(string)) {
					deps[v.(string)] = field.Name
				} else {
					value.SetString(v.(string))
				}
			case reflect.Bool:
				value.SetBool(v.(bool))
			case reflect.Int:
				value.SetInt(int64(v.(int)))
			case reflect.Ptr:
				ptr := reflect.New(fieldType.Elem())
				ptr.Elem().Set(reflect.ValueOf(v))
				value.Set(ptr)
			case reflect.Slice:
				if compareFieldName == "host_routes" {
					originSlice := make([]HostRoute, value.Len())
					for _, item := range v.([]interface{}) {
						destination := item.(map[string]interface{})["destination"].(string)
						nexthop := item.(map[string]interface{})["nexthop"].(string)
						obj := HostRoute{DestinationCIDR: destination, NextHop: nexthop}
						originSlice = append(originSlice, obj)
					}
					value.Set(reflect.ValueOf(originSlice))
				}
				if compareFieldName == "allocation_pools" {
					originSlice := make([]AllocationPool, value.Len())
					for _, item := range v.([]interface{}) {
						start := item.(map[string]interface{})["start"].(string)
						end := item.(map[string]interface{})["end"].(string)
						obj := AllocationPool{start, end}
						originSlice = append(originSlice, obj)
					}
					value.Set(reflect.ValueOf(originSlice))
				}
			}
		}
	}
	return opts, deps
}
