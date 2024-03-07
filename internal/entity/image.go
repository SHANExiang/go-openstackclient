package entity

import (
	"fmt"
	"go-openstackclient/consts"
	"time"
)

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
type ImageVisibility string

const (
	// ImageVisibilityPublic all users
	ImageVisibilityPublic ImageVisibility = "public"

	// ImageVisibilityPrivate users with tenantId == tenantId(owner)
	ImageVisibilityPrivate ImageVisibility = "private"

	// ImageVisibilityShared images are visible to:
	// - users with tenantId == tenantId(owner)
	// - users with tenantId in the member-list of the image
	// - users with tenantId in the member-list with member_status == 'accepted'
	ImageVisibilityShared ImageVisibility = "shared"

	// ImageVisibilityCommunity images:
	// - all users can see and boot it
	// - users with tenantId in the member-list of the image with
	//	 member_status == 'accepted' have this image in their default image-list.
	ImageVisibilityCommunity ImageVisibility = "community"
)

// CreateImageOpts represents options used to create an image.
type CreateImageOpts struct {
	// Name is the name of the new image.
	Name string `json:"name" required:"true"`

	// Id is the the image ID.
	ID string `json:"id,omitempty"`

	// Visibility defines who can see/use the image.
	Visibility *ImageVisibility `json:"visibility,omitempty"`

	// Hidden is whether the image is listed in default image list or not.
	Hidden *bool `json:"os_hidden,omitempty"`

	// Tags is a set of image tags.
	Tags []string `json:"tags,omitempty"`

	// ContainerFormat is the format of the
	// container. Valid values are ami, ari, aki, bare, and ovf.
	ContainerFormat string `json:"container_format,omitempty"`

	// DiskFormat is the format of the disk. If set,
	// valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi,
	// and iso.
	DiskFormat string `json:"disk_format,omitempty"`

	// MinDisk is the amount of disk space in
	// GB that is required to boot the image.
	MinDisk int `json:"min_disk,omitempty"`

	// MinRAM is the amount of RAM in MB that
	// is required to boot the image.
	MinRAM int `json:"min_ram,omitempty"`

	// protected is whether the image is not deletable.
	Protected *bool `json:"protected,omitempty"`

	// properties is a set of properties, if any, that
	// are associated with the image.
	Properties map[string]string `json:"-"`
}

func (opts *CreateImageOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.Image)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}
