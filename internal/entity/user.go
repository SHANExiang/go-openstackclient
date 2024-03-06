package entity

type User struct {
	DefaultProjectId string `json:"default_project_id"`
	Description      string `json:"description"`
	DomainId         string `json:"domain_id"`
	Enabled          bool   `json:"enabled"`
	Id               string `json:"id"`
	Links            struct {
		Self string `json:"self"`
	} `json:"links"`
	Name    string `json:"name"`
	Options struct {
		IgnorePasswordExpiry bool `json:"ignore_password_expiry"`
	} `json:"options"`
	PasswordExpiresAt string `json:"password_expires_at"`
}

type UserMap struct {
	User `json:"user"`
}
