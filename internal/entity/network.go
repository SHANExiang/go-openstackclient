package entity

import (
	"fmt"
	"reflect"
	"go-openstackclient/consts"
	"strings"
	"time"
)

type Network struct {
	AdminStateUp            bool          `json:"admin_state_up"`
	AvailabilityZoneHints   []interface{} `json:"availability_zone_hints"`
	AvailabilityZones       []string      `json:"availability_zones"`
	CreatedAt               time.Time     `json:"created_at"`
	Description             string        `json:"description"`
	Id                      string        `json:"id"`
	Ipv4AddressScope        interface{}   `json:"ipv4_address_scope"`
	Ipv6AddressScope        interface{}   `json:"ipv6_address_scope"`
	Mtu                     int           `json:"mtu"`
	Name                    string        `json:"name"`
	PortSecurityEnabled     bool          `json:"port_security_enabled"`
	ProjectId               string        `json:"project_id"`
	ProviderNetworkType     string        `json:"provider:network_type"`
	ProviderPhysicalNetwork interface{}   `json:"provider:physical_network"`
	ProviderSegmentationId  int           `json:"provider:segmentation_id"`
	QosPolicyId             interface{}   `json:"qos_policy_id"`
	RevisionNumber          int           `json:"revision_number"`
	RouterExternal          bool          `json:"router:external"`
	Shared                  bool          `json:"shared"`
	Status                  string        `json:"status"`
	Subnets                 []string      `json:"subnets"`
	Tags                    []interface{} `json:"tags"`
	TenantId                string        `json:"tenant_id"`
	UpdatedAt               time.Time     `json:"updated_at"`
}

type NetworkMap struct {
	Network `json:"network"`
}

type Networks struct {
	Nets               []Network `json:"networks"`
	Count              int       `json:"count"`
}

type CreateNetworkOpts struct {
	AdminStateUp          *bool    `json:"admin_state_up,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Description           string   `json:"description,omitempty"`
	Shared                *bool    `json:"shared,omitempty"`
	TenantID              string   `json:"tenant_id,omitempty"`
	ProjectID             string   `json:"project_id,omitempty"`
	AvailabilityZoneHints []string `json:"availability_zone_hints,omitempty"`
	RouterExternal        bool     `json:"router:external,omitempty"`
}

func (opts *CreateNetworkOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.NETWORK)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

// AssignProps second output parameter is dependent resources slice
func (opts *CreateNetworkOpts) AssignProps(props map[string]interface{}) (*CreateNetworkOpts, map[string]string) {
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
				value.SetString(v.(string))
			case reflect.Bool:
				value.SetBool(v.(bool))
			case reflect.Ptr:
				ptr := reflect.New(fieldType.Elem())
				ptr.Elem().Set(reflect.ValueOf(v))
				value.Set(ptr)
			case reflect.Slice:
				stringSlice := InterfaceSliceToStringSlice(v.([]interface{}))
				value.Set(reflect.ValueOf(stringSlice))
			}
		}
	}
	return opts, deps
}
