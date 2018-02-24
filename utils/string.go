package utils

import (
	"math/rand"
	"regexp"
	"time"
)

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

var token68Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~+/")

func GenerateToken(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = token68Letters[rand.Intn(len(token68Letters))]
	}
	return string(b)
}
