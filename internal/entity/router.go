package entity

import (
	"fmt"
	"reflect"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"strings"
	"time"
)

type Router struct {
	AdminStateUp          bool          `json:"admin_state_up"`
	AvailabilityZoneHints []interface{} `json:"availability_zone_hints"`
	AvailabilityZones     []string      `json:"availability_zones"`
	CreatedAt             time.Time     `json:"created_at"`
	Description           string        `json:"description"`
	Distributed           bool          `json:"distributed"`
	GatewayInfo                         `json:"external_gateway_info"`
	FlavorId         string        `json:"flavor_id"`
	Ha               bool          `json:"ha"`
	Id               string        `json:"id"`
	Name             string             `json:"name"`
	Routes           []interface{}      `json:"routes"`
	RevisionNumber   int                `json:"revision_number"`
	Status           string             `json:"status"`
	UpdatedAt        time.Time          `json:"updated_at"`
	ProjectId        string             `json:"project_id"`
	TenantId         string             `json:"tenant_id"`
	ServiceTypeId    interface{}        `json:"service_type_id"`
	Tags             []string           `json:"tags"`
	ConntrackHelpers []interface{}      `json:"conntrack_helpers"`
	RouterInterfaces  []RouterInterface `json:"router_interfaces"`
}

type RouterMap struct {
	Router `json:"router"`
}

type Routers struct {
	Rs             []Router `json:"routers"`
	Count              int  `json:"count"`
}

type RouterInterface struct {
	Id        string   `json:"id"`
	NetworkId string   `json:"network_id"`
	PortId    string   `json:"port_id"`
	SubnetId  string   `json:"subnet_id"`
	SubnetIds []string `json:"subnet_ids"`
	ProjectId string   `json:"project_id"`
	TenantId  string   `json:"tenant_id"`
	Tags      []string `json:"tags"`
}

type ExternalFixedIP struct {
	IPAddress string `json:"ip_address,omitempty"`
	SubnetID  string `json:"subnet_id"`
}

type GatewayInfo struct {
	NetworkID        string            `json:"network_id,omitempty"`
	EnableSNAT       bool              `json:"enable_snat,omitempty"`
	ExternalFixedIPs []ExternalFixedIP `json:"external_fixed_ips,omitempty"`
	QosPolicyId      string            `json:"qos_policy_id,omitempty"`
}

type CreateRouterOpts struct {
	Name                  string       `json:"name,omitempty"`
	Description           string       `json:"description,omitempty"`
	AdminStateUp          *bool        `json:"admin_state_up,omitempty"`
	Distributed           *bool        `json:"distributed,omitempty"`
	TenantID              string       `json:"tenant_id,omitempty"`
	ProjectID             string       `json:"project_id,omitempty"`
	GatewayInfo           *GatewayInfo `json:"external_gateway_info,omitempty"`
	AvailabilityZoneHints []string     `json:"availability_zone_hints,omitempty"`
}

type Route struct {
	NextHop         string `json:"nexthop"`
	DestinationCIDR string `json:"destination"`
}

func (opts *CreateRouterOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.ROUTER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *CreateRouterOpts) AssignProps(props map[string]interface{}) (*CreateRouterOpts, map[string]string) {
	var deps = make(map[string]string)
	typ := reflect.TypeOf(opts)
	val := reflect.ValueOf(opts).Elem()
	for i := 0;i < typ.Elem().NumField();i++ {
		field := typ.Elem().Field(i)
		tag := field.Tag.Get("json")
		compareFieldName := strings.Split(tag, ",")[0]
		value := val.Field(i)
		fieldType := field.Type
		if v, ok := props[compareFieldName]; ok {
			switch fieldType.Kind() {
			case reflect.String:
				value.SetString(v.(string))
			case reflect.Bool:
				value.SetBool(v.(bool))
			case reflect.Ptr:
				if compareFieldName == "external_gateway_info" {
					obj := GatewayInfo{}
					if networkId, ok := v.(map[string]interface{})["network_id"]; ok {
						if networkId.(string) == "local" {
							obj.NetworkID = configs.CONF.ExternalNetwork
						} else {
							obj.NetworkID = networkId.(string)
						}
					}
					if enableSnat, ok := v.(map[string]interface{})["enable_snat"]; ok {
						obj.EnableSNAT = enableSnat.(bool)
					}
					if fixedIPs, ok := v.(map[string]interface{})["external_fixed_ips"]; ok {
						fixedIPObjs := make([]ExternalFixedIP, len(fixedIPs.([]interface{})))
						for _, ip := range fixedIPs.([]interface{}) {
							fixedIpObj := ExternalFixedIP{}
							if subnetId, ok := ip.(map[string]interface{})["subnet_id"]; ok {
								fixedIpObj.SubnetID = subnetId.(string)
							}
							if ipAddress, ok := ip.(map[string]interface{})["ip_address"]; ok {
								fixedIpObj.IPAddress = ipAddress.(string)
							}
							fixedIPObjs = append(fixedIPObjs, fixedIpObj)
						}
						obj.ExternalFixedIPs = fixedIPObjs
					}
					v = obj
				}
				ptr := reflect.New(fieldType.Elem())
				ptr.Elem().Set(reflect.ValueOf(v))
				value.Set(ptr)
			case reflect.Slice:
				value.Set(reflect.ValueOf(val))
			}
		}
	}
	return opts, deps
}

type UpdateRouterOpts struct {
	Name         string       `json:"name,omitempty"`
	Description  *string      `json:"description,omitempty"`
	AdminStateUp *bool        `json:"admin_state_up,omitempty"`
	Distributed  *bool        `json:"distributed,omitempty"`
	GatewayInfo  *GatewayInfo `json:"external_gateway_info,omitempty"`
	Routes       *[]Route     `json:"routes,omitempty"`
}


func (opts *UpdateRouterOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.ROUTER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type AddRouterInterfaceOpts struct {
	SubnetID              string       `json:"subnet_id,omitempty"`
	RouterId              string       `json:"router_id,-"`
}

func (opts *AddRouterInterfaceOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

func (opts *AddRouterInterfaceOpts) AssignProps(props map[string]interface{}) (*AddRouterInterfaceOpts, map[string]string) {
	var deps = make(map[string]string)
	typ := reflect.TypeOf(opts)
	val := reflect.ValueOf(opts).Elem()
	for i := 0;i < typ.Elem().NumField();i++ {
		field := typ.Elem().Field(i)
		tag := field.Tag.Get("json")
		compareFieldName := strings.Split(tag, ",")[0]
		value := val.Field(i)
		fieldType := field.Type
		if v, ok := props[compareFieldName]; ok {
			switch fieldType.Kind() {
			case reflect.String:
				if IsUUID(v.(string)) {
					value.SetString(v.(string))
				} else {
					deps[v.(string)] = field.Name
				}
			}
		}
	}
	return opts, deps
}

