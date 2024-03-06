package service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
	"go-openstackclient/internal"
	"go-openstackclient/internal/entity"
	"go-openstackclient/utils"
	"log"
	"sync"
	"time"
)

var supportedOctaviaResourceTypes = [...]string{
	consts.LOADBALANCER, consts.LISTENER, consts.POOL, consts.MEMBER,
	consts.HEALTHMONITOR, consts.L7POLICY, consts.L7RULE,
}

func initOctaviaOutputChannels() map[string]chan internal.Output {
	outputChannel := make(map[string]chan internal.Output)
	for _, resourceType := range supportedOctaviaResourceTypes {
		outputChannel[resourceType] = make(chan internal.Output, 0)
	}
	return outputChannel
}

type Octavia struct {
	internal.Request
	headers            map[string]string
	tag                string
	snowflake          *utils.Snowflake
	ExternalNetwork    string
	DeleteChannels     map[string]chan internal.Output
	mu                 sync.Mutex
}

func NewLB(token string, client *fasthttp.Client) *Octavia {
	return &Octavia{
		Request: internal.Request{
			UrlPrefix: fmt.Sprintf("http://%s:%d/v2.0/", configs.CONF.Host, consts.OctaviaPort),
			Client: client,
		},
		headers:         map[string]string{"X-Auth-Token": token},
		tag:             configs.CONF.Host + "_",
		snowflake:       utils.NewSnowflake(uint16(1)),
		ExternalNetwork: configs.CONF.ExternalNetwork,
		DeleteChannels:  initOctaviaOutputChannels(),
	}
}

func (o *Octavia) MakeDeleteChannel(resourceType string, length int) chan internal.Output {
	defer o.mu.Unlock()
	o.mu.Lock()

	o.DeleteChannels[resourceType] = make(chan internal.Output, length)
	return o.DeleteChannels[resourceType]
}

// load balancer

func (o *Octavia) CreateLoadbalancer(opts entity.CreateLoadbalancerOpts) string {
	opts.Name = fmt.Sprintf("%s_%s", opts.Name, consts.LOADBALANCER)
	urlSuffix := "lbaas/loadbalancers"
	reqBody := opts.ToRequestBody()
	resp := o.Post(o.headers, urlSuffix, reqBody)
	var lb entity.LoadbalancerMap
	_ = json.Unmarshal(resp, &lb)

	//cache.RedisClient.SetMap(l.tag + consts.LOADBALANCERS, lb.Loadbalancer.Id, lb)

	o.makeSureLbActive(lb.Loadbalancer.Id)
	log.Println("==============create loadbalancer success", lb.Loadbalancer.Id)
	return lb.Loadbalancer.Id
}

