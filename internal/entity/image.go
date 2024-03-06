package entity

import "time"

type ImageMap struct {
	Status          string        `json:"status"`
	Name            string        `json:"name"`
	Tags            []interface{} `json:"tags"`
	ContainerFormat string        `json:"container_format"`
	CreatedAt       time.Time     `json:"created_at"`
	Size            interface{}   `json:"size"`
	DiskFormat      string        `json:"disk_format"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Visibility      string        `json:"visibility"`
	Locations       []interface{} `json:"locations"`
	Self            string        `json:"self"`
	MinDisk         int           `json:"min_disk"`
	Protected       bool          `json:"protected"`
	Id              string        `json:"id"`
	File            string        `json:"file"`
	Checksum        interface{}   `json:"checksum"`
	OsHashAlgo      interface{}   `json:"os_hash_algo"`
	OsHashValue     interface{}   `json:"os_hash_value"`
	OsHidden        bool          `json:"os_hidden"`
	Owner           string        `json:"owner"`
	VirtualSize     interface{}   `json:"virtual_size"`
	MinRam          int           `json:"min_ram"`
	Schema          string        `json:"schema"`
}

type Images struct {
	Is []struct {
		Status          string        `json:"status"`
		Name            string        `json:"name"`
		Tags            []interface{} `json:"tags"`
		ContainerFormat string        `json:"container_format"`
		CreatedAt       time.Time     `json:"created_at"`
		DiskFormat      string        `json:"disk_format"`
		UpdatedAt       time.Time     `json:"updated_at"`
		Visibility      string        `json:"visibility"`
		Self            string        `json:"self"`
		MinDisk         int           `json:"min_disk"`
		Protected       bool          `json:"protected"`
		Id              string        `json:"id"`
		File            string        `json:"file"`
		Checksum        string        `json:"checksum"`
		OsHashAlgo      string        `json:"os_hash_algo"`
		OsHashValue     string        `json:"os_hash_value"`
		OsHidden        bool          `json:"os_hidden"`
		Owner           string        `json:"owner"`
		Size            int           `json:"size"`
		MinRam          int           `json:"min_ram"`
		Schema          string        `json:"schema"`
		VirtualSize     interface{}   `json:"virtual_size"`
	} `json:"images"`
	Schema string `json:"schema"`
	First  string `json:"first"`
}

type ImageMember struct {
	CreatedAt time.Time `json:"created_at"`
	ImageId   string    `json:"image_id"`
	MemberId  string    `json:"member_id"`
	Schema    string    `json:"schema"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}
