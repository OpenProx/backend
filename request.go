package backend

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateRequestToken generates a token for a check request
func GenerateRequestToken(proxy, uid, checkid int) (string, error) {
	claims := struct {
		Proxy   int `json:"proxy"`
		ID      int `json:"id"`
		CheckID int `json:"checkid"`
		jwt.StandardClaims
	}{
		proxy,
		uid,
		checkid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			Issuer:    "Server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWTSecret")))
}

// DecodeRequestToken decodes a token for a check request
func DecodeRequestToken(ptoken string) (int, int, int, error) { // TODO: Return Values to Struct!
	token, err := jwt.Parse(ptoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWTSecret")), nil
	})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return int(claims["proxy"].(float64)), int(claims["id"].(float64)), int(claims["checkid"].(float64)), nil
		}
	}

	return 0, 0, 0, err
}
