package service

var (
	ResourcesMap      = make(map[string]Resource)
)

type Resource struct {
    Name               string
    Type               string
    PropsObj           interface{}
    Dependencies       map[string]string
    Done               bool
}