func (o *Octavia) deleteLoadbalancer(ipId string) internal.Output {
	outputObj := internal.Output{ParametersMap: map[string]string{"loadbalancer_id": ipId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
			outputObj.Response = err
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/loadbalancers/%s", ipId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
    o.makeSureLbDeleted(ipId)
	return outputObj
}

func (l *Octavia) getLoadbalancer(ipId string) entity.LoadbalancerMap {
	urlSuffix := fmt.Sprintf("lbaas/loadbalancers/%s", ipId)
	resp := l.Get(l.headers, urlSuffix)
	var lb entity.LoadbalancerMap
	_ = json.Unmarshal(resp, &lb)
	return lb
}

func (l *Octavia) ListLoadbalancers() entity.Loadbalancers {
	urlSuffix := "lbaas/loadbalancers"
	resp := l.List(l.headers, urlSuffix)
	var lbs entity.Loadbalancers
	_ = json.Unmarshal(resp, &lbs)
	log.Println("==============List loadbalancers success, there had", len(lbs.LBs))
	return lbs
}

func (o *Octavia) DeleteLoadbalancers() {
	lbs := o.ListLoadbalancers()
	ch := o.MakeDeleteChannel(consts.LOADBALANCER, len(lbs.LBs))

	for _, lb := range lbs.LBs {
		tempLb := lb
		go func() {
			ch <- o.deleteLoadbalancer(tempLb.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Loadbalancers were deleted completely")
}

// listener

func (o *Octavia) CreateListener(opts entity.CreateListenerOpts) string {
	urlSuffix := "lbaas/listeners"
	opts.Name = fmt.Sprintf("%s_%s", opts.Name, consts.LISTENER)
	reqBody := opts.ToRequestBody()
	resp := o.Post(o.headers, urlSuffix, reqBody)
	var listener entity.ListenerMap
	_ = json.Unmarshal(resp, &listener)

	//cache.RedisClient.SetMap(l.tag + consts.LISTENERS, listener.Listener.Id, listener)

	o.makeSureLbActive(opts.LoadbalancerID)
	log.Println("==============Create listener success", listener.Listener.Id)
	return listener.Listener.Id
}

func (o *Octavia) deleteListener(listenerId string) internal.Output {
	outputObj := internal.Output{ParametersMap: map[string]string{"listener_id": listenerId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/listeners/%s", listenerId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	return outputObj
}

func (o *Octavia) getListener(ipId string) entity.ListenerMap {
	urlSuffix := fmt.Sprintf("lbaas/listeners/%s", ipId)
	resp := o.Get(o.headers, urlSuffix)
	var listener entity.ListenerMap
	_ = json.Unmarshal(resp, &listener)
	return listener
}

func (o *Octavia) ListListeners() entity.Listeners {
	urlSuffix := "lbaas/listeners"
	resp := o.List(o.headers, urlSuffix)
	var listeners entity.Listeners
	_ = json.Unmarshal(resp, &listeners)
	log.Println("==============List listeners success, there had", len(listeners.Liss))
	return listeners
}

func (o *Octavia) DeleteListeners() {
	listeners := o.ListListeners()
	ch := o.MakeDeleteChannel(consts.LISTENER, len(listeners.Liss))

	for _, listener := range listeners.Liss {
		temp := listener
		go func() {
			ch <- o.deleteListener(temp.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Listeners were deleted completely")
}

func (o *Octavia) makeSureLbActive(lbId string) entity.LoadbalancerMap {
	lb := o.getLoadbalancer(lbId)

    for lb.Loadbalancer.ProvisioningStatus != consts.ACTIVE {
    	time.Sleep(5 * time.Second)
    	lb = o.getLoadbalancer(lbId)
	}
	return lb
}

func (o *Octavia) makeSureLbDeleted(lbId string) {
	lb := o.getLoadbalancer(lbId)
	for &lb != nil {
		time.Sleep(5 * time.Second)
		lb = o.getLoadbalancer(lbId)
	}
	log.Println("*******************Lb was deleted success")
}

// pool

func (o *Octavia) CreatePool(opts entity.CreatePoolOpts) string {
	urlSuffix := "lbaas/pools"
	opts.Name = fmt.Sprintf("%s_%s", opts.Name, consts.POOL)
	reqBody := opts.ToRequestBody()
	resp := o.Post(o.headers, urlSuffix, reqBody)
	var pool entity.PoolMap
	_ = json.Unmarshal(resp, &pool)

	o.makeSurePoolActive(pool.Pool.Id)
	for _, lb := range pool.Pool.Loadbalancers {
		o.makeSureLbActive(lb.Id)
	}
	log.Println("==============Create pool success", pool.Pool.Id)
	return pool.Pool.Id
}

func (o *Octavia) deletePool(poolId string) internal.Output {
	pool := o.getPool(poolId)
	outputObj := internal.Output{ParametersMap: map[string]string{"pool_id": poolId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf( "lbaas/pools/%s", poolId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	for _, lb := range pool.Loadbalancers {
		o.makeSureLbActive(lb.Id)
	}
	return outputObj
}

func (l *Octavia) getPool(ipId string) entity.PoolMap {
	urlSuffix := fmt.Sprintf("lbaas/pools/%s", ipId)
	resp := l.Get(l.headers, urlSuffix)
	var pool entity.PoolMap
	_ = json.Unmarshal(resp, &pool)
	return pool
}

func (l *Octavia) ListPools() entity.Pools {
	urlSuffix := "lbaas/pools"
	resp := l.List(l.headers, urlSuffix)
	var pools entity.Pools
	_ = json.Unmarshal(resp, &pools)
	log.Println("==============List pools success, there had", len(pools.Ps))
	return pools
}

func (o *Octavia) DeletePools() {
	pools := o.ListPools()
	ch := o.MakeDeleteChannel(consts.POOL, len(pools.Ps))
	for _, pool := range pools.Ps {
		//for _, member := range pool.Members {
		//	o.deletePoolMember(pool, member.(map[string]interface{})["id"].(string))
		//}
		temp := pool
		go func() {
			ch <- o.deletePool(temp.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Pool were deleted completely")
}

func (l *Octavia) makeSurePoolActive(poolId string) entity.PoolMap {
	pool := l.getPool(poolId)

	for pool.Pool.ProvisioningStatus != consts.ACTIVE {
		time.Sleep(5 * time.Second)
		pool = l.getPool(poolId)
	}
	return pool
}

// pool member

func (o *Octavia) CreatePoolMember(poolId string, opts entity.CreateMemberOpts) string {
	urlSuffix := fmt.Sprintf("lbaas/pools/%s/members", poolId)
	opts.Name = fmt.Sprintf("%s_%s", opts.Name, consts.POOL)
	reqBody := opts.ToRequestBody()
	resp := o.Post(o.headers, urlSuffix, reqBody)
	var member entity.MemberMap
	_ = json.Unmarshal(resp, &member)

	o.makeSurePoolActive(poolId)
	log.Println("==============Create member success", member.Member.Id)
	return member.Member.Id
}

func (o *Octavia) deletePoolMember(pool entity.Pool, memberId string) internal.Output {
	defer o.mu.Unlock()
	o.mu.Lock()
	outputObj := internal.Output{ParametersMap: map[string]string{"member_id": memberId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/pools/%s/members/%s", pool.Id, memberId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	o.makeSurePoolActive(pool.Id)
	for _, lb := range pool.Loadbalancers {
		o.makeSureLbActive(lb.Id)
	}
	return outputObj
}

func (l *Octavia) getPoolMember(poolId, memberId string) entity.MemberMap {
	urlSuffix := fmt.Sprintf("lbaas/pools/%s/members/%s", poolId, memberId)
	resp := l.Get(l.headers, urlSuffix)
	var member entity.MemberMap
	_ = json.Unmarshal(resp, &member)
	return member
}

func (o *Octavia) ListPoolMembers(poolId string) entity.Members {
	urlSuffix := fmt.Sprintf("lbaas/pools/%s/members", poolId)
	resp := o.List(o.headers, urlSuffix)
	var ms entity.Members
	_ = json.Unmarshal(resp, &ms)
	log.Println("==============List pool members success, there had", len(ms.Ms))
	return ms
}


func (o *Octavia) DeleteMembers() {
	pools := o.ListPools()
	var memberNumber int
	for _, pool := range pools.Ps {
		memberNumber += len(pool.Members)
	}
	ch := o.MakeDeleteChannel(consts.MEMBER, memberNumber)
	for _, pool := range pools.Ps {
		tempPool := pool
		for _, member := range pool.Members {
			temp := member
			go func() {
				ch <- o.deletePoolMember(tempPool, temp.Id)
			}()
		}
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("Pool members were deleted completely")
}

// health monitor

func (l *Octavia) CreateHealthMonitor(opts entity.CreateHealthMonitorOpts) string {
	urlSuffix := "lbaas/healthmonitors"
	opts.Name = fmt.Sprintf("%s_%s", opts.Name, consts.HEALTHMONITOR)
	reqBody := opts.ToRequestBody()
	resp := l.Post(l.headers, urlSuffix, reqBody)
	var healthMonitor entity.HealthMonitorMap
	_ = json.Unmarshal(resp, &healthMonitor)

	log.Println("==============Create health monitor success", healthMonitor.Healthmonitor.Id)
	return healthMonitor.Healthmonitor.Id
}

func (o *Octavia) deleteHealthMonitor(healthmonitorId string) internal.Output {
	outputObj := internal.Output{ParametersMap: map[string]string{"health_monitor_id": healthmonitorId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/healthmonitors/%s", healthmonitorId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	return outputObj
}

func (l *Octavia) getHealthMonitor(healthmonitorId string) entity.HealthMonitorMap {
	urlSuffix := fmt.Sprintf("lbaas/healthmonitors/%s", healthmonitorId)
	resp := l.Get(l.headers, urlSuffix)
	var healthmonitor entity.HealthMonitorMap
	_ = json.Unmarshal(resp, &healthmonitor)
	return healthmonitor
}

func (l *Octavia) ListHealthMonitors() entity.HealthMonitors {
	urlSuffix := "lbaas/healthmonitors"
	resp := l.List(l.headers, urlSuffix)
	var hms entity.HealthMonitors
	_ = json.Unmarshal(resp, &hms)
	log.Println("==============List health monitors success, there had", len(hms.HMs))
	return hms
}

func (o *Octavia) DeleteHealthmonitors() {
	healthmonitors := o.ListHealthMonitors()
	ch := o.MakeDeleteChannel(consts.HEALTHMONITOR, len(healthmonitors.HMs))
	for _, healthmonitor := range healthmonitors.HMs {
		temp := healthmonitor
		go func() {
			ch <- o.deleteHealthMonitor(temp.Id)
		}()
	}
	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	for _, healthmonitor := range healthmonitors.HMs {
		for _, pool := range healthmonitor.Pools {
			o.makeSurePoolActive(pool.Id)
		}
	}
	log.Println("Health monitors were deleted completely")
}

// L7 policy

func (l *Octavia) CreateL7Policy(listenerId string) string {
	urlSuffix := "lbaas/l7policies"
	reqBody := fmt.Sprintf("{\"l7policy\": {\"description\": \"Redirect requests to example.com\", \"admin_state_up\": true, \"listener_id\": \"%+v\", \"redirect_url\": \"http://www.example.com\", \"redirect_http_code\": 301, \"name\": \"redirect-example.com\", \"action\": \"REDIRECT_TO_URL\", \"position\": 1, \"tags\": [\"test_tag\"]}}", listenerId)
	resp := l.Post(l.headers, urlSuffix, reqBody)
	var l7policy entity.L7PolicyMap
	_ = json.Unmarshal(resp, &l7policy)

	//cache.RedisClient.SetMap(l.tag + consts.L7POLICIES, l7policy.L7Policy.Id, l7policy)
	l.makeSureL7PolicyActive(l7policy.L7Policy.Id)
	log.Println("==============create l7policy success", l7policy.L7Policy.Id)
	return l7policy.L7Policy.Id
}

func (o *Octavia) deleteL7Policy(l7policyId string) internal.Output {
	outputObj := internal.Output{ParametersMap: map[string]string{"l7Policy_id": l7policyId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/l7policies/%s", l7policyId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	return outputObj
}

func (l *Octavia) getL7Policy(l7policyId string) entity.L7PolicyMap {
	urlSuffix := fmt.Sprintf("lbaas/l7policies/%s", l7policyId)
	resp := l.Get(l.headers, urlSuffix)
	var l7policy entity.L7PolicyMap
	_ = json.Unmarshal(resp, &l7policy)
	return l7policy
}

func (l *Octavia) ListL7Policies() entity.L7Policies {
	urlSuffix := "lbaas/l7policies"
	resp := l.List(l.headers, urlSuffix)
	var l7ps entity.L7Policies
	_ = json.Unmarshal(resp, &l7ps)
	log.Println("==============List l7 policies success, there had", len(l7ps.L7Ps))
	return l7ps
}

func (l *Octavia) makeSureL7PolicyActive(l7PolicyId string) entity.L7PolicyMap {
	l7Policy := l.getL7Policy(l7PolicyId)

	for l7Policy.L7Policy.ProvisioningStatus != consts.ACTIVE {
		time.Sleep(5 * time.Second)
		l7Policy = l.getL7Policy(l7PolicyId)
	}
	//l.SyncMap(l.tag + consts.L7POLICIES, l7PolicyId, nil, l7Policy)
	return l7Policy
}

func (o *Octavia) DeleteL7Policies() {
	l7policies := o.ListL7Policies()
	ch := o.MakeDeleteChannel(consts.L7POLICY, len(l7policies.L7Ps))
	for _, l7policy := range l7policies.L7Ps {
		temp := l7policy
		go func() {
			ch <- o.deleteL7Policy(temp.Id)
		}()
	}

	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {}
	}
	log.Println("L7 policies were deleted completely")
}

// L7 rule

func (l *Octavia) CreateL7Rule(policyId string) string {
	urlSuffix := fmt.Sprintf("lbaas/l7policies/%s/rules", policyId)
	reqBody := fmt.Sprintf("{\"rule\": {\"compare_type\": \"REGEX\", \"invert\": false, \"type\": \"PATH\", \"value\": \"/images*\", \"admin_state_up\": true, \"tags\": [\"test_tag\"]}}")
	resp := l.Post(l.headers, urlSuffix, reqBody)
	var l7Rule entity.L7RuleMap
	_ = json.Unmarshal(resp, &l7Rule)

	//cache.RedisClient.SetMap(l.tag + consts.L7RULES, l7Rule.Rule.Id, l7Rule)
	log.Println("==============create l7Rule success", l7Rule.Rule.Id)
	return l7Rule.Rule.Id
}

func (o *Octavia) deleteL7Rule(l7PolicyId, l7RuleId string) internal.Output {
	outputObj := internal.Output{ParametersMap: map[string]string{"l7Rule_id": l7RuleId}}
	defer func() {
		if err := recover(); err != nil {
			log.Println("catch error：", err)
			outputObj.Success = false
		}
	}()

	urlSuffix := fmt.Sprintf("lbaas/l7policies/%s/rules/%s", l7PolicyId, l7RuleId)
	outputObj.Success, outputObj.Response = o.Delete(o.headers, urlSuffix)
	return outputObj
}

func (o *Octavia) getL7Rule(l7PolicyId, l7RuleId string) entity.L7RuleMap {
	urlSuffix := fmt.Sprintf("lbaas/l7policies/%s/rules/%s", l7PolicyId, l7RuleId)
	resp := o.Get(o.headers, urlSuffix)
	var l7Rule entity.L7RuleMap
	_ = json.Unmarshal(resp, &l7Rule)
	return l7Rule
}

func (o *Octavia) DeleteL7Rules() {
	l7policies := o.ListL7Policies()
	var ruleNumber int
	for _, policy := range l7policies.L7Ps {
        ruleNumber += len(policy.Rules)
	}
	ch := o.MakeDeleteChannel(consts.L7RULE, ruleNumber)
	for _, l7policy := range l7policies.L7Ps {
		tempPolicy := l7policy
		for _, rule := range l7policy.Rules {
			temp := rule
			go func() {
				ch <- o.deleteL7Rule(tempPolicy.Id, temp.Id)
			}()
		}
	}

	if len(ch) != cap(ch) {
		for len(ch) != cap(ch) {
		}
	}
	log.Println("L7 rules were deleted completely")
}