package service

import (
	"fmt"
	"go-openstackclient/consts"
	"go-openstackclient/internal/client"
	"go-openstackclient/internal/entity"
	"strings"
)

type HookOpts struct {
	ConstructReq  func(opts interface{}) client.Request
}

type RequestOption struct {
	Action                string
	Resource              string
    ResourceLocation      string
    RequestSuffix         string
	Body                  entity.CreateUpdateOptions
	Headers               map[string]string
}

type ExtraOption struct {
	ParentID                 string
	Resource                 string
	ResourceLocation         string
    ResourceSuffix           string
}


func constructNetworkRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	netOpts := RequestOption{
		Action: CREATE,
		Resource: consts.NETWORK,
		ResourceLocation: consts.NETWORKS,
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
	return netOpts
}

func constructSubnetRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	requestOpts := RequestOption{
		Action: CREATE,
		Resource: consts.SUBNET,
		ResourceLocation: consts.SUBNETS,
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
	return requestOpts
}

func constructPortRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.PORT,
		ResourceLocation: consts.PORTS,
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructSgRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.SECURITYGROUP,
		ResourceLocation: strings.Replace(consts.SECURITYGROUPS, "_", "-", 1),
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructRouterRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.ROUTER,
		ResourceLocation: consts.ROUTERS,
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructRouterInterfaceRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: UPDATE,
		Resource: consts.ROUTER,
		ResourceLocation: fmt.Sprintf("%s/%s/add_router_interface", consts.ROUTERS, opts.(*entity.AddRouterInterfaceOpts).RouterId),
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructSetRouterGatewayRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: UPDATE,
		Resource: consts.ROUTER,
		ResourceLocation: fmt.Sprintf("%s/%s", consts.ROUTERS, extraOption.ParentID),
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructSgRuleRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.SECURITYGROUPRULE,
		ResourceLocation: strings.Replace(consts.SECURITYGROUPRULES, "_", "-", 1),
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructInstanceRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.SERVER,
		ResourceLocation: consts.SERVERS,
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructUpdateRouterRequestOpts(opts entity.CreateUpdateOptions, extraOption *ExtraOption) RequestOption {
	return RequestOption{
		Action: UPDATE,
		Resource: consts.ROUTER,
		ResourceLocation: fmt.Sprintf("%s/%s", consts.ROUTERS, extraOption.ParentID),
		RequestSuffix: "",
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}


func constructListRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: GET,
		Resource: extraOpts.Resource,
		ResourceLocation: extraOpts.ResourceLocation,
		RequestSuffix: extraOpts.ResourceSuffix,
		Body: nil,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructDeleteRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: DELETE,
		Resource: extraOpts.Resource,
		ResourceLocation: extraOpts.ResourceLocation,
		RequestSuffix: extraOpts.ResourceSuffix,
		Body: nil,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructVolumeTypeRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.VOLUME,
		ResourceLocation: extraOpts.ResourceLocation,
		RequestSuffix: extraOpts.ResourceSuffix,
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructCreateVolumeRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.VOLUME,
		ResourceLocation: extraOpts.ResourceLocation,
		RequestSuffix: extraOpts.ResourceSuffix,
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}

func constructCreateSnapshotRequestOpts(opts entity.CreateUpdateOptions, extraOpts *ExtraOption) RequestOption {
	return RequestOption{
		Action: CREATE,
		Resource: consts.SNAPSHOT,
		ResourceLocation: extraOpts.ResourceLocation,
		RequestSuffix: extraOpts.ResourceSuffix,
		Body: opts,
		Headers: map[string]string{consts.AuthToken: defaultController.Token()},
	}
}
