package backend

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// InitRouter inits the router of the server
func (i *Instance) InitRouter() {
	i.Router = echo.New()
	i.Router.Logger.SetOutput(ioutil.Discard)

	i.Router.Use(i.ParseJWT)

	i.Router.GET("user", i.GetUserRoute)
	i.Router.GET("new", i.GetNewIdentityRoute)
	i.Router.GET("check", i.GetCheckRequestRoute)

	i.Router.POST("submit", i.PostProxies)
	i.Router.POST("check", i.PostCheckResultRoute)
}

// GetUserRoute returns user infos
func (i *Instance) GetUserRoute(c echo.Context) error {
	u := c.Get("User")
	if u == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, c.Get("User"))
}

// GetNewIdentityRoute creates a new identity
func (i *Instance) GetNewIdentityRoute(c echo.Context) error {
	new := CreateUser()
	err := i.Database.Save(new)
	if err != nil {
		i.Log.WithFields(logrus.Fields{
			"ID": new.ID,
		}).WithError(err).Error("Error while new identity save")
		return c.NoContent(http.StatusBadRequest)
	}

	claims := struct {
		ID int `json:"id"`
		jwt.StandardClaims
	}{
		new.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
			Issuer:    "Server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWTSecret")))
	if err != nil {
		i.Log.WithFields(logrus.Fields{
			"ID": new.ID,
		}).WithError(err).Error("Error while new identity save")
		return c.NoContent(http.StatusBadRequest)
	}

	i.Log.WithFields(logrus.Fields{
		"ID": new.ID,
	}).Info("New identity created")
	return c.JSON(http.StatusOK, map[string]string{
		"token": ss,
	})
}

// PostProxies submits new proxies
func (i *Instance) PostProxies(c echo.Context) error {
	u := c.Get("User").(*User)
	if u == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	d := struct {
		Proxies []string `json:"proxies"`
	}{}
	c.Bind(&d)

	if len(d.Proxies) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	if len(i.IncomingProxy)+len(d.Proxies) > cap(i.IncomingProxy) {
		return c.NoContent(http.StatusTooManyRequests)
	}

	i.Log.WithFields(logrus.Fields{
		"ID":      u.ID,
		"Proxies": len(d.Proxies),
	}).Info("Proxies added to queue")
	i.IncomingProxy <- AddRequest{Proxies: d.Proxies, By: u.ID}
	return c.NoContent(http.StatusOK)
}

// GetCheckRequestRoute returns a check request
func (i *Instance) GetCheckRequestRoute(c echo.Context) error {
	u := c.Get("User").(*User)
	if u == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	req, err := i.GetCheckableProxy(u.ID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if req == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, req)
}

// PostCheckResultRoute posts the user check result
func (i *Instance) PostCheckResultRoute(c echo.Context) error {
	u := c.Get("User").(*User)
	if u == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	chk := CheckResult{}
	c.Bind(&chk)

	if len(chk.Token) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	i.IncomingResult <- chk
	return c.NoContent(http.StatusOK)
}
