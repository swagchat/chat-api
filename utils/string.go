package utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var token68Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~+/")

// GenerateClientSecret is generate clientSecret
func GenerateClientSecret(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = token68Letters[rand.Intn(len(token68Letters))]
	}
	return string(b)
}

// GenerateUUID is generate UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	return uuid
}

// GenerateClientID is generate clientID
func GenerateClientID() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "", -1)
}

// GetFileNameWithoutExt is get filename without extention
func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

// BasicAuth is generate basic authorization
func BasicAuth(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
