package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

func CSRF() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("CSRF_AUTH_KEY")),
		csrf.Secure(true),
		csrf.HttpOnly(true),
	)
	return adapter.Wrap(csrfMiddleware)
}