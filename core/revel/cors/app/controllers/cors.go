package controllers

import (
	"errors"
	"github.com/revel/revel"
	"strings"
)

var (
	allowedOrigins []string
)

func Init(origins string) {
	originsArr := strings.Split(origins, ",")
	allowedOrigins = originsArr
}

type Cors struct {
	*revel.Controller
}

func (c *Cors) Handle() revel.Result {
	return c.RenderText("")
}

func SetCORSHeaders(c *revel.Controller) error {
	header, exists := c.Request.Header["Origin"]
	if !exists {
		return errors.New("No Origin specified.")
	}

	origin := header[0]

	if allowed := isOriginAllowed(origin); !allowed {
		return errors.New("Origin not allowed.")
	}

	c.Response.Out.Header().Add("Access-Control-Allow-Methods", "POST,GET,DELETE")
	c.Response.Out.Header().Add("Access-Control-Allow-Origin", origin)
	c.Response.Out.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Response.Out.Header().Add("Access-Control-Allow-Credentials", "true")
	return nil
}

func isOriginAllowed(origin string) bool {
	for i := 0; i < len(allowedOrigins); i++ {
		switch allowedOrigins[i] {
		case "*", origin:
			return true
		}
	}
	return false
}
