package main

import (
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	token, err := GenerateJWT("wen", "hello")
	if err != nil {
		log.Println("err: ", err)
		return
	}
	log.Println("token=", token)
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Token string `json:"token"`
	jwt.StandardClaims
}

func GenerateJWT(name, usertoken string) (jwttoken string, err error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		Name:  name,
		Token: usertoken, // just some name
		// Admin: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	jwttoken, err = token.SignedString([]byte("secret"))
	if err != nil {
		err = fmt.Errorf("token sign err: %v\n", err)
		return
	}
	return
}
