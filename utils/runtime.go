package utils

import (
	"runtime"
	"strings"
)


//GetCurrentMethodName get method name
func GetCurrentMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return strings.Split(funcName, ".")[1]
}
