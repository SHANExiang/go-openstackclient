package entity

import (
	"encoding/json"
	"fmt"
	"go-openstackclient/consts"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var t time.Time

func isZero(v reflect.Value) bool {
	//fmt.Printf("\n\nchecking isZero for value: %+v\n", v)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return false
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(t) {
			if v.Interface().(time.Time).IsZero() {
				return true
			}
			return false
		}
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	//fmt.Printf("zero type for value: %+v\n\n\n", z)
	return v.Interface() == z.Interface()
}

func ProtocolAnyToNull(optsMap map[string]interface{}, parent string) map[string]interface{} {
	if val, ok := optsMap["protocol"]; ok && val.(string) == "any" && parent == consts.FIREWALLRULE {
		optsMap["protocol"] = nil
	}
	return optsMap
}

func BuildRequestBody(opts interface{}, parent string) (string, error) {
	optsValue := reflect.ValueOf(opts)
	if optsValue.Kind() == reflect.Ptr {
		optsValue = optsValue.Elem()
	}

	optsType := reflect.TypeOf(opts)
	if optsType.Kind() == reflect.Ptr {
		optsType = optsType.Elem()
	}

	optsMap := make(map[string]interface{})
	if optsValue.Kind() == reflect.Struct {
		//fmt.Printf("optsValue.Kind() is a reflect.Struct: %+v\n", optsValue.Kind())
		for i := 0; i < optsValue.NumField(); i++ {
			v := optsValue.Field(i)
			f := optsType.Field(i)

			if f.Name != strings.Title(f.Name) {
				//fmt.Printf("Skipping field: %s...\n", f.Name)
				continue
			}

			//fmt.Printf("Starting on field: %s...\n", f.Name)

			zero := isZero(v)
			//fmt.Printf("v is zero?: %v\n", zero)

			// if the field has a required tag that's set to "true"
			if requiredTag := f.Tag.Get("required"); requiredTag == "true" {
				//fmt.Printf("Checking required field [%s]:\n\tv: %+v\n\tisZero:%v\n", f.Name, v.Interface(), zero)
				// if the field's value is zero, return a missing-argument error
				if zero {
					// if the field has a 'required' tag, it can't have a zero-value
					return "nil", fmt.Errorf("missing required argument %s", f.Name)
				}
			}

			jsonTag := f.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			if v.Kind() == reflect.Slice || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice) {
				sliceValue := v
				if sliceValue.Kind() == reflect.Ptr {
					sliceValue = sliceValue.Elem()
				}

				for i := 0; i < sliceValue.Len(); i++ {
					element := sliceValue.Index(i)
					if element.Kind() == reflect.Struct || (element.Kind() == reflect.Ptr && element.Elem().Kind() == reflect.Struct) {
						_, err := BuildRequestBody(element.Interface(), "")
						if err != nil {
							return "nil", err
						}
					}
				}
			}
			if v.Kind() == reflect.Struct || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
				if zero {
					//fmt.Printf("value before change: %+v\n", optsValue.Field(i))
					if jsonTag != "" {
						jsonTagPieces := strings.Split(jsonTag, ",")
						if len(jsonTagPieces) > 1 && jsonTagPieces[1] == "omitempty" {
							if v.CanSet() {
								if !v.IsNil() {
									if v.Kind() == reflect.Ptr {
										v.Set(reflect.Zero(v.Type()))
									}
								}
								//fmt.Printf("value after change: %+v\n", optsValue.Field(i))
							}
						}
					}
					continue
				}

				//fmt.Printf("Calling BuildRequestBody with:\n\tv: %+v\n\tf.Name:%s\n", v.Interface(), f.Name)
				_, err := BuildRequestBody(v.Interface(), f.Name)
				if err != nil {
					return "nil", err
				}
			}
		}

		//fmt.Printf("opts: %+v \n", opts)

		b, err := json.Marshal(opts)
		if err != nil {
			return "nil", err
		}

		//fmt.Printf("string(b): %s\n", string(b))

		err = json.Unmarshal(b, &optsMap)
		if err != nil {
			return "nil", err
		}
		optsMap = ProtocolAnyToNull(optsMap, parent)
		//fmt.Printf("optsMap: %+v\n", optsMap)

		if parent != "" {
			optsMap = map[string]interface{}{parent: optsMap}
		}
		//fmt.Printf("optsMap after parent added: %+v\n", optsMap)
		b, err = json.Marshal(optsMap)
		if err != nil {
			return "nil", err
		}
		return string(b), nil
	}
	// Return an error if the underlying type of 'opts' isn't a struct.
	return "nil", fmt.Errorf("options type is not a struct")
}

func InterfaceSliceToStringSlice(value []interface{}) []string {
	stringSlice := make([]string, len(value))
	for i, v := range value {
		if str, ok := v.(string); ok {
			stringSlice[i] = str
		} else if in, ok := v.(int); ok {
			stringSlice[i] = strconv.Itoa(in)
		}else {
			// 处理类型不匹配的情况
			log.Printf("Element at index %d is not a string\n", i)
		}
	}
	return stringSlice
}

func IsUUID(uuid string) bool {
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	match, _ := regexp.MatchString(uuidRegex, uuid)
	return match
}

