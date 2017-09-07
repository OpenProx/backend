package backend

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a token for a check request
func (i *Instance) GenerateToken(proxy, uid int) (string, error) {
	claims := struct {
		Proxy int    `json:"proxy"`
		ID    int    `json:"id"`
		Key   string `json:"key"`
		jwt.StandardClaims
	}{
		proxy,
		uid,
		fmt.Sprint(rand.Int63()),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			Issuer:    "Server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWTSecret")))
}
