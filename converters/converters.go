package converters

// copy from helpers

import (
	"errors"
	"fmt"
	"strconv"
)

func StringArrayToInterfaceArray(strings []string) []interface{} {
	var intfs []interface{}
	for _, s := range strings {
		intfs = append(intfs, s)
	}
	return intfs
}

func InterfaceToBool(i interface{}) (bool, error) {
	res := false

	switch i.(type) {
	case bool:
		res = i.(bool)
	default:
		return res, errors.New(fmt.Sprint(i) + " is not a boolean.")
	}

	return res, nil
}

func InterfaceToInt64(i interface{}) (int64, error) {
	if res, err := strconv.ParseInt(fmt.Sprint(i), 10, 64); err != nil {
		return res, err
	} else {
		return res, nil
	}
}

func BackQuote(object string) string {
	return "`" + object + "`"
}

func CreateStatusMap(key string, value interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m[key] = value
	return m
}

func MapInterfaceToMapString(interfaceMap map[string]interface{}) map[string]string {
	stringMap := make(map[string]string)
	for key, val := range interfaceMap {
		stringMap[key] = fmt.Sprint(val)
	}
	return stringMap
}
