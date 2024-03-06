package entity

type RbacPolicy struct {
	TargetTenant string `json:"target_tenant"`
	TenantId     string `json:"tenant_id"`
	ObjectType   string `json:"object_type"`
	ObjectId     string `json:"object_id"`
	Action       string `json:"action"`
	ProjectId    string `json:"project_id"`
	Id           string `json:"id"`
}

type RbacPolicyMap struct {
	RbacPolicy `json:"rbac_policy"`
}
