package entity

type Snapshot struct {
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
	Id          string `json:"id"`
	Metadata    struct {
		Key string `json:"key"`
	} `json:"metadata"`
	Name            string      `json:"name"`
	Size            int         `json:"size"`
	Status          string      `json:"status"`
	UpdatedAt       interface{} `json:"updated_at"`
	VolumeId        string      `json:"volume_id"`
	GroupSnapshotId interface{} `json:"group_snapshot_id"`
	UserId          string      `json:"user_id"`
	ConsumesQuota   bool        `json:"consumes_quota"`
}

type SnapshotMap struct {
    Snapshot `json:"snapshot"`
}

type Snapshots struct {
	Ss                  []Snapshot `json:"snapshots"`
}
