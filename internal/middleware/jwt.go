package middleware

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	UserId int
}

// const TOKEN_EXP = int64(time.Hour * 2) // 2 hours expiry

func CreateJwtToken(secret_key string) (string, error) {
	claims := jwt.MapClaims{
		"name":       "aryha",
		"expired_at": time.Now().Add(time.Hour * 2).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secret_key))
	if err != nil {
		log.Fatal(err)
	}
	return token, nil
}

func ValidateHMAC(r *http.Request, secret_key string) (bool, error) {
	x_api_key_id := r.Header.Get("X-Api-Key-ID")
	timestamp := r.Header.Get("Timestamp")
	initialSignature := r.Header.Get("Signature")

	method := r.Method

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonBody string
	if len(body) > 0 {
		jsonBody = string(body)
		jsonBody = strings.ReplaceAll(jsonBody, "\n", "")
		jsonBody = strings.ReplaceAll(jsonBody, "\\", "")
		jsonBody = strings.ReplaceAll(jsonBody, " ", "")
	} else {
		jsonBody = ""
	}

	fmt.Printf("JSON Body RAW: %v\n", jsonBody)
	fmt.Printf("Timestamp: %v\n", timestamp)
	fmt.Printf("ApiKeyID: %v\n", x_api_key_id)
	fmt.Printf("Method: %v\n", method)

	body_md5 := md5.Sum([]byte(jsonBody))
	body_md5_string := base64.StdEncoding.EncodeToString(body_md5[:])

	hmac_signature := timestamp + ":" + x_api_key_id + ":" + method + ":" + body_md5_string

	fmt.Printf("HMAC Signature: %v\n", hmac_signature)

	hmac := hmac.New(sha256.New, []byte(secret_key))
	hmac.Write([]byte(hmac_signature))
	expectedHMAC := hmac.Sum(nil)

	hmac_base64 := base64.StdEncoding.EncodeToString(expectedHMAC)

	fmt.Printf("HMAC Signature dalam bentuk base64: %v\n", hmac_base64)

	signature := "#" + x_api_key_id + ":#" + hmac_base64

	signature_base64 := base64.StdEncoding.EncodeToString([]byte(signature))

	fmt.Printf("Signature dalam bentuk base64: %v\n", signature_base64)
	fmt.Printf("Initial Signature: %v\n", initialSignature)

	if signature_base64 == initialSignature {
		return true, nil
	}

	return false, nil

}
