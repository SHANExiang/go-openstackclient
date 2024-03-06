package utils

import (
	"fmt"
	"regexp"
)

func ParseFields() {
	testStr := "\n\tadmin_state_up (Optional)\n\n\tbody\n\n\tboolean\n\n\tThe administrative state of the network, which is up (true) or down (false).\n\n\tdns_domain (Optional)\n\n\tbody\n\n\tstring\n\n\tA valid DNS domain.\n\n\tmtu (Optional)\n\n\tbody\n\n\tinteger\n\n\tThe maximum transmission unit (MTU) value to address fragmentation. Minimum value is 68 for IPv4, and 1280 for IPv6.\n\n\tname (Optional)\n\n\tbody\n\n\tstring\n\n\tHuman-readable name of the network.\n\n\tport_security_enabled (Optional)\n\n\tbody\n\n\tboolean\n\n\tThe port security status of the network. Valid values are enabled (true) and disabled (false). This value is used as the default value of port_security_enabled field of a newly created port.\n\n\tproject_id (Optional)\n\n\tbody\n\n\tstring\n\n\tThe ID of the project that owns the resource. Only administrative and users with advsvc role can specify a project ID other than their own. You cannot change this value through authorization policies.\n\n\tprovider:network_type (Optional)\n\n\tbody\n\n\tstring\n\n\tThe type of physical network that this network should be mapped to. For example, flat, vlan, vxlan, or gre. Valid values depend on a networking back-end.\n\n\tprovider:physical_network (Optional)\n\n\tbody\n\n\tstring\n\n\tThe physical network where this network should be implemented. The Networking API v2.0 does not provide a way to list available physical networks. For example, the Open vSwitch plug-in configuration file defines a symbolic name that maps to specific bridges on each compute host.\n\n\tprovider:segmentation_id (Optional)\n\n\tbody\n\n\tinteger\n\n\tThe ID of the isolated segment on the physical network. The network_type attribute defines the segmentation model. For example, if the network_type value is vlan, this ID is a vlan identifier. If the network_type value is gre, this ID is a gre key.\n\n\tqos_policy_id (Optional)\n\n\tbody\n\n\tstring\n\n\tThe ID of the QoS policy associated with the network.\n\n\trouter:external (Optional)\n\n\tbody\n\n\tboolean\n\n\tIndicates whether the network has an external routing facility that’s not managed by the networking service.\n\n\tsegments (Optional)\n\n\tbody\n\n\tarray\n\n\tA list of provider segment objects.\n\n\tshared (Optional)\n\n\tbody\n\n\tboolean\n\n\tIndicates whether this resource is shared across all projects. By default, only administrative users can change this value.\n\n\ttenant_id (Optional)\n\n\tbody\n\n\tstring\n\n\tThe ID of the project that owns the resource. Only administrative and users with advsvc role can specify a project ID other than their own. You cannot change this value through authorization policies.\n\n\tvlan_transparent (Optional)\n\n\tbody\n\n\tboolean\n\n\tIndicates the VLAN transparency mode of the network, which is VLAN transparent (true) or not VLAN transparent (false).\n\n\tdescription (Optional)\n\n\tbody\n\n\tstring\n\n\tA human-readable description for the resource. Default is an empty string.\n\n\tis_default (Optional)\n\n\tbody\n\n\tboolean\n\n\tThe network is default or not.\n\n\tavailability_zone_hints (Optional)"

	pattern := `\n\t(.*?)\s\(Optional\)`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindAllStringSubmatch(testStr, -1)
	for _, match := range matches {
		fmt.Println(match[1])
	}
}
