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

// GenerateClientSecret generate a clientSecret
func GenerateClientSecret(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = token68Letters[rand.Intn(len(token68Letters))]
	}
	return string(b)
}

// GenerateUUID generate a UUID
func GenerateUUID() string {
	uuid := uuid.NewV4().String()
	return uuid
}

// GenerateClientID generate a clientID
func GenerateClientID() string {
	uuid := uuid.NewV4().String()
	return strings.Replace(uuid, "-", "", -1)
}

// GetFileNameWithoutExt get a filename without extension
func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

// BasicAuth generate a basic authorization
func BasicAuth(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
