package entity


type CreateUpdateOptions interface {
	ToRequestBody()     string
}
