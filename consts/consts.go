package consts

import (
	"time"
)


type Action string
type RuleType string
type CompareType string
type LBMethod string
type Protocol string
type TLSVersion string


const (
    POST                       = "POST"
    PUT                        = "PUT"
    GET                        = "GET"
    DELETE                     = "DELETE"
    PATCH                      = "PATCH"
    ContentTypeJson            = "application/json"
    ImagePatchJson             = "application/openstack-images-v2.1-json-patch"

    NOVA                       = "nova"
    CINDER                     = "cinder"

    VOLUME	                   = "volume"
    VOLUMES	                   = "volumes"
    SNAPSHOT                   = "snapshot"
    SNAPSHOTS                  = "snapshots"
    QOSSPEC                    = "qos_spec"
    QOSSPECS                   = "qos_specs"
    VOLUMETYPE                 = "type"
    VOLUMETYPES                = "types"
    ATTACHMENT                 = "attachment"
    NETWORKS                   = "networks"
    NETWORK                    = "network"
    SUBNETS                    = "subnets"
    SUBNET                     = "subnet"
    PORT                       = "port"
    PORTS                      = "ports"
    QOS_POLICIES               = "qos_policies"
    QOS_POLICY                 = "qos_policy"
    POLICY                     = "policy"
    BANDWIDTH_LIMIT_RULES      = "bandwidth_limit_rules"
    BANDWIDTH_LIMIT_RULE       = "bandwidth_limit_rule"
    DSCP_MARKING_RULES         = "dscp_marking_rules"
    DSCP_MARKING_RULE          = "dscp_marking_rule"
    MINIMUM_BANDWIDTH_RULES    = "minimum_bandwidth_rules"
    MINIMUM_BANDWIDTH_RULE     = "minimum_bandwidth_rule"
    ROUTERS                    = "routers"
    ROUTER                     = "router"
    ROUTERINTERFACE            = "router_interface"
    ROUTERGATEWAY              = "router_gateway"
    ROUTERROUTE                = "router_route"
    NETWORKROUTERINTERFACE     = "network:router_interface"
    FLOATINGIPS                = "floatingips"
    FLOATINGIP                 = "floatingip"
    NETWORKFLOATINGIP          = "network:floatingip"
    PORTFORWARDING             = "port_forwarding"
    PORTFORWARDINGS            = "port_forwardings"
    SERVER                     = "server"
    SERVERS                    = "servers"
    AGGREGATE                  = "aggregate"
    ADDHOST                    = "add_host"
    REMOVEHOST                 = "remove_host"
    SETMETADATA                = "set_metadata"
    REMOTECONSOLE              = "remote_console"
    FIREWALLGROUPS             = "firewall_groups"
    FIREWALLPOLICIES           = "firewall_policies"
    FIREWALLRULES              = "firewall_rules"
    FIREWALL                   = "firewall"
    FIREWALLPOLICY             = "firewall_policy"
    FIREWALLRULE               = "firewall_rule"
    PROJECT	                   = "project"
    PROJECTS	               = "projects"
    USER	                   = "user"
    USERS	                   = "users"
    TOKEN                      = "token"
    SECURITYGROUP              = "security_group"
    SECURITYGROUPS             = "security_groups"
    SECURITYGROUPRULES         = "security_group_rules"
    SECURITYGROUPRULE          = "security_group_rule"
    RBACPOLICIES               = "rbac_policies"
    VPNSERVICE                 = "vpnservice"
    VPNSERVICES                = "vpnservices"
    ENDPOINTGROUPS             = "endpoint_groups"
    IKEPOLICIES                = "ikepolicies"
    IPSECPOLICIES              = "ipsecpolicies"
    IPSECCONNECTIONS           = "ipsec_site_connections"
    LOADBALANCERS              = "loadbalancers"
    LOADBALANCER               = "loadbalancer"
    LISTENERS                  = "listeners"
    LISTENER                   = "listener"
    POOLS                      = "pools"
    POOL                       = "pool"
    MEMBERS                    = "members"
    MEMBER                     = "member"
    HEALTHMONITORS             = "healthmonitors"
    HEALTHMONITOR              = "healthmonitor"
    L7POLICIES                 = "l7policies"
    L7POLICY                   = "l7policy"
    L7RULES                    = "l7rules"
    L7RULE                     = "l7rule"
    VpcConnections             = "vpc_connections"
    VpcConnection              = "vpc_connection"
    Images                     = "images"
    Snats                      = "snats"
    Snat                       = "snat"
    Dnats                      = "dnats"
    Dnat                       = "dnat"

    ACTIVE                     = "ACTIVE"
    Available                  = "available"
    Error                      = "error"

    KeystonePort               = 5000
    CinderPort                 = 8776
    NeutronPort                = 9696
    NovaPort                   = 8774
    GlancePort                 = 9292
    OctaviaPort                = 9876
    SDNPort                    = 18002
    SDNMDCPort                 = 31943

    ADMIN                      = "admin"
	AuthToken                  = "X-Auth-Token"
	Timeout                    = 2 * 60 * time.Second
	IntervalTime               = 5 * time.Second


    ProtocolAny                = "any"
    ProtocolTCP                = "tcp"
    ProtocolICMP               = "icmp"
    ProtocolUDP                = "udp"
    ActionAllow                = "allow"
    ActionDeny                 = "deny"
    DirectionIngress           = "ingress"
    DirectionEgress            = "egress"
    EtherTypeV4                = "IPv4"
    EtherTypeV6                = "IPv6"
)
