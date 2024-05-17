package main

import (
	"go-openstackclient/configs"
	"go-openstackclient/internal/service"
)

func init() {
	configs.Viper()
}

func main() {
	//service.CreateNetworkHelper()
	//service.CreateInstanceHelper("599771ab-5682-49f5-a291-cf674aad91fb", "sdn_test")


	service.ListNetworkHelper()
}

