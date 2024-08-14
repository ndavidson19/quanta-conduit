package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	adapter "github.com/gwatts/gin-adapter"
)

func SecureHeaders() gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
		SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})

	return adapter.Wrap(secureMiddleware.Handler)
}