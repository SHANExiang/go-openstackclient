package service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/entity"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	defaultController = Controller{}
	defaultName       = "sdn_test"
)

func CreateNetworkHelper() string {
	defaultOpts := defaultNetworkOpts()
    res := wrapper(constructNetworkRequestOpts)(defaultOpts, nil)
	defer fasthttp.ReleaseResponse(res)
	var network entity.NetworkMap
	_ = json.Unmarshal(res.Body(), &network)

	log.Println("==============Create internal network success", network.Network.Id)
	return network.Id
}

func CreateSubnetHelper(networkId string) string {
	createOpts := defaultSubnetOpts(networkId)
	res := wrapper(constructSubnetRequestOpts)(createOpts, nil)
	defer fasthttp.ReleaseResponse(res)

	var subnet entity.SubnetMap
	_ = json.Unmarshal(res.Body(), &subnet)

	log.Println("==============Create internal subnet success", subnet.Id)
	return subnet.Id
}

func CreateSecurityGroupHelper() string {
	defaultSg := defaultSgOpts()
	res := wrapper(constructSgRequestOpts)(defaultSg, nil)
	defer fasthttp.ReleaseResponse(res)

	var sg entity.Sg
	_ = json.Unmarshal(res.Body(), &sg)

	log.Println("==============Create security group success", sg.Id)
	return sg.Id
}

func CreateSecurityRuleHelper(opts entity.CreateUpdateOptions) {
	res := wrapper(constructSgRuleRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var sgRule entity.SgRule
	_ = json.Unmarshal(res.Body(), &sgRule)

	log.Println("==============Create security group rule success", sgRule.Id)
}

func CreateSecurityRuleICMP(sgId string) {
	ingressICMP := defaultICMPIngressSgRuleOpts(sgId)
	CreateSecurityRuleHelper(ingressICMP)

	egressICMP := defaultICMPEgressSgRuleOpts(sgId)
	CreateSecurityRuleHelper(egressICMP)
}


func CreateSecurityRuleSSH(sgId string) {
	ingressSSH := defaultSSHIngressSgRuleOpts(sgId)
	CreateSecurityRuleHelper(ingressSSH)

	egressSSH := defaultSSHEgressSgRuleOpts(sgId)
	CreateSecurityRuleHelper(egressSSH)
}

func CreatePortHelper(networkId, subnetId string) {
	createOpts := defaultPortOpts(networkId, subnetId)
    res := wrapper(constructPortRequestOpts)(createOpts, nil)
	defer fasthttp.ReleaseResponse(res)

	var port entity.PortMap
	_ = json.Unmarshal(res.Body(), &port)

	log.Println("==============Create internal port success", port.Id)
}

func CreateRouterHelper() string {
	createOpts := defaultRouterOpts()
    res := wrapper(constructRouterRequestOpts)(createOpts, nil)
	defer fasthttp.ReleaseResponse(res)

	var router entity.RouterMap
	_ = json.Unmarshal(res.Body(), &router)

	log.Println("==============Create router success", router.Id)
	return router.Id
}

func SetRouterGatewayHelper(routerId, extNetId string) {
	createOpts := defaultRouterGatewayOpts(extNetId)
    res := wrapper(constructSetRouterGatewayRequestOpts)(createOpts, &ExtraOption{ParentID: routerId})
	defer fasthttp.ReleaseResponse(res)

	var router entity.RouterMap
	_ = json.Unmarshal(res.Body(), &router)

	log.Println("==============Set router gateway success", router.Id)
}

func AddRouterInterfaceHelper(routerId, subnetId string) {
    opts := defaultRouterInterfaceOpts(routerId, subnetId)
	res := wrapper(constructRouterInterfaceRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var routerInterface entity.RouterInterface
	_ = json.Unmarshal(res.Body(), &routerInterface)

	log.Println("==============Add router interface success")
}


func makeSureInstanceActive(instanceId string) {
	instance := GetInstanceDetail(instanceId)
	if instance == nil {
		time.Sleep(2 * time.Second)
		for instance == nil {
			instance = GetInstanceDetail(instanceId)
			time.Sleep(2 * time.Second)
		}
	}
	timeout := 2 * 60 * time.Second
	done := make(chan bool, 1)
	go func() {
		state := instance.Server.Status
		for state != "ACTIVE" {
			time.Sleep(10 * time.Second)
			instance = GetInstanceDetail(instanceId)
			state =instance.Server.Status
		}
		done <- true
	}()
	select {
	case <-done:
		log.Println("*******************Create instance success")
	case <-time.After(timeout):
		log.Println("*******************Create instance timeout")
	}
}

func CreateInstanceHelper(netId, sgName string) string {
	EnsureSgExist(sgName)
	opts := defaultInstanceOpts(netId, sgName)
	res := wrapper(constructInstanceRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var server entity.ServerMap
	_ = json.Unmarshal(res.Body(), &server)
	makeSureInstanceActive(server.Id)
	return server.Id
}

func GetInstanceDetail(instanceId string) *entity.ServerMap {
	res := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		ParentID: "", Resource: consts.SERVER,
		ResourceLocation: fmt.Sprintf("%s/%s", consts.SERVERS, instanceId),
		ResourceSuffix: ""})
	defer fasthttp.ReleaseResponse(res)

	var server entity.ServerMap
	_ = json.Unmarshal(res.Body(), &server)
	log.Println("==============List server success", string(res.Body()))
	return &server
}


func GetSgsByName(sgName string) *entity.Sgs {
    res := wrapper(constructListRequestOpts)(nil, &ExtraOption{
    	ParentID: "", Resource: consts.SECURITYGROUP,
    	ResourceLocation: strings.Replace(consts.SECURITYGROUPS, "_", "-", 1),
        ResourceSuffix: fmt.Sprintf("name=%s", sgName)})
    defer fasthttp.ReleaseResponse(res)

	var sgs entity.Sgs
	_ = json.Unmarshal(res.Body(), &sgs)
	log.Println("==============List sgs success", sgs)
	return &sgs
}

func EnsureSgExist(sgName string) {
    sgs := GetSgsByName(sgName)
    if len(sgs.Sgs) == 0 {
    	sgId := CreateSecurityGroupHelper()
    	CreateSecurityRuleICMP(sgId)
    	CreateSecurityRuleSSH(sgId)
	}
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

