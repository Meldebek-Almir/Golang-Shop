package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var secretKey = func() string {
	body, err := os.ReadFile("./secret.txt")
	if err != nil {
		fmt.Println(err)
		return "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasfdasdfsadfasdf"
	}
	return string(body)
}

func CreateJWT(claims map[string]string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	claims["iat"] = time.Now().Format(time.RFC3339)
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	data := []byte(encodedHeader + "." + encodedClaims)
	h := hmac.New(sha256.New, []byte(secretKey()))
	_, err = h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	token := strings.Join([]string{encodedHeader, encodedClaims, signature}, ".")

	return token, nil
}

func VerifyJWT(token string) (bool, map[string]string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, nil, errors.New("len not 3")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false, nil, err
	}
	var header map[string]interface{}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return false, nil, err
	}

	var claims map[string]string
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, nil, err
	}
	err = json.Unmarshal(claimsBytes, &claims)
	if err != nil {
		return false, nil, err
	}
	if header["alg"] != "HS256" {
		return false, nil, errors.New("invalid algorithm type")
	}

	data := []byte(parts[0] + "." + parts[1])
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false, nil, err
	}
	h := hmac.New(sha256.New, []byte(secretKey()))
	_, err = h.Write(data)
	if err != nil {
		return false, nil, err
	}
	if !hmac.Equal(signature, h.Sum(nil)) {
		return false, nil, errors.New("dont equal")
	}

	exp, ok := claims["exp"]
	if !ok {
		return false, nil, errors.New("dont have exp")
	}

	expiry, err := time.Parse(time.RFC3339, exp)
	if err != nil {
		return false, nil, err
	}

	if expiry.Before(time.Now()) {
		return false, nil, errors.New("time before")
	}

	return true, claims, nil
}
