package utils

import (
	"reflect"
	"strconv"
)

func Join(stringSlice []string, delimiter string, enclosure string) string {
	resultString := ""
	for i, v := range stringSlice {
		if i == (len(stringSlice) - 1) {
			resultString = AppendStrings(resultString, enclosure, v, enclosure)
		} else {
			resultString = AppendStrings(resultString, enclosure, v, enclosure, delimiter)
		}
	}
	return resultString
}

func MakePrepareForInExpression(target interface{}) (string, map[string]interface{}) {
	bindParams := make(map[string]interface{})
	query := ""
	rv := reflect.ValueOf(target)

	switch rv.Kind() {
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if len(query) > 0 {
				query += ", "
			}
			bindParams["id"+strconv.Itoa(i)] = rv.Index(i)
			query += ":id" + strconv.Itoa(i)
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if len(query) > 0 {
				query += ", "
			}
			bindParams["id"+strconv.Itoa(i)] = rv.Index(i).Interface()
			query += ":id" + strconv.Itoa(i)
		}
	}
	return query, bindParams
}

func RemoveDuplicate(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

func MergeMap(baseMap map[string]interface{}, mergeMaps ...map[string]interface{}) map[string]interface{} {
	for _, mergeMap := range mergeMaps {
		for k, v := range mergeMap {
			baseMap[k] = v
		}
	}
	return baseMap
}

func SearchStringValueInSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
