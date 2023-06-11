package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	expiration = flag.Int64("e", 30, "Days until token expiration")
	subject    = flag.String("s", "jay", "Subject of the token")
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Sub string `json:"sub"`
}

func main() {
	flag.Parse()

	tokenHeader := header{
		Alg: "HS256",
		Typ: "JWT",
	}

	tokenPayload := payload{
		Exp: time.Now().AddDate(0, 0, int(*expiration)).Unix(),
		Iat: time.Now().Unix(),
		Sub: *subject,
	}

	headerJson, _ := json.Marshal(tokenHeader)
	payloadJson, _ := json.Marshal(tokenPayload)

	// Base64 encoding
	headerB64 := base64.StdEncoding.EncodeToString(headerJson)
	payloadB64 := base64.StdEncoding.EncodeToString(payloadJson)

	signingString := headerB64 + "." + payloadB64

	// HMAC-SHA256 signing
	mac := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	mac.Write([]byte(signingString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Assembling the final token
	token := signingString + "." + signature

	fmt.Println("Generated JWT:", token)
}
