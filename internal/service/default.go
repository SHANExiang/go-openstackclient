package service

import (
	"fmt"
	"go-openstackclient/configs"
	"go-openstackclient/internal/entity"
	"math/rand"
	"time"
)

var (
	defaultName       = "sdn_test"
	defaultController = NewController(defaultName)
)

func CreateNetworkHelper() string {
	defaultOpts := defaultNetworkOpts()
    network := defaultController.CreateNetwork(defaultOpts)
	return network.Id
}

func ListNetworkHelper() entity.Networks {
	networks := defaultController.ListNetworks()
	return networks
}

func CreateSubnetHelper(networkId string) string {
	defaultOpts := defaultSubnetOpts(networkId)
	subnet := defaultController.CreateSubnet(defaultOpts)
	return subnet.Id
}

func CreateSecurityGroupHelper() string {
	defaultSg := defaultSgOpts()
    sg := defaultController.CreateSecurityGroup(defaultSg)
	return sg.Id
}

func CreateSecurityRuleICMP(sgId string) {
	ingressICMP := defaultICMPIngressSgRuleOpts(sgId)
	defaultController.CreateSecurityRule(ingressICMP)

	egressICMP := defaultICMPEgressSgRuleOpts(sgId)
	defaultController.CreateSecurityRule(egressICMP)
}


func CreateSecurityRuleSSH(sgId string) {
	ingressSSH := defaultSSHIngressSgRuleOpts(sgId)
	defaultController.CreateSecurityRule(ingressSSH)

	egressSSH := defaultSSHEgressSgRuleOpts(sgId)
	defaultController.CreateSecurityRule(egressSSH)
}

func CreatePortHelper(networkId, subnetId string) string {
	defaultOpts := defaultPortOpts(networkId, subnetId)
    port := defaultController.CreatePort(defaultOpts)
    return port.Id
}

func CreateRouterHelper() string {
	createOpts := defaultRouterOpts()
    router := defaultController.CreateRouter(createOpts)
	return router.Id
}

func SetRouterGatewayHelper(routerId, extNetId string) {
	createOpts := defaultRouterGatewayOpts(extNetId)
    defaultController.SetRouterGateway(createOpts, routerId)
}

func AddRouterInterfaceHelper(routerId, subnetId string) {
    opts := defaultRouterInterfaceOpts(routerId, subnetId)
	defaultController.AddRouterInterface(opts)
}

func CreateInstanceHelper(netId, sgName string) string {
	defaultController.EnsureSgExist(sgName)
	opts := defaultInstanceOpts(netId, sgName)
	server := defaultController.CreateInstance(opts)
	return server.Id
}

func CreateQosPolicyHelper() string {
	opts := defaultQosPolicyRequestOpts()
	qosPolicy := defaultController.createQosPolicy(opts)
	return qosPolicy.Id
}

func CreateBandwidthLimitRuleHelper(qosPolicyId string) {
	ingressOpts := defaultBandwidthLimitRuleIngressRequestOpts()
	egressOpts := defaultBandwidthLimitRuleEgressRequestOpts()
	defaultController.CreateBandwidthLimitRule(ingressOpts, qosPolicyId)
	defaultController.CreateBandwidthLimitRule(egressOpts, qosPolicyId)
}

func defaultNetworkOpts() entity.CreateUpdateOptions {
	return &entity.CreateNetworkOpts{Name: defaultName}
}

func defaultSubnetOpts(netId string) entity.CreateUpdateOptions {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(200)
	cidr := fmt.Sprintf("192.%d.%d.0/24", randomNum, randomNum)
	gatewayIp1 := fmt.Sprintf("192.%d.%d.1", randomNum, randomNum)
	subnetOpts := &entity.CreateSubnetOpts{
		NetworkID: netId,
		CIDR: cidr,
		IPVersion: 4,
		GatewayIP: &gatewayIp1,
		DNSNameservers: []string{"114.114.114.114"},
	}
	return subnetOpts
}

func defaultPortOpts(netId, subnetId string) entity.CreateUpdateOptions {
	fixedIp1 := entity.FixedIP{SubnetId: subnetId}
	fixedIp2 := entity.FixedIP{SubnetId: subnetId}
	opts := &entity.CreatePortOpts{
		FixedIp: []entity.FixedIP{fixedIp1, fixedIp2},
		NetworkId: netId,
	}
	return opts
}

func defaultSgOpts() entity.CreateUpdateOptions {
	opts := &entity.CreateSecurityGroupOpts{Name: defaultName}
	return opts
}

func defaultICMPIngressSgRuleOpts(sgId string) entity.CreateUpdateOptions {
	ingress := &entity.CreateSecurityRuleOpts{
		Direction: "ingress", EtherType: "IPv4",
		Protocol: "icmp", SecGroupID: sgId}
	return ingress
}

func defaultICMPEgressSgRuleOpts(sgId string) entity.CreateUpdateOptions {
	egress := &entity.CreateSecurityRuleOpts{
		Direction: "egress", EtherType: "IPv4",
		Protocol: "icmp", SecGroupID: sgId}
	return egress
}

func defaultSSHIngressSgRuleOpts(sgId string) entity.CreateUpdateOptions {
	return &entity.CreateSecurityRuleOpts{
		Direction: "ingress", EtherType: "IPv4", Protocol: "tcp",
		SecGroupID: sgId, PortRangeMin: 22, PortRangeMax: 22}
}

func defaultSSHEgressSgRuleOpts(sgId string) entity.CreateUpdateOptions {
	return &entity.CreateSecurityRuleOpts{
		Direction: "egress", EtherType: "IPv4", Protocol: "tcp",
		SecGroupID: sgId, PortRangeMin: 22, PortRangeMax: 22}
}

