package configs


type Openstack struct {
	Host                  string
	AdminPassword         string
	ProjectName           string
	UserName              string
	UserPassword          string
	ExternalNetwork       string
	ImageId               string
	FlavorId              string
}

type SDN struct {
	SDNHost               string
	SDNUserName           string
	SDNPassword           string
}

type Server struct {
	Openstack
	SDN
}

var CONF  Server
