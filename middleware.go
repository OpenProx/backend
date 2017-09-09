package backend

import (
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ParseJWT parses the JWT Token if present
func (i *Instance) ParseJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("User", nil)

		rawToken := c.Request().Header.Get("X-JWT-Token")
		token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWTSecret")), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var u User
				if err := i.Database.One("ID", int(claims["id"].(float64)), &u); err == nil {
					c.Set("User", &u)
				}
			}
		}

		return next(c)
	}
}
