package service

import (
	"go-openstackclient/consts"
	"log"
)


var ResourceDependencies = map[string][]string{
	consts.SECURITYGROUP: []string{},
	consts.QOS_POLICY: []string{},
	consts.ROUTER: []string{},
	consts.VOLUME: []string{},
	consts.SECURITYGROUPRULE: []string{consts.SECURITYGROUP},
	consts.BANDWIDTH_LIMIT_RULE: []string{consts.QOS_POLICY},
	consts.DSCP_MARKING_RULE: []string{consts.QOS_POLICY},
	consts.MINIMUM_BANDWIDTH_RULE: []string{consts.QOS_POLICY},
	consts.NETWORK: []string{consts.QOS_POLICY},
	consts.SUBNET: []string{consts.NETWORK},
	consts.PORT: []string{consts.SUBNET, consts.SECURITYGROUP, consts.QOS_POLICY},
	consts.ROUTERINTERFACE: []string{consts.ROUTER, consts.PORT},
	consts.ROUTERGATEWAY: []string{consts.ROUTER, consts.PORT},
	consts.ROUTERROUTE: []string{consts.ROUTERINTERFACE, consts.ROUTERGATEWAY},
	consts.SNAPSHOT: []string{consts.VOLUME},
	consts.SERVER: []string{consts.SECURITYGROUP, consts.PORT, consts.VOLUME},
	consts.FLOATINGIP: []string{consts.SERVER, consts.ROUTERGATEWAY, consts.ROUTERINTERFACE},
	consts.PORTFORWARDING: []string{consts.FLOATINGIP},
	consts.FIREWALLRULE: []string{consts.PORT},
	consts.FIREWALLPOLICY: []string{consts.FIREWALLRULE},
	consts.FIREWALL: []string{consts.FIREWALLPOLICY, consts.ROUTER},
	consts.VpcConnection: []string{consts.ROUTERINTERFACE, consts.ROUTERGATEWAY, consts.FIREWALL},
	consts.LOADBALANCER: []string{consts.SUBNET},
	consts.LISTENER: []string{consts.LOADBALANCER},
	consts.POOL: []string{consts.LISTENER},
	consts.MEMBER: []string{consts.POOL},
	consts.HEALTHMONITOR: []string{consts.POOL},
	consts.L7POLICY: []string{},
	consts.L7RULE: []string{consts.L7POLICY},
}

var OrderResources = [...]string{
	consts.SECURITYGROUP,
	consts.QOS_POLICY,
	consts.ROUTER,
	consts.VOLUME,
	consts.SECURITYGROUPRULE,
	consts.BANDWIDTH_LIMIT_RULE,
	consts.DSCP_MARKING_RULE,
	consts.MINIMUM_BANDWIDTH_RULE,
	consts.NETWORK,
	consts.SUBNET,
	consts.PORT,
	consts.ROUTERINTERFACE,
	consts.ROUTERGATEWAY,
	consts.ROUTERROUTE,
	consts.SNAPSHOT,
	consts.SERVER,
	consts.FLOATINGIP,
	consts.PORTFORWARDING,
	consts.FIREWALLRULE,
	consts.FIREWALLPOLICY,
	consts.FIREWALL,
	consts.VpcConnection,
	consts.LOADBALANCER,
	consts.LISTENER,
	consts.POOL,
	consts.MEMBER,
	consts.HEALTHMONITOR,
	consts.L7POLICY,
	consts.L7RULE,
}

type Dependency struct {
	resourceName        string
	dependents          []string
	requiredBy          []string
    beDependentChannel  chan struct{}
}

type Node struct {
	resourceType             string
	dependencies             []*Node
	requiredBy               []*Node
	monitorDeleteChannel     chan struct{}
	monitorCreateChannel     chan struct{}
	resources                []string
}

func InitNodes() map[string]*Node {
	nodeMap := make(map[string]*Node)
	for _, resource := range OrderResources {
		depNodes := make([]*Node, 0)
		deps := ResourceDependencies[resource]
		if len(deps) > 0 {
			for _, resourceType := range deps {
				if node, ok := nodeMap[resourceType]; ok {
					depNodes = append(depNodes, node)
				}
			}
		}
		nodeMap[resource] = &Node{
			resourceType: resource,
			dependencies: depNodes,
			resources: make([]string, 0),
		}
	}

	for resourceType, node := range nodeMap {
		for _, dep := range node.dependencies {
			dep.requiredBy = append(dep.requiredBy, nodeMap[resourceType])
			nodeMap[dep.resourceType] = dep
		}
	}

	for _, node := range nodeMap {
		node.monitorDeleteChannel = make(chan struct{}, len(node.requiredBy))
		node.monitorCreateChannel = make(chan struct{}, len(node.dependencies))
	}

	return nodeMap
}

func nodeAssociateResources(nodeMap map[string]*Node) map[string]*Node {
	nodes := make(map[string]*Node)
	for resourceType, node := range nodeMap {
		nodes[resourceType] = nodeMap[resourceType]
		for key, resource := range ResourcesMap {
			if node.resourceType == resource.Type {
				nodes[resourceType].resources = append(nodes[resourceType].resources, key)
			}
		}
	}
	return nodes
}

func initDependencies() map[string]Dependency {
	dependencies := make(map[string]Dependency)
	beDependentMap := make(map[string][]string)
	for key, resource := range ResourcesMap {
		dependency := Dependency{
			resourceName: key,
			dependents: make([]string, 0),
		}
		for dep, _ := range resource.Dependencies {
			beDependentMap[dep] = append(beDependentMap[dep], key)
			dependency.dependents = append(dependency.dependents, dep)
		}
		dependencies[key] = dependency
	}

	for key, res := range ResourcesMap {
		dependency := dependencies[key]
		if deps, ok := beDependentMap[key]; ok {
			dependency.requiredBy = deps
		}
		if res.Type == consts.NETWORK {
			channelLen := 0
			for _, dep := range beDependentMap[key] {
				if ResourcesMap[dep].Type == consts.SUBNET {
					channelLen++
				}
			}
			dependency.beDependentChannel = make(chan struct{}, channelLen)
		} else {
			dependency.beDependentChannel = make(chan struct{}, len(dependency.requiredBy))
		}
		dependencies[key] = dependency
		log.Printf("%+v===", dependency)
		log.Printf("%+v===", cap(dependency.beDependentChannel))
	}
	return dependencies
}
