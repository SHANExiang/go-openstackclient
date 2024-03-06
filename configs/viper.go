package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"runtime"
)

func Viper() *viper.Viper {
	v := viper.New()
	currentPath := getCurrentAbPathByCaller()
	v.SetConfigFile(currentPath + "\\openstack.yaml")

	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fail to read yaml err:%v \n", err))
	}
	if err = v.Unmarshal(&CONF); err != nil {
		fmt.Println(err)
	}
	return v
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}