package entity

type VolumeType struct {
	Description string `json:"description"`
	ExtraSpecs  struct {
		Capabilities         string    `json:"capabilities,omitempty"`
		VolumeBackendName    string    `json:"volume_backend_name"`
	} `json:"extra_specs"`
	Id                         string      `json:"id"`
	IsPublic                   bool        `json:"is_public"`
	Name                       string      `json:"name"`
	OsVolumeTypeAccessIsPublic bool        `json:"os-volume-type-access:is_public"`
	QosSpecsId                 interface{} `json:"qos_specs_id"`
}

type VolumeTypes struct {
	VTs    []VolumeType       `json:"volume_types"`
}
