package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// secret key used to sign tokens
var (
	SecretKey = []byte("secret")
)

func GenerateToken(uid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// create a map to store claims
	claims := token.Claims.(jwt.MapClaims)
	// set token claims
	claims["user_id"] = uid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating Key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, val := range claims {
			println(key, "====>", val)
		}
		uid := claims["user_id"].(string)
		return uid, nil
	} else {
		return "", err
	}
}
