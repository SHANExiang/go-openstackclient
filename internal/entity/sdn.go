package entity

import (
	"fmt"
	"time"
)

type CreateSDNTokenOpts struct {
	UserName            string        	`json:"userName"`
	Password            string          `json:"password"`
}

type SDNToken struct {
	Data struct {
		TokenId     string `json:"token_id"`
		ExpiredDate string `json:"expiredDate"`
	} `json:"data"`
	Errcode string `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (opts *CreateSDNTokenOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}

type SDNNetworks struct {
	HuaweiAcNeutronNetworks struct {
		Network []struct {
			Uuid                                      string `json:"uuid"`
			Mtu                                       int    `json:"mtu"`
			VlanTransparent                           bool   `json:"vlan-transparent"`
			Name                                      string `json:"name"`
			Description                               string `json:"description"`
			UpdatedAt                                 string `json:"updated-at"`
			CreatedAt                                 string `json:"created-at"`
			TenantName                                string `json:"tenant-name"`
			AdminStateUp                              bool   `json:"admin-state-up"`
			CloudName                                 string `json:"cloud-name"`
			Shared                                    bool   `json:"shared"`
			TenantId                                  string `json:"tenant-id"`
			HuaweiAcNeutronProviderExtPhysicalNetwork string `json:"huawei-ac-neutron-provider-ext:physical-network"`
			HuaweiAcNeutronProviderExtNetworkType     string `json:"huawei-ac-neutron-provider-ext:network-type"`
			HuaweiAcNeutronProviderExtSegmentationId  string `json:"huawei-ac-neutron-provider-ext:segmentation-id"`
			HuaweiAcNeutronL3ExtExternal              bool   `json:"huawei-ac-neutron-l3-ext:external"`
		} `json:"network"`
	} `json:"huawei-ac-neutron:networks"`
}

// TrafficProfile 带宽通道
type TrafficProfile struct {
	Id                         string `json:"id"`
	Name                       string `json:"name"`
	Description                string `json:"description"`
	ReferenceMode              string `json:"referenceMode"`
	FabricWholeConnectionLimit struct {
		TotalConnection int `json:"totalConnection"`
	} `json:"fabricWholeConnectionLimit"`
	FabricWholeBandWidthLimit struct {
		UpstreamBandwidth   int `json:"upstreamBandwidth"`
		DownstreamBandwidth int `json:"downstreamBandwidth"`
	} `json:"fabricWholeBandWidthLimit"`
	TenantId string `json:"tenantId"`
}

type TrafficProfiles struct {
	QosDtoList []struct {
		Id                         string `json:"id"`
		Name                       string `json:"name"`
		Description                string `json:"description"`
		ReferenceMode              string `json:"referenceMode"`
		FabricWholeConnectionLimit struct {
			TotalConnection int `json:"totalConnection"`
		} `json:"fabricWholeConnectionLimit"`
		FabricWholeBandWidthLimit struct {
			UpstreamBandwidth   int `json:"upstreamBandwidth"`
			DownstreamBandwidth int `json:"downstreamBandwidth"`
		} `json:"fabricWholeBandWidthLimit"`
		TenantId   string `json:"tenantId"`
		CreateTime string `json:"createTime"`
		UpdateTime string `json:"updateTime"`
	} `json:"qosDtoList"`
	TotalNum int `json:"totalNum"`
}

type QosRule struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Enable        bool   `json:"enable"`
	SourceAddress struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"sourceAddress"`
	DestinationAddress struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"destinationAddress"`
	Service struct {
		Protocol   string `json:"protocol"`
		Dport      string `json:"dport"`
		Sport      string `json:"sport"`
		IcmpType   string `json:"icmpType"`
		IcmpCode   string `json:"icmpCode"`
		ProtocolId int    `json:"protocolId"`
	} `json:"service"`
	TrafficProfile string `json:"trafficProfile"`
	TenantId       string `json:"tenantId"`
	LogicVasId     string `json:"logicVasId"`
}

type LogicVass struct {
	Vas []struct {
		Id                 string      `json:"id"`
		Name               string      `json:"name"`
		Description        interface{} `json:"description"`
		Type               string      `json:"type"`
		LogicNetworkId     string      `json:"logicNetworkId"`
		LogicRouterId      string      `json:"logicRouterId"`
		TenantId           string      `json:"tenantId"`
		DesignatedName     string      `json:"designatedName"`
		Ipv6Enable         bool        `json:"ipv6Enable"`
		RouteDeployDisable bool        `json:"routeDeployDisable"`
		VaspoolId          []string    `json:"vaspoolId"`
		AutoGenerateLink   interface{} `json:"autoGenerateLink"`
		Additional         struct {
			Producer string      `json:"producer"`
			CreateAt time.Time   `json:"createAt"`
			UpdateAt interface{} `json:"updateAt"`
		} `json:"additional"`
	} `json:"vas"`
	TotalNum  int `json:"totalNum"`
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}