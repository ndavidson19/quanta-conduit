package middleware

import (
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "conduit/internal/auth"
    "github.com/golang-jwt/jwt"
    "golang.org/x/time/rate"
)

// Middleware wraps all custom middleware
type Middleware struct {
    authService *auth.AuthService
}

func NewMiddleware(authService *auth.AuthService) *Middleware {
    return &Middleware{
        authService: authService,
    }
}

func (m *Middleware) SetupMiddleware(e *echo.Echo) {
    e.Use(middleware.Recover())
    e.Use(middleware.RequestID())
    e.Use(m.SecureHeaders)
    e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
        TokenLookup: "header:X-CSRF-Token",
        CookieName:  "csrf",
        CookieMaxAge: 3600,
    }))
    e.Use(m.RateLimiter)
    e.Use(m.CSP)
    e.Use(m.JWTAuth)
}

// SecureHeaders adds security headers to every response
func (m *Middleware) SecureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
        c.Response().Header().Set("X-Frame-Options", "DENY")
        c.Response().Header().Set("X-Content-Type-Options", "nosniff")
        c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        return next(c)
    }
}

// RateLimiter implements a simple rate limiting middleware
func (m *Middleware) RateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(1*time.Second), 10) // 10 requests per second
    return func(c echo.Context) error {
        if !limiter.Allow() {
            return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
        }
        return next(c)
    }
}

// CSP adds Content Security Policy headers
func (m *Middleware) CSP(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("Content-Security-Policy", 
            "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; frame-ancestors 'none'; form-action 'self';")
        return next(c)
    }
}

// JWTAuth handles JWT authentication
func (m *Middleware) JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user, err := m.authService.GetUserFromRequest(c.Request())
        if err != nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or missing auth token")
        }
        c.Set("user", user)
        return next(c)
    }
}