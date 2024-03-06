package service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
	"go-openstackclient/internal/entity"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var actionMapMethod = map[string]string{
	CREATE: fasthttp.MethodPost,
	UPDATE: fasthttp.MethodPut,
	GET: fasthttp.MethodGet,
	LIST: fasthttp.MethodGet,
	DELETE: fasthttp.MethodDelete,
}

type Controller struct {
	neutron           *Neutron
	keystone          *Keystone
	nova              *Nova
	cinder            *Cinder
	token             string
	projectName       string
	projectID         string
	DeleteChannels    map[string]chan Output
	mu                sync.Mutex
}

func NewController(projectName string) *Controller {
    return &Controller{
    	projectName: projectName,
	}
}

func wrapper(fn func(options entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption) func(options entity.CreateUpdateOptions, extraOption *ExtraOption) *fasthttp.Response {
	return func(param entity.CreateUpdateOptions, extraOpt *ExtraOption) *fasthttp.Response {
		opts := fn(param, extraOpt)
		resp := defaultController.Do(opts)
		return resp
	}
}

func (s *Controller) GetProjectId(projectName string) string {
	opts := &ExtraOption{
		Resource: consts.PROJECT, ResourceLocation: consts.PROJECTS,
		ResourceSuffix: fmt.Sprintf("name=%s", projectName),
	}
	res := wrapper(constructListRequestOpts)(nil, opts)
	defer fasthttp.ReleaseResponse(res)

	var projects entity.Projects
	_ = json.Unmarshal(res.Body(), &projects)
	if len(projects.Ps) == 0 {
		return ""
	}
	return projects.Ps[0].Id
}

func (s *Controller) Project() string {
    if len(s.projectID) == 0 {
    	s.projectID = s.GetProjectId(s.projectName)
	}
	return s.projectID
}

func (s *Controller) Do(option RequestOption) *fasthttp.Response {
	req := s.buildRequest(option)
    service := s.targetService(option)

    resp := service.Call(req)
    return resp
}

func (s *Controller) Token() string {
	if len(s.token) == 0 {
		createOpts := &entity.AuthOption{
			Auth: entity.Auth{
				Identity: entity.Identity{
					Methods: []string{"password"},
					Password: entity.Password{
						Userr: entity.Userr{
							Name: configs.CONF.UserName,
							Password: configs.CONF.UserPassword,
							Domain: entity.Domain{Name: "default"},
						},
					},
				},
				Scope: entity.Scope{
					Projectt: entity.Projectt{
						Name: configs.CONF.ProjectName,
						Domain: entity.Domain{Name: "default"},
					},
				},
			},
		}
		opts := RequestOption{
			Action: CREATE,
			Resource: consts.TOKEN,
			ResourceLocation: "auth/tokens",
			RequestSuffix: "",
			Body: createOpts,
			Headers: make(map[string]string),
		}
		resp := defaultController.Do(opts)
		defer fasthttp.ReleaseResponse(resp)

		token := string(resp.Header.Peek("X-Subject-Token"))
		log.Println("==============Get auth token success")
		s.token = token
	}
	return s.token
}

func (s *Controller) buildRequest(option RequestOption) client.Request {
	method := s.actionMapMethod(option.Action)
	endpoint := s.getEndpoint(option)
	service := s.targetService(option)
	var req client.Request
	if option.Body != nil {
		reqBody := option.Body.ToRequestBody()
		req = service.Client().NewRequest(
			endpoint, method, consts.ContentTypeJson, option.Headers, reqBody)
	} else {
		req = service.Client().NewRequest(
			endpoint, method, consts.ContentTypeJson, option.Headers, nil)
	}
	return req
}

func (s *Controller) targetService(option RequestOption) Service {
	var ok bool
    if _, ok = s.Neutron().SupportedResources()[option.Resource]; ok {
    	return s.neutron
	} else if _, ok = s.Keystone().SupportedResources()[option.Resource]; ok {
		return s.keystone
	} else if _, ok = s.Nova().SupportedResources()[option.Resource]; ok {
		return s.nova
	} else if _, ok = s.Cinder().SupportedResources()[option.Resource]; ok {
		return s.cinder
	} else {
		return nil
	}
}

func (s *Controller) Neutron() *Neutron {
	if s.neutron == nil {
		s.neutron = newNeutron()
	}
	return s.neutron
}

func (s *Controller) Nova() *Nova {
	if s.nova == nil {
		s.nova = newNova()
	}
	return s.nova
}

func (s *Controller) Keystone() *Keystone {
	if s.keystone == nil {
		s.keystone = newKeystone()
	}
	return s.keystone
}

func (s *Controller) Cinder() *Cinder {
	if s.cinder == nil {
		s.cinder = newCinder()
	}
	return s.cinder
}

func (s *Controller) actionMapMethod(action string) string {
	if method, ok := actionMapMethod[action]; ok {
		return method
	} else {
		panic("The action not supported")
	}
}

func (s *Controller) getEndpoint(option RequestOption) string {
     service := s.targetService(option)
     url := fmt.Sprintf("%s/%s", service.HttpPrefix(), option.ResourceLocation)
     if len(option.RequestSuffix) != 0 {
     	url = fmt.Sprintf("%s?%s", url, option.RequestSuffix)
	 }
	 return url
}


func (s *Controller) MakeDeleteChannel(resourceType string, length int) chan Output {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DeleteChannels[resourceType] = make(chan Output, length)
	return s.DeleteChannels[resourceType]
}

// network

func (s *Controller) CreateNetwork(opts entity.CreateUpdateOptions) entity.NetworkMap {
	res := wrapper(constructNetworkRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)
	var network entity.NetworkMap
	_ = json.Unmarshal(res.Body(), &network)

	log.Println("==============Create internal network success", network.Network.Id)
	return network
}

func (s *Controller) ListNetworks() entity.Networks {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.NETWORK, ResourceLocation: consts.NETWORKS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var networks entity.Networks
	_ = json.Unmarshal(resp.Body(), &networks)
	log.Println("==============List network success, there had", networks.Count)
	return networks
}

func (s *Controller) DeleteNetwork(ipId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"network_id": ipId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.NETWORK, ResourceLocation: fmt.Sprintf("%s/%s", consts.NETWORKS, ipId)})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteNetworks() {
	networks := s.ListNetworks()
	ch := s.MakeDeleteChannel(consts.NETWORK, len(networks.Nets))

	for _, network := range networks.Nets {
		tempNetwork := network
		go func() {
			ch <- s.DeleteNetwork(tempNetwork.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Networks were deleted completely")
}

// subnet

func (s *Controller) CreateSubnet(opts entity.CreateUpdateOptions) entity.SubnetMap {
	res := wrapper(constructSubnetRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)
	var subnet entity.SubnetMap
	_ = json.Unmarshal(res.Body(), &subnet)

	log.Println("==============Create internal subnet success", subnet.Id)
	return subnet
}


func (s *Controller) ListSubnet() entity.Subnets {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SUBNET, ResourceLocation: consts.SUBNETS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var subnets entity.Subnets
	_ = json.Unmarshal(resp.Body(), &subnets)
	log.Println("==============List subnet success")
	return subnets
}

func (s *Controller) DeleteSubnet(ipId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"subnet_id": ipId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("%s/%s", consts.SUBNETS, ipId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.NETWORK, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteSubnets() {
	subnets := s.ListSubnet()
	ch := s.MakeDeleteChannel(consts.SUBNET, len(subnets.Ss))

	for _, subnet := range subnets.Ss {
		tempSubnet := subnet
		go func() {
			ch <- s.DeleteSubnet(tempSubnet.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Subnets were deleted completely")
}

// security group

func (s *Controller) CreateSecurityGroup(opts entity.CreateUpdateOptions) entity.Sg {
	res := wrapper(constructSgRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var sg entity.Sg
	_ = json.Unmarshal(res.Body(), &sg)

	log.Println("==============Create security group success", sg.Id)
	return sg.Id
}


func (s *Controller) GetSgsByName(sgName string) *entity.Sgs {
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

func (s *Controller) EnsureSgExist(sgName string) {
	sgs := s.GetSgsByName(sgName)
	if len(sgs.Sgs) == 0 {
		sgId := CreateSecurityGroupHelper()
		CreateSecurityRuleICMP(sgId)
		CreateSecurityRuleSSH(sgId)
	}
}

func (s *Controller) getSecurityGroup(sgId string) interface{} {
	urlSuffix := fmt.Sprintf("security-groups/%s", sgId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SECURITYGROUP, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var sg entity.Sg
	_ = json.Unmarshal(resp.Body(), &sg)
	log.Println(fmt.Sprintf("get sg==%+v", sg))
	return sg
}


func (s *Controller) listSecurityGroups() entity.Sgs {
	urlSuffix := fmt.Sprintf("project_id=%s", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SECURITYGROUP, ResourceLocation: consts.SECURITYGROUPS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var sgs entity.Sgs
	_ = json.Unmarshal(resp.Body(), &sgs)
	log.Println("==============List sg success, there had", len(sgs.Sgs))
	return sgs
}

func (s *Controller) deleteSecurityGroup(id string) Output {
	outputObj := Output{ParametersMap: map[string]string{"security_group_id": id}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.SECURITYGROUP,
		ResourceLocation: fmt.Sprintf("%s/%s", strings.Replace(consts.SECURITYGROUPS, "_", "-", 1), id)})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteSecurityGroups() {
	sgs := s.listSecurityGroups()
	ch := s.MakeDeleteChannel(consts.SECURITYGROUP, len(sgs.Sgs))
	for _, sg := range sgs.Sgs {
		tempSg := sg
		go func() {
			ch <- s.deleteSecurityGroup(tempSg.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Security groups were deleted completely")
}


// security group rule

func (s *Controller) listSecurityGroupRules() entity.SgRules {
	urlSuffix := fmt.Sprintf("project_id=%s", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SECURITYGROUPRULE, ResourceLocation: consts.SECURITYGROUPRULES,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var sgRules entity.SgRules
	_ = json.Unmarshal(resp.Body(), &sgRules)

	log.Println("==============List sg rule success, there had", len(sgRules.Srs))
	return sgRules
}

func (s *Controller) deleteSecurityGroupRule(id string) Output {
	outputObj := Output{ParametersMap: map[string]string{"security_group_id": id}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("security-group-rules/%s", id)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SECURITYGROUPRULE, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteSecurityGroupRules() {
	sgRules := s.listSecurityGroupRules()
	ch := s.MakeDeleteChannel(consts.SECURITYGROUPRULE, len(sgRules.Srs))
	for _, sgRule := range sgRules.Srs {
		tempSgRule := sgRule
		go func() {
			ch <- s.deleteSecurityGroupRule(tempSgRule.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Security group rules were deleted completely")
}

func (s *Controller) CreateSecurityRule(opts entity.CreateUpdateOptions) {
	res := wrapper(constructSgRuleRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var sgRule entity.SgRule
	_ = json.Unmarshal(res.Body(), &sgRule)

	log.Println("==============Create security group rule success", sgRule.Id)
}

// router

func (s *Controller) CreateRouter(opts entity.CreateUpdateOptions) entity.RouterMap {
	res := wrapper(constructRouterRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var router entity.RouterMap
	_ = json.Unmarshal(res.Body(), &router)

	log.Println("==============Create router success", router.Id)
	return router
}

func (s *Controller) SetRouterGateway(opts entity.CreateUpdateOptions, routerId string) entity.RouterMap {
	res := wrapper(constructSetRouterGatewayRequestOpts)(opts, &ExtraOption{ParentID: routerId})
	defer fasthttp.ReleaseResponse(res)

	var router entity.RouterMap
	_ = json.Unmarshal(res.Body(), &router)

	log.Println("==============Set router gateway success", router.Id)
	return router
}

func (s *Controller) AddRouterInterface(opts entity.CreateUpdateOptions) {
	res := wrapper(constructRouterInterfaceRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var routerInterface entity.RouterInterface
	_ = json.Unmarshal(res.Body(), &routerInterface)

	log.Println("==============Add router interface success")
}

func (s *Controller) RemoveRouterInterface(routerId, subnetId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"router_id": routerId, "subnetId": subnetId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
			log.Println("==============Remove router interface failed")
		} else {
			log.Println("==============Remove router interface success")
		}
	}()
	urlSuffix := fmt.Sprintf("routers/%s/remove_router_interface", routerId)
	resp := wrapper(constructRouterInterfaceRequestOpts)(
		defaultRouterInterfaceOpts(routerId, subnetId), &ExtraOption{
			Resource: consts.ROUTER, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.Body()
	outputObj.Success = true
	return outputObj
}

func (s *Controller) DeleteRouterInterfaces() {
	interfacePorts := s.listRouterInterfacePorts()
	ch := s.MakeDeleteChannel(consts.ROUTERINTERFACE, len(interfacePorts.Ps))

	for _, port := range interfacePorts.Ps {
		routerId := port.DeviceId
		fixedIps := port.FixedIps
		for _, fixedIp := range fixedIps {
			tempFixedIp := fixedIp
			go func() {
				ch <- s.RemoveRouterInterface(routerId, tempFixedIp.SubnetId)
			}()
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Router interfaces were deleted completely")
}

func (s *Controller) ListRouters() entity.Routers {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.ROUTER, ResourceLocation: consts.ROUTERS, ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var routers entity.Routers
	_ = json.Unmarshal(resp.Body(), &routers)
	log.Println("==============List routers success, there had", routers.Count)
	return routers
}

func (s *Controller) listRouterInterfacePorts() entity.Ports {
	urlSuffix := fmt.Sprintf("device_owner=network:router_interface&project_id=%s", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.ROUTERINTERFACE, ResourceLocation: consts.PORTS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var ports entity.Ports
	_ = json.Unmarshal(resp.Body(), &ports)
	log.Println("==============List router interface port success, there had", ports.Count)
	return ports
}


func (s *Controller) updateRouterNoRoutes(id string) Output {
	outputObj := Output{ParametersMap: map[string]string{"router_id": id}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	opts := &entity.UpdateRouterOpts{Routes: new([]entity.Route)}
	resp := wrapper(constructRouterRequestOpts)(opts, &ExtraOption{
		Resource: consts.ROUTER, ResourceLocation: fmt.Sprintf("%s/%s", consts.ROUTERS, id)})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = ""
	outputObj.Success = true
	return outputObj
}

func (s *Controller) DeleteRouterRoutes() {
	routers := s.ListRouters()
	length := 0
	for _, router := range routers.Rs {
		if len(router.Routes) != 0 {
			length++
		}

	}
	ch := s.MakeDeleteChannel(consts.ROUTERROUTE, length)

	for _, router := range routers.Rs {
		if len(router.Routes) != 0 {
			go func() {
				ch <- s.updateRouterNoRoutes(router.Id)
			}()
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Router routes were deleted completely")
}

func (s *Controller) DeleteRouter(routerId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"router_id": routerId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
			log.Println("==============Delete router failed", routerId)
		} else {
			log.Println("==============Delete router success", routerId)
		}
	}()

	urlSuffix := fmt.Sprintf("%s/%s", consts.ROUTERS, routerId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.ROUTER, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteRouters() {
	routers := s.ListRouters()
	ch := s.MakeDeleteChannel(consts.ROUTER, len(routers.Rs))
	for _, router := range routers.Rs {
		tempRouter := router
		go func() {
			ch <- s.DeleteRouter(tempRouter.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Routers were deleted completely")
}

func (s *Controller) ClearRouterGateway(routerId, extNetId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"router_id": routerId, "ext_net_id": extNetId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()
	urlSuffix := fmt.Sprintf("routers/%s", routerId)
	opts := &entity.UpdateRouterOpts{
		GatewayInfo: &entity.GatewayInfo{}}

	resp := wrapper(constructRouterRequestOpts)(opts, &ExtraOption{
		ParentID: routerId, Resource: consts.ROUTER, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Success = true
	outputObj.Response = ""
	log.Println("==============Clear router gateway success, router", routerId)
	return outputObj
}

func (s *Controller) DeleteRouterGateways() {
	routers := s.ListRouters()
	var length int
	for _, router := range routers.Rs {
		if !reflect.DeepEqual(router.GatewayInfo, nil) {
			length++
		}
	}
	ch := s.MakeDeleteChannel(consts.ROUTERGATEWAY, length)
	for _, router := range routers.Rs {
		if !reflect.DeepEqual(router.GatewayInfo, nil) {
			go func() {
				ch <- s.ClearRouterGateway(router.Id, router.GatewayInfo.NetworkID)
			}()
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Router gateways were deleted completely")
}

func (s *Controller) GetRouter(routerId string) entity.RouterMap {
	urlSuffix := fmt.Sprintf("routers/%s", routerId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.ROUTER, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var router entity.RouterMap
	_ = json.Unmarshal(resp.Body(), &router)
	log.Println("==============Get router success", routerId)
	return router
}

// server

func (s *Controller) makeSureInstanceActive(instanceId string) {
	instance := s.GetInstanceDetail(instanceId)
	if instance == nil {
		time.Sleep(2 * time.Second)
		for instance == nil {
			instance = s.GetInstanceDetail(instanceId)
			time.Sleep(2 * time.Second)
		}
	}
	timeout := 2 * 60 * time.Second
	done := make(chan bool, 1)
	go func() {
		state := instance.Server.Status
		for state != "ACTIVE" {
			time.Sleep(10 * time.Second)
			instance = s.GetInstanceDetail(instanceId)
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

func (s *Controller) CreateInstance(opts entity.CreateUpdateOptions) entity.ServerMap {
	res := wrapper(constructInstanceRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var server entity.ServerMap
	_ = json.Unmarshal(res.Body(), &server)
	s.makeSureInstanceActive(server.Id)
	return server
}

func (s *Controller) GetInstanceDetail(instanceId string) *entity.ServerMap {
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

// port

func (s *Controller) CreatePort(opts entity.CreateUpdateOptions) entity.PortMap {
	res := wrapper(constructPortRequestOpts)(opts, nil)
	defer fasthttp.ReleaseResponse(res)

	var port entity.PortMap
	_ = json.Unmarshal(res.Body(), &port)

	log.Println("==============Create internal port success", port.Id)
	return port
}

func (s *Controller) GetPort(portId string) entity.PortMap {
	urlSuffix := fmt.Sprintf("ports/%s", portId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.PORT, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var port entity.PortMap
	_ = json.Unmarshal(resp.Body(), &port)
	log.Println("==============Get port success", portId)
	return port
}

func (s *Controller) GetPortIP(portId string) string {
	port := s.GetPort(portId)
	return port.FixedIps[0].IpAddress
}

func (s *Controller) ListPort() entity.Ports {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.PORT, ResourceLocation: consts.PORTS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var ports entity.Ports
	_ = json.Unmarshal(resp.Body(), &ports)
	log.Println("==============List port success")
	return ports
}

func (s *Controller) DeletePort(ipId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"port_id": ipId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("%s/%s", consts.PORTS, ipId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.NETWORK, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeletePorts() {
	ports := s.ListPort()
	ch := s.MakeDeleteChannel(consts.PORT, len(ports.Ps))

	for _, port := range ports.Ps {
		tempPort := port
		go func() {
			ch <- s.DeletePort(tempPort.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Ports were deleted completely")
}


// floating ip

func (s *Controller) GetFIP(fipId string) entity.FipMap {
	urlSuffix := fmt.Sprintf("floatingips/%s", fipId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.FLOATINGIP, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var fip entity.FipMap
	_ = json.Unmarshal(resp.Body(), &fip)
	log.Println("==============Get fip success", fipId)
	return fip
}

func (s *Controller) ListFIPs() entity.Fips {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}

	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.FLOATINGIP, ResourceLocation: consts.FLOATINGIPS,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var fs entity.Fips
	_ = json.Unmarshal(resp.Body(), &fs)
	log.Println("==============List fip success, there had", fs.Count)
	return fs
}

func (s *Controller) DeleteFIP(fipId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"floatingip_id": fipId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.FLOATINGIP,
		ResourceLocation: fmt.Sprintf("%s/%s", consts.FLOATINGIPS, fipId)})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteFloatingips() {
	fips := s.ListFIPs()
	ch := s.MakeDeleteChannel(consts.FLOATINGIP, len(fips.Fs))

	for _, fip := range fips.Fs {
		tempFip := fip
		go func() {
			ch <- s.DeleteFIP(tempFip.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Floatingips were deleted completely")
}

// port forwarding

func (s *Controller) GetPortForwarding(fipId string, pfId string) entity.PortForwardingMap {
	urlSuffix := fmt.Sprintf("%s/%s/%s/%s", consts.FLOATINGIPS, fipId, consts.PORTFORWARDINGS, pfId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.PORTFORWARDING, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var pf entity.PortForwardingMap
	_ = json.Unmarshal(resp.Body(), &pf)

	log.Printf("==============Get port forwarding success %+v\n", pf)
	return pf
}

func (s *Controller) ListPortForwarding(fipId string) entity.PortForwardings {
	urlSuffix := fmt.Sprintf("%s/%s/%s", consts.FLOATINGIPS, fipId, consts.PORTFORWARDINGS)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.PORTFORWARDING, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var pfs entity.PortForwardings
	_ = json.Unmarshal(resp.Body(), &pfs)

	log.Printf("==============List port forwarding success %+v\n", pfs)
	return pfs
}

func (s *Controller) DeletePortForwarding(fipId string, pfId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"floatingip_id": fipId, "port_forwarding_id": pfId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()
	urlSuffix := fmt.Sprintf("%s/%s/%s/%s", consts.FLOATINGIPS, fipId, consts.PORTFORWARDINGS, pfId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.PORTFORWARDING, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeletePortForwardings() {
	fips := s.ListFIPs()
	var pfsMap = make(map[string]entity.PortForwardings)
	var length int
	for _, fip := range fips.Fs {
		tmpPfs := s.ListPortForwarding(fip.Id)
		pfsMap[fip.Id] = tmpPfs
		length += len(tmpPfs.Pfs)
	}

	ch := s.MakeDeleteChannel(consts.PORTFORWARDING, length)
	for fipId, pfs := range pfsMap {
		for _, pf := range pfs.Pfs {
			tmpPf := pf
			go func() {
				ch <- s.DeletePortForwarding(fipId, tmpPf.Id)
			}()
		}
	}

	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Port forwarding were deleted completely")
}

// qos policy

func (s *Controller) createQosPolicy(opts entity.CreateUpdateOptions) entity.QosPolicyMap {
	res := wrapper(constructSgRuleRequestOpts)(opts, &ExtraOption{
		ResourceLocation: "qos/policies"})
	defer fasthttp.ReleaseResponse(res)

	var qosPolicy entity.QosPolicyMap
	_ = json.Unmarshal(res.Body(), &qosPolicy)
	log.Println("==============Create qos policy success", qosPolicy.Id)
	return qosPolicy
}

func (s *Controller) CreateBandwidthLimitRule(opts entity.CreateUpdateOptions, qosPolicyId string) entity.BandwidthLimitRuleMap {
	res := wrapper(constructBandwidthLimitRuleRequestOpts)(opts, &ExtraOption{ParentID: qosPolicyId})
	defer fasthttp.ReleaseResponse(res)

	var rule entity.BandwidthLimitRuleMap
	_ = json.Unmarshal(res.Body(), &rule)

	log.Println("==============Create bandwidth_limit_rule success")
	return rule
}

func (s *Controller) DeleteBandwidthLimitRules() {
	qoss := s.listQoss()
	var length int
	for _, qos := range qoss.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "bandwidth_limit" {
				length++
			}
		}
	}
	ch := s.MakeDeleteChannel(consts.BANDWIDTH_LIMIT_RULE, length)
	for _, qos := range qoss.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "bandwidth_limit" {
				go func() {
					ch <- s.DeleteQosRule(rule.Type, qos.Id, rule.Id)
				}()
			}
		}
	}

	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Bandwidth limit rules were deleted completely")
}


func (s *Controller) DeleteDscpMarkingRules() {
	qoss := s.listQoss()
	var length int
	for _, qos := range qoss.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "dscp_marking" {
				length++
			}
		}
	}
	ch := s.MakeDeleteChannel(consts.DSCP_MARKING_RULE, length)
	for _, qos := range qoss.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "dscp_marking" {
				go func() {
					ch <- s.DeleteQosRule(rule.Type, qos.Id, rule.Id)
				}()
			}
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Dscp marking rules were deleted completely")
}

func (s *Controller) DeleteMinimumBandwidthRules() {
	qoses := s.listQoss()
	var length int
	for _, qos := range qoses.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "minimum_bandwidth" {
				length++
			}
		}
	}
	ch := s.MakeDeleteChannel(consts.MINIMUM_BANDWIDTH_RULE, length)
	for _, qos := range qoses.Qps {
		rules := qos.Rules
		for _, rule := range rules {
			if rule.Type == "minimum_bandwidth" {
				go func() {
					ch <- s.DeleteQosRule(rule.Type, qos.Id, rule.Id)
				}()
			}
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Minimum bandwidth rules were deleted completely")
}

func (s *Controller) listQoss() entity.QosPolicies {
	var urlSuffix = ""
	if s.projectName != consts.ADMIN {
		urlSuffix = fmt.Sprintf("project_id=%s", s.projectID)
	}
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.QOS_POLICY, ResourceLocation: consts.QOS_POLICIES,
		ResourceSuffix: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var qos entity.QosPolicies
	_ = json.Unmarshal(resp.Body(), &qos)
	log.Println("==============List qos policy success, there had", qos.Count)
	return qos
}

func (s *Controller) DeleteQos(qosId string) Output {
	outputObj := Output{ParametersMap: map[string]string{"qos_policy_id": qosId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("qos/policies/%s", qosId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.NETWORK, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

func (s *Controller) DeleteQosPolicies() {
	qoses := s.listQoss()
	ch := s.MakeDeleteChannel(consts.QOS_POLICY, len(qoses.Qps))
	for _, qos := range qoses.Qps {
		tempQos := qos
		go func() {
			ch <- s.DeleteQos(tempQos.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Qos policies were deleted completely")
}

func (s *Controller) DeleteQosRule(ruleType, qosId, ruleId string) Output {
	var identity string
	switch ruleType {
	case "bandwidth_limit":
		identity = consts.BANDWIDTH_LIMIT_RULES
	case "dscp_marking":
		identity = consts.DSCP_MARKING_RULES
	case "minimum_bandwidth":
		identity = consts.MINIMUM_BANDWIDTH_RULES
	}
	outputObj := Output{ParametersMap: map[string]string{"qos_policy_id": qosId, "rule_id": ruleId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("qos/policies/%s/%s/%s", qosId, identity, ruleId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.QOS_POLICY, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
	return outputObj
}

// rbac policy

// vpn


// Volume type
func (s *Controller) createVolumeType(opts entity.CreateUpdateOptions) entity.VolumeType {
	PostUrl := fmt.Sprintf("%s/types", s.projectID)
	res := wrapper(constructVolumeTypeRequestOpts)(opts, &ExtraOption{
		ResourceLocation: PostUrl})
	defer fasthttp.ReleaseResponse(res)

	var volumeType entity.VolumeType
	_ = json.Unmarshal(res.Body(), &volumeType)
	return volumeType
}

func (s *Controller) VolumeTypeAssociateQos(opts entity.CreateUpdateOptions, qosId, volTypeId string) {
	URL := fmt.Sprintf("/%s/qos-specs/%s/associate?vol_type_id=%s", s.projectID, qosId, volTypeId)
	resp := wrapper(constructListRequestOpts)(opts, &ExtraOption{
		Resource: consts.QOSSPEC, ResourceLocation: URL})
	defer fasthttp.ReleaseResponse(resp)
}

func (s *Controller) GetVolumeType(volumeTypeId string) entity.VolumeType {
	suffix := fmt.Sprintf("%s/types/%s", s.projectID, volumeTypeId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUMETYPE, ResourceLocation: suffix})
	defer fasthttp.ReleaseResponse(resp)

	var volumeType entity.VolumeType
	_ = json.Unmarshal(resp.Body(), &volumeType)

	return volumeType
}

func (s *Controller) ListVolumeTypes() entity.VolumeTypes {
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUMETYPE, ResourceLocation: fmt.Sprintf("%s/types", s.projectID)})
	defer fasthttp.ReleaseResponse(resp)

	var volumeTypes entity.VolumeTypes
	_ = json.Unmarshal(resp.Body(), &volumeTypes)

	return volumeTypes
}

func (s *Controller) ListVolumeQos() entity.QosSpecss {
	urlSuffix := fmt.Sprintf("/%s/qos-specs", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.QOSSPEC, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var qss entity.QosSpecss
	_ = json.Unmarshal(resp.Body(), &qss)
	log.Println("==============Get qos specs success")
	return qss
}

func (s *Controller) deleteVolumeType(typeId string)  {
	urlSuffix := fmt.Sprintf("/%s/types/%s", s.projectID, typeId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUMETYPE, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	log.Println("success delete", typeId)
}


// volume

// CreateVolume create volume
func (s *Controller) CreateVolume(opts entity.CreateUpdateOptions) string {
	urlSuffix := fmt.Sprintf("/%s/volumes", s.projectID)
	resp := wrapper(constructCreateVolumeRequestOpts)(opts, &ExtraOption{
		ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var volume entity.VolumeMap
	_ = json.Unmarshal(resp.Body(), &volume)

	s.MakeSureVolumeAvailable(volume.Id)
	log.Println("==============Create volume success", volume.Id)
	return volume.Id
}


func (s *Controller) GetVolume(volumeId string) entity.VolumeMap {
	urlSuffix := fmt.Sprintf("%s/volumes/%s", s.projectID, volumeId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUME, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var volume entity.VolumeMap
	_ = json.Unmarshal(resp.Body(), &volume)
	log.Println("==============Get volume success", volumeId)
	return volume
}

func (s *Controller) MakeSureVolumeAvailable(volumeId string) {
	volume := s.GetVolume(volumeId)
	done := make(chan bool, 1)
	go func() {
		state := volume.Status
		for state != consts.Available && state != consts.Error {
			time.Sleep(consts.IntervalTime)
			volume = s.GetVolume(volumeId)
			state = volume.Status
		}
		done <- true
	}()
	select {
	case <-done:
		log.Println("*******************Create Volume success")
	case <-time.After(consts.Timeout):
		log.Fatalln("*******************Create volume timeout")
	}
}

func (s *Controller) DeleteVolume(volumeId string, ch chan Output) {
	outputObj := Output{ParametersMap: map[string]string{"volume_id": volumeId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
		ch <- outputObj
	}()

	urlSuffix := fmt.Sprintf("/%s/volumes/%s", s.projectID, volumeId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUME, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
}

func (s *Controller) ListVolumes() entity.Volumes {
	urlSuffix := fmt.Sprintf("/%s/volumes", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.VOLUME, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var volumes entity.Volumes
	_ = json.Unmarshal(resp.Body(), &volumes)
	log.Println("==============List volume success, there had", len(volumes.Vs))
	return volumes
}


func (s *Controller) DeleteAttachment(attachmentId string) {
	urlSuffix := fmt.Sprintf("%s/attachments/%s", s.projectID, attachmentId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.ATTACHMENT, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() == fasthttp.StatusOK {
		log.Println("==============Delete attachment success", attachmentId)
		return
	}
	log.Println("==============Delete attachment failed", attachmentId)
}

func (s *Controller) DeleteVolumes() {
	volumes := s.ListVolumes()
	ch := s.MakeDeleteChannel(consts.VOLUME, len(volumes.Vs))
	for _, volume := range volumes.Vs {
		if len(volume.Attachments) != 0 {
			for _, attachment := range volume.Attachments {
				s.DeleteAttachment(attachment.AttachmentId)
			}
		}
		go s.DeleteVolume(volume.Id, ch)
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Volumes were deleted completely")
}


// snapshot

// CreateSnapshot create snapshot from volume
func (s *Controller) CreateSnapshot(opts entity.CreateUpdateOptions) string {
	urlSuffix := fmt.Sprintf("/%s/snapshots", s.projectID)
	resp := wrapper(constructCreateSnapshotRequestOpts)(opts, &ExtraOption{
		ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var snapshot entity.SnapshotMap
	_ = json.Unmarshal(resp.Body(), &snapshot)
	s.makeSureSnapshotAvailable(snapshot.Id)
	return snapshot.Id
}

func (s *Controller) GetSnapshot(snapshotId string) entity.SnapshotMap {
	urlSuffix := fmt.Sprintf("/%s/snapshots/%s", s.projectID, snapshotId)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SNAPSHOT, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var snapshot entity.SnapshotMap
	_ = json.Unmarshal(resp.Body(), &snapshot)
	log.Println("==============Get snapshot success", snapshotId)
	return snapshot
}

func (s *Controller) makeSureSnapshotAvailable(snapshotId string) {
	snapshot := s.GetSnapshot(snapshotId)
	done := make(chan bool, 1)
	go func() {
		state := snapshot.Status
		for state != consts.Available && state != consts.Error {
			time.Sleep(consts.IntervalTime)
			snapshot = s.GetSnapshot(snapshotId)
			state = snapshot.Status
		}
		done <- true
	}()
	select {
	case <-done:
		log.Println("*******************Create snapshot success")
	case <-time.After(consts.Timeout):
		log.Fatalln("*******************Create snapshot timeout")
	}
}

func (s *Controller) listProjectSnapshots() entity.Snapshots {
	urlSuffix := fmt.Sprintf("/%s/snapshots", s.projectID)
	resp := wrapper(constructListRequestOpts)(nil, &ExtraOption{
		Resource: consts.SNAPSHOT, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	var ss entity.Snapshots
	_ = json.Unmarshal(resp.Body(), &ss)
	log.Println("==============List snapshot success")
	return ss
}

func (s *Controller) DeleteSnapshot(snapshotId string, ch chan Output) {
	outputObj := Output{ParametersMap: map[string]string{"snapshot_id": snapshotId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
		ch <- outputObj
	}()

	urlSuffix := fmt.Sprintf("/%s/snapshots/%s", s.projectID, snapshotId)
	resp := wrapper(constructDeleteRequestOpts)(nil, &ExtraOption{
		Resource: consts.SNAPSHOT, ResourceLocation: urlSuffix})
	defer fasthttp.ReleaseResponse(resp)

	outputObj.Response = resp.StatusCode()
	if resp.StatusCode() == fasthttp.StatusOK {
		outputObj.Success = true
	}
}

func (s *Controller) DeleteSnapshots() {
	snapshots := s.listProjectSnapshots()
	ch := s.MakeDeleteChannel(consts.SNAPSHOT, len(snapshots.Ss))
	for _, snapshot := range snapshots.Ss {
		go s.DeleteSnapshot(snapshot.Id, ch)
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Snapshots were deleted completely")
}
