package entity

import (
	"fmt"
)

type Identity struct {
	Methods  []string `json:"methods"`
	Password          `json:"password"`
}

type Password struct {
    Userr       `json:"user"`
}

type Domain struct {
	Name string `json:"name"`
}

type Userr struct {
	Domain           `json:"domain"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
}

type Auth struct {
     Identity           `json:"identity"`
     Scope              `json:"scope"`
}

type Projectt struct {
	Name string       `json:"name"`
	Domain            `json:"domain"`
}

type Scope struct {
	Projectt           `json:"project"`
}

type AuthOption struct {
	Auth              `json:"auth"`
}


func (opts *AuthOption) ToRequestBody() string {
	reqBody, err := BuildRequestBody(opts, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to build request body %s", err))
	}
	return reqBody
}