func defaultInstanceOpts(netId, sgName string) entity.CreateUpdateOptions {
	instanceOpts := &entity.CreateInstanceOpts{
		FlavorRef:      configs.CONF.FlavorId,
		ImageRef:       configs.CONF.ImageId,
		Networks:       []entity.ServerNet{{UUID: netId}},
		AdminPass:      "Wang.123",
		SecurityGroups: []entity.ServerSg{{Name: sgName}},
		Name:           defaultName,
		UserData: "Q29udGVudC1UeXBlOiBtdWx0aXBhcnQvbWl4ZWQ7IGJvdW5kYXJ5PSI9PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0iIApNSU1FLVZlcnNpb246IDEuMAoKLS09PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0KQ29udGVudC1UeXBlOiB0ZXh0L2Nsb3VkLWNvbmZpZzsgY2hhcnNldD0idXMtYXNjaWkiIApNSU1FLVZlcnNpb246IDEuMApDb250ZW50LVRyYW5zZmVyLUVuY29kaW5nOiA3Yml0CkNvbnRlbnQtRGlzcG9zaXRpb246IGF0dGFjaG1lbnQ7IGZpbGVuYW1lPSJzc2gtcHdhdXRoLXNjcmlwdC50eHQiIAoKI2Nsb3VkLWNvbmZpZwpkaXNhYmxlX3Jvb3Q6IGZhbHNlCnNzaF9wd2F1dGg6IHRydWUKcGFzc3dvcmQ6IFdhbmcuMTIzCgotLT09PT09PT09PT09PT09PTIzMDk5ODQwNTk3NDM3NjI0NzU9PQpDb250ZW50LVR5cGU6IHRleHQveC1zaGVsbHNjcmlwdDsgY2hhcnNldD0idXMtYXNjaWkiIApNSU1FLVZlcnNpb246IDEuMApDb250ZW50LVRyYW5zZmVyLUVuY29kaW5nOiA3Yml0CkNvbnRlbnQtRGlzcG9zaXRpb246IGF0dGFjaG1lbnQ7IGZpbGVuYW1lPSJwYXNzd2Qtc2NyaXB0LnR4dCIgCgojIS9iaW4vc2gKZWNobyAncm9vdDpXYW5nLjEyMycgfCBjaHBhc3N3ZAoKLS09PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0KQ29udGVudC1UeXBlOiB0ZXh0L3gtc2hlbGxzY3JpcHQ7IGNoYXJzZXQ9InVzLWFzY2lpIiAKTUlNRS1WZXJzaW9uOiAxLjAKQ29udGVudC1UcmFuc2Zlci1FbmNvZGluZzogN2JpdApDb250ZW50LURpc3Bvc2l0aW9uOiBhdHRhY2htZW50OyBmaWxlbmFtZT0iZW5hYmxlLWZzLWNvbGxlY3Rvci50eHQiIAoKIyEvYmluL3NoCnFlbXVfZmlsZT0iL2V0Yy9zeXNjb25maWcvcWVtdS1nYSIKaWYgWyAtZiAke3FlbXVfZmlsZX0gXTsgdGhlbgogICAgc2VkIC1pIC1yICJzL14jP0JMQUNLTElTVF9SUEM9LyNCTEFDS0xJU1RfUlBDPS8iICIke3FlbXVfZmlsZX0iCiAgICBoYXNfZ3FhPSQoc3lzdGVtY3RsIGxpc3QtdW5pdHMgLS1mdWxsIC1hbGwgLXQgc2VydmljZSAtLXBsYWluIHwgZ3JlcCAtbyBxZW11LWd1ZXN0LWFnZW50LnNlcnZpY2UpCiAgICBpZiBbWyAtbiAke2hhc19ncWF9IF1dOyB0aGVuCiAgICAgICAgc3lzdGVtY3RsIHJlc3RhcnQgcWVtdS1ndWVzdC1hZ2VudC5zZXJ2aWNlCiAgICBmaQpmaQoKLS09PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0tLQ==",
		BlockDeviceMappingV2: []entity.BlockDeviceMapping{{
			BootIndex: 0, Uuid: configs.CONF.ImageId, SourceType: "image",
			DestinationType: "volume", VolumeSize: 20, DeleteOnTermination: true,
		}},
	}
	return instanceOpts
}

func defaultRouterOpts() entity.CreateUpdateOptions {
	routerOpts := &entity.CreateRouterOpts{
		Name: defaultName, Description: defaultName,
	}
	return routerOpts
}

func defaultRouterGatewayOpts(externalNetId string) entity.CreateUpdateOptions {
	updateRouterOpts := &entity.UpdateRouterOpts{
		GatewayInfo: &entity.GatewayInfo{
			NetworkID: externalNetId}}
	return updateRouterOpts
}

func defaultRouterInterfaceOpts(routerId, subnetId string) entity.CreateUpdateOptions {
	addInterfaceOpts := &entity.AddRouterInterfaceOpts{
		SubnetID: subnetId, RouterId: routerId,
	}
	return addInterfaceOpts
}

func defaultQosPolicyRequestOpts() entity.CreateUpdateOptions {
	opts := &entity.CreateQosPolicyOpts{
		Name: defaultName,
	}
	return opts
}

func defaultBandwidthLimitRuleIngressRequestOpts() entity.CreateUpdateOptions {
	opts := &entity.CreateBandwidthLimitRuleOpts{
		MaxKBps: 10240, MaxBurstKBps: 0, Direction: "ingress",
	}
	return opts
}

func defaultBandwidthLimitRuleEgressRequestOpts() entity.CreateUpdateOptions {
	opts := &entity.CreateBandwidthLimitRuleOpts{
		MaxKBps: 20480, MaxBurstKBps: 0, Direction: "egress",
	}
	return opts
}
