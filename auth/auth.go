package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"silic0n-wiki/config"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateToken(byteLength int) (string, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func SignToken(token string) string {
	mac := hmac.New(sha256.New, []byte(config.AppConfig.Secret))
	mac.Write([]byte(token))
	signature := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s.%s", token, signature)
}

func VerifySignedToken(signedToken string) (string, bool) {
	idx := strings.LastIndex(signedToken, ".")
	if idx <= 0 || idx == len(signedToken)-1 {
		return "", false
	}

	token := signedToken[:idx]
	providedSig := signedToken[idx+1:]

	mac := hmac.New(sha256.New, []byte(config.AppConfig.Secret))
	mac.Write([]byte(token))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(providedSig), []byte(expectedSig)) {
		return "", false
	}

	return token, true
}

func GenerateCSRFToken(sessionToken string) string {
	mac := hmac.New(sha256.New, []byte(config.AppConfig.Secret))
	mac.Write([]byte("csrf:" + sessionToken))
	return hex.EncodeToString(mac.Sum(nil))
}

func ValidateCSRFToken(csrfToken, sessionToken string) bool {
	expected := GenerateCSRFToken(sessionToken)
	return hmac.Equal([]byte(csrfToken), []byte(expected))
}
