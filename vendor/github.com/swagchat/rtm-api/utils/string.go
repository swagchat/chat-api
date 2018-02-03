package utils

import "regexp"

func AppendStrings(strings ...string) string {
	buf := make([]byte, 0)
	for _, str := range strings {
		buf = append(buf, str...)
	}
	return string(buf)
}

func IsValidId(id string) bool {
	r := regexp.MustCompile(`(?m)^[0-9a-zA-Z-]+$`)
	return r.MatchString(id)
}
