package entity

type SyncResultsMap struct {
	Count       int `json:"count"`
	SyncResults []struct {
		Status     string `json:"status"`
		Oper       string `json:"oper"`
		UpdateTime string `json:"update_time"`
		ResType    string `json:"res_type"`
		TenantName string `json:"tenant_name"`
		ResName    string `json:"res_name"`
		CreateTime string `json:"create_time"`
		Id         string `json:"id"`
	} `json:"sync_results"`
}

type SyncSummaryResult struct {
	Count       int `json:"count"`
	SyncResults []struct {
		TotalSyncResources   string  `json:"total_sync_resources"`
		SyncEndTime          string  `json:"sync_end_time"`
		TotalCreateResources int     `json:"total_create_resources"`
		TotalTimeUsed        float64 `json:"total_time_used"`
		TotalUpdateResources int     `json:"total_update_resources"`
		SyncStartTime        string  `json:"sync_start_time"`
		SyncId               int     `json:"sync_id"`
		Id                   string  `json:"id"`
		TotalDeleteResources int     `json:"total_delete_resources"`
	} `json:"sync_results"`
}
