package entity


type Attachment struct {
	ServerId     string      `json:"server_id"`
	AttachmentId string      `json:"attachment_id"`
	AttachedAt   string      `json:"attached_at"`
	HostName     interface{} `json:"host_name"`
	VolumeId     string      `json:"volume_id"`
	Device       string      `json:"device"`
	Id           string      `json:"id"`
}

type Volume struct {
	Attachments        []Attachment  `json:"attachments"`
	AvailabilityZone   string        `json:"availability_zone"`
	Bootable           string        `json:"bootable"`
	ConsistencygroupId interface{}   `json:"consistencygroup_id"`
	CreatedAt          string        `json:"created_at"`
	Description        interface{}   `json:"description"`
	Encrypted          bool          `json:"encrypted"`
	Id                 string        `json:"id"`
	Links              []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Metadata struct {
	} `json:"metadata"`
	MigrationStatus   interface{} `json:"migration_status"`
	Multiattach       bool        `json:"multiattach"`
	Name              interface{} `json:"name"`
	ReplicationStatus interface{} `json:"replication_status"`
	Size              int         `json:"size"`
	SnapshotId        interface{} `json:"snapshot_id"`
	SourceVolid       interface{} `json:"source_volid"`
	Status            string      `json:"status"`
	UpdatedAt         interface{} `json:"updated_at"`
	UserId            string      `json:"user_id"`
	VolumeType        string      `json:"volume_type"`
	GroupId           interface{} `json:"group_id"`
	ProviderId        interface{} `json:"provider_id"`
	ServiceUuid       interface{} `json:"service_uuid"`
	SharedTargets     bool        `json:"shared_targets"`
	ClusterName       interface{} `json:"cluster_name"`
	VolumeTypeId      string      `json:"volume_type_id"`
	ConsumesQuota     bool        `json:"consumes_quota"`
}

type VolumeMap struct {
	 Volume `json:"volume"`
}

type Volumes struct {
	Vs            []Volume `json:"volumes"`
}

type VolumeToImage struct {
	OsVolumeUploadImage struct {
		ContainerFormat    string      `json:"container_format"`
		DiskFormat         string      `json:"disk_format"`
		DisplayDescription interface{} `json:"display_description"`
		Id                 string      `json:"id"`
		ImageId            string      `json:"image_id"`
		ImageName          string      `json:"image_name"`
		Protected          bool        `json:"protected"`
		Size               int         `json:"size"`
		Status             string      `json:"status"`
		UpdatedAt          string      `json:"updated_at"`
		Visibility         string      `json:"visibility"`
		VolumeType         interface{} `json:"volume_type"`
	} `json:"os-volume_upload_image"`
}
