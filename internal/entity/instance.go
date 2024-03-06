package entity

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"strings"
	"time"
)

type ServerMap struct {
	Server `json:"server"`
}

type Server struct {
	OSDCFDiskConfig                string      `json:"OS-DCF:diskConfig"`
	OSEXTAZAvailabilityZone        string      `json:"OS-EXT-AZ:availability_zone"`
	OSEXTSRVATTRHost               string      `json:"OS-EXT-SRV-ATTR:host"`
	OSEXTSRVATTRHypervisorHostname string      `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	OSEXTSRVATTRInstanceName       string      `json:"OS-EXT-SRV-ATTR:instance_name"`
	OSEXTSTSPowerState             int         `json:"OS-EXT-STS:power_state"`
	OSEXTSTSTaskState              interface{} `json:"OS-EXT-STS:task_state"`
	OSEXTSTSVmState                string      `json:"OS-EXT-STS:vm_state"`
	OSSRVUSGLaunchedAt             string      `json:"OS-SRV-USG:launched_at"`
	OSSRVUSGTerminatedAt           interface{} `json:"OS-SRV-USG:terminated_at"`
	AccessIPv4                     string      `json:"accessIPv4"`
	AccessIPv6                     string      `json:"accessIPv6"`
	Addresses                      struct {
		DxIntNet81 []struct {
			OSEXTIPSMACMacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
			OSEXTIPSType       string `json:"OS-EXT-IPS:type"`
			Addr               string `json:"addr"`
			Version            int    `json:"version"`
		} `json:"dx_int_net81"`
	} `json:"addresses"`
	ConfigDrive string    `json:"config_drive"`
	Created     time.Time `json:"created"`
	Flavor      struct {
		Id    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"flavor"`
	HostId string `json:"hostId"`
	Id     string `json:"id"`
	Image  struct {
		Id    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"image"`
	KeyName interface{} `json:"key_name"`
	Links   []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Metadata struct {
	} `json:"metadata"`
	Name                             string        `json:"name"`
	OsExtendedVolumesVolumesAttached []interface{} `json:"os-extended-volumes:volumes_attached"`
	Progress                         int           `json:"progress"`
	SecurityGroups                   []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	Status   string    `json:"status"`
	TenantId string    `json:"tenant_id"`
	Updated  time.Time `json:"updated"`
	UserId   string    `json:"user_id"`
}

type Servers struct {
	Servers            []Server `json:"servers"`
	Count              int      `json:"count"`
}

type Personality []*File

type File struct {
	// Path of the file.
	Path string

	// Contents of the file. Maximum content size is 255 bytes.
	Contents []byte
}

// MarshalJSON marshals the escaped file, base64 encoding the contents.
func (f *File) MarshalJSON() ([]byte, error) {
	file := struct {
		Path     string `json:"path"`
		Contents string `json:"contents"`
	}{
		Path:     f.Path,
		Contents: base64.StdEncoding.EncodeToString(f.Contents),
	}
	return json.Marshal(file)
}

type ServerNet struct {
	UUID          string      `json:"uuid,omitempty"`
	Port          string      `json:"port,omitempty"`
}

type ServerSg struct {
	Name          string       `json:"name"`
}

type BlockDeviceMapping struct {
	BootIndex           int    `json:"boot_index"`
	Uuid                string `json:"uuid,omitempty"`
	SourceType          string `json:"source_type,omitempty"`
	VolumeSize          int    `json:"volume_size,omitempty"`
	DestinationType     string `json:"destination_type,omitempty"`
	VolumeType          string `json:"volume_type,omitempty"`
	DeleteOnTermination bool   `json:"delete_on_termination,omitempty"`
}

// CreateInstanceOpts specifies server creation parameters.
type CreateInstanceOpts struct {
	// Name is the name to assign to the newly launched server.
	Name string `json:"name" required:"true"`

	// ImageRef is the ID or full URL to the image that contains the
	// server's OS and initial state.
	// Also optional if using the boot-from-volume extension.
	ImageRef string `json:"imageRef"`

	// FlavorRef is the ID or full URL to the flavor that describes the server's specs.
	FlavorRef string `json:"flavorRef"`

	// SecurityGroups lists the names of the security groups to which this server
	// should belong.
	SecurityGroups []ServerSg `json:"security_groups,omitempty"`

	// UserData contains configuration information or scripts to use upon launch.
	// Create will base64-encode it for you, if it isn't already.
	UserData string `json:"user_data,omitempty"`

	// AvailabilityZone in which to launch the server.
	AvailabilityZone string `json:"availability_zone,omitempty"`

	// Networks dictates how this server will be attached to available networks.
	// By default, the server will be attached to all isolated networks for the
	// tenant.
	// Starting with microversion 2.37 networks can also be an "auto" or "none"
	// string.
	Networks []ServerNet `json:"networks,omitempty"`

	// Metadata contains key-value pairs (up to 255 bytes each) to attach to the
	// server.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Personality includes files to inject into the server at launch.
	// Create will base64-encode file contents for you.
	Personality Personality `json:"personality,omitempty"`

	// ConfigDrive enables metadata injection through a configuration drive.
	ConfigDrive *bool `json:"config_drive,omitempty"`

	// AdminPass sets the root user password. If not set, a randomly-generated
	// password will be created and returned in the response.
	AdminPass string `json:"adminPass,omitempty"`

	// AccessIPv4 specifies an IPv4 address for the instance.
	AccessIPv4 string `json:"accessIPv4,omitempty"`

	// AccessIPv6 specifies an IPv6 address for the instance.
	AccessIPv6 string `json:"accessIPv6,omitempty"`

	// Min specifies Minimum number of servers to launch.
	Min int `json:"min_count,omitempty"`

	// Max specifies Maximum number of servers to launch.
	Max int `json:"max_count,omitempty"`

	// Tags allows a server to be tagged with single-word metadata.
	// Requires microversion 2.52 or later.
	Tags []string `json:"tags,omitempty"`

	BlockDeviceMappingV2       []BlockDeviceMapping           `json:"block_device_mapping_v2,omitempty"`
}


func (opts *CreateInstanceOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.SERVER)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type ComputeServices struct {
	Services []struct {
		Status         string      `json:"status"`
		Binary         string      `json:"binary"`
		Host           string      `json:"host"`
		Zone           string      `json:"zone"`
		State          string      `json:"state"`
		DisabledReason interface{} `json:"disabled_reason"`
		Id             int         `json:"id"`
		UpdatedAt      string      `json:"updated_at"`
	} `json:"services"`
}

type AggregateMap struct {
	Aggregate struct {
		AvailabilityZone string        `json:"availability_zone"`
		CreatedAt        string        `json:"created_at"`
		Deleted          bool          `json:"deleted"`
		DeletedAt        interface{}   `json:"deleted_at"`
		Hosts            []interface{} `json:"hosts"`
		Id               int           `json:"id"`
		Metadata         struct {
			AvailabilityZone string `json:"availability_zone"`
			Key              string `json:"key"`
		} `json:"metadata"`
		Name      string `json:"name"`
		UpdatedAt string `json:"updated_at"`
		Uuid      string `json:"uuid"`
	} `json:"aggregate"`
}

type CreateAggregateOpts struct {
	// The name of the host aggregate.
	Name string `json:"name" required:"true"`

	// The availability zone of the host aggregate.
	// You should use a custom availability zone rather than
	// the default returned by the os-availability-zone API.
	// The availability zone must not include ‘:’ in its name.
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts *CreateAggregateOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.AGGREGATE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type UpdateAggregateOpts struct {
	// The name of the host aggregate.
	Name string `json:"name,omitempty"`

	// The availability zone of the host aggregate.
	// You should use a custom availability zone rather than
	// the default returned by the os-availability-zone API.
	// The availability zone must not include ‘:’ in its name.
	AvailabilityZone string `json:"availability_zone,omitempty"`
}

func (opts *UpdateAggregateOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.AGGREGATE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type AddHostOpts struct {
	// The name of the host.
	Host string `json:"host" required:"true"`
}

func (opts *AddHostOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.ADDHOST)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type RemoveHostOpts struct {
	// The name of the host.
	Host string `json:"host" required:"true"`
}

func (opts *RemoveHostOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.REMOVEHOST)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type SetMetadataOpts struct {
	Metadata map[string]interface{} `json:"metadata" required:"true"`
}

func (opts *SetMetadataOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.SETMETADATA)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type ConsoleProtocol string

const (
	// ConsoleProtocolVNC represents the VNC console protocol.
	ConsoleProtocolVNC ConsoleProtocol = "vnc"

	// ConsoleProtocolSPICE represents the SPICE console protocol.
	ConsoleProtocolSPICE ConsoleProtocol = "spice"

	// ConsoleProtocolRDP represents the RDP console protocol.
	ConsoleProtocolRDP ConsoleProtocol = "rdp"

	// ConsoleProtocolSerial represents the Serial console protocol.
	ConsoleProtocolSerial ConsoleProtocol = "serial"

	// ConsoleProtocolMKS represents the MKS console protocol.
	ConsoleProtocolMKS ConsoleProtocol = "mks"
)

// ConsoleType represents valid remote console type.
// It can be used to create a remote console with one of the pre-defined type.
type ConsoleType string

const (
	// ConsoleTypeNoVNC represents the VNC console type.
	ConsoleTypeNoVNC ConsoleType = "novnc"

	// ConsoleTypeXVPVNC represents the XVP VNC console type.
	ConsoleTypeXVPVNC ConsoleType = "xvpvnc"

	// ConsoleTypeRDPHTML5 represents the RDP HTML5 console type.
	ConsoleTypeRDPHTML5 ConsoleType = "rdp-html5"

	// ConsoleTypeSPICEHTML5 represents the SPICE HTML5 console type.
	ConsoleTypeSPICEHTML5 ConsoleType = "spice-html5"

	// ConsoleTypeSerial represents the Serial console type.
	ConsoleTypeSerial ConsoleType = "serial"

	// ConsoleTypeWebMKS represents the Web MKS console type.
	ConsoleTypeWebMKS ConsoleType = "webmks"
)

type CreateConsoleOpts struct {
	// Protocol specifies the protocol of a new remote console.
	Protocol ConsoleProtocol `json:"protocol" required:"true"`

	// Type specifies the type of a new remote console.
	Type ConsoleType `json:"type" required:"true"`
}

func (opts *CreateConsoleOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.REMOTECONSOLE)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type RemoteConsoleMap struct {
	RemoteConsole struct {
		Protocol string `json:"protocol"`
		Type     string `json:"type"`
		Url      string `json:"url"`
	} `json:"remote_console"`
}

type ImageSnapshot struct {
	ImageId             string       `json:"image_id"`
}


// AssignProps second output parameter is dependent resources slice
func (opts *CreateInstanceOpts) AssignProps(props map[string]interface{}) (*CreateInstanceOpts, map[string]string) {
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
				if compareFieldName == "imageRef" && v.(string) == "local" {
					value.SetString(configs.CONF.ImageId)
				} else if compareFieldName == "flavorRef" && v.(string) == "local" {
					value.SetString(configs.CONF.FlavorId)
				} else {
					value.SetString(v.(string))
				}
			case reflect.Bool:
				value.SetBool(v.(bool))
			case reflect.Ptr:
				ptr := reflect.New(fieldType.Elem())
				ptr.Elem().Set(reflect.ValueOf(v))
				value.Set(ptr)
			case reflect.Slice:
				if compareFieldName == "networks" {
					originSlice := make([]ServerNet, value.Len())
					for _, item := range v.([]interface{}) {
						uuid := item.(map[string]interface{})["uuid"].(string)
						if !IsUUID(uuid) {
							deps[uuid] = field.Name + "/UUID"
						} else {
							obj := ServerNet{UUID: uuid}
							originSlice = append(originSlice, obj)
						}
					}
					value.Set(reflect.ValueOf(originSlice))
				}
				if compareFieldName == "security_groups" {
					originSlice := make([]ServerSg, value.Len())
					for _, item := range v.([]interface{}) {
						name := item.(map[string]interface{})["name"].(string)
						obj := ServerSg{name}
						originSlice = append(originSlice, obj)
					}
					value.Set(reflect.ValueOf(originSlice))
				}
				if compareFieldName == "block_device_mapping_v2" {
					originSlice := make([]BlockDeviceMapping, value.Len())
					for _, item := range v.([]interface{}) {
						obj := BlockDeviceMapping{}
						data, _ := json.Marshal(item)
						_ = json.Unmarshal(data, &obj)
						if sourceType, ok := item.(map[string]interface{})["source_type"]; ok {
							if sourceType.(string) == "image" && obj.Uuid == "local" {
								obj.Uuid = configs.CONF.ImageId
							} else {
								obj.Uuid = sourceType.(string)
							}
						}
						originSlice = append(originSlice, obj)
					}
					value.Set(reflect.ValueOf(originSlice))
				}
			}
		}
	}
	return opts, deps
}