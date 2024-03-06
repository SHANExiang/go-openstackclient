package entity

import (
	"fmt"
	"go-openstackclient/consts"
)

type Project struct {
	IsDomain    bool   `json:"is_domain"`
	Description string `json:"description"`
	Links       struct {
		Self string `json:"self"`
	} `json:"links"`
	Tags     []interface{} `json:"tags"`
	Enabled  bool          `json:"enabled"`
	Id       string        `json:"id"`
	ParentId string        `json:"parent_id"`
	Options  struct {
	} `json:"options"`
	DomainId string `json:"domain_id"`
	Name     string `json:"name"`
}

type ProjectMap struct {
	Project    `json:"project"`
}

type Projects struct {
	Ps        []Project    `json:"projects"`
}

// CreateProjectOpts represents parameters used to create a project.
type CreateProjectOpts struct {
	// DomainID is the ID this project will belong under.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the project status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// IsDomain indicates if this project is a domain.
	IsDomain *bool `json:"is_domain,omitempty"`

	// Name is the name of the project.
	Name string `json:"name" required:"true"`

	// ParentID specifies the parent project of this new project.
	ParentID string `json:"parent_id,omitempty"`

	// Description is the description of the project.
	Description string `json:"description,omitempty"`

	// Tags is a list of tags to associate with the project.
	Tags []string `json:"tags,omitempty"`

	// Extra is free-form extra key/value pairs to describe the project.
	Extra map[string]interface{} `json:"-"`
}

func (opts *CreateProjectOpts) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, consts.PROJECT)
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}
