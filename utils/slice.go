package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Join is join strings
func Join(stringSlice []string, delimiter string, enclosure string) string {
	resultString := ""
	for i, v := range stringSlice {
		if i == (len(stringSlice) - 1) {
			resultString = fmt.Sprintf("%s%s%s%s", resultString, enclosure, v, enclosure)
		} else {
			resultString = fmt.Sprintf("%s%s%s%s%s", resultString, enclosure, v, enclosure, delimiter)
		}
	}
	return resultString
}

// RemoveDuplicateString removes duplicate string value
func RemoveDuplicateString(args []string) []string {
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

// RemoveDuplicateInt32 removes duplicate int32 value
func RemoveDuplicateInt32(args []int32) []int32 {
	results := make([]int32, 0, len(args))
	encountered := map[int32]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

// MergeMap merges map
func MergeMap(baseMap map[string]interface{}, mergeMaps ...map[string]interface{}) map[string]interface{} {
	if baseMap == nil {
		return nil
	}

	for _, mergeMap := range mergeMaps {
		for k, v := range mergeMap {
			baseMap[k] = v
		}
	}
	return baseMap
}

// SearchStringValueInSlice is search string value in slice
func SearchStringValueInSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func CommaSeparatedStringsToInt32(v string) []int32 {
	vv := strings.Split(v, ",")
	int32Sli := make([]int32, len(vv))
	if len(vv) == 0 {
		return int32Sli
	}

	for i, v := range vv {
		intValue, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			continue
		}
		int32Sli[i] = int32(intValue)
	}
	return int32Sli
}
