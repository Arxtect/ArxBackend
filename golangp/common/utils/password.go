package utils

import (
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/toheart/functrace"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	defer functrace.Trace([]interface {
	}{password})()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	defer functrace.Trace([]interface {
	}{hashedPassword, candidatePassword})()
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func GenerateRandomString(length int) string {
	defer functrace.Trace([]interface {
	}{length})()
	randomBytes := make([]byte, (length*5+7)/8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		logger.Warning("Error while generating random bytes for prompt %s", err.Error())
		return ""
	}

	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return encoded[:length]
}
