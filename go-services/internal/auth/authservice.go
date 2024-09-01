package auth

import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"
    "golang.org/x/oauth2"
)

type Config struct {
    GoogleOAuthConfig oauth2.Config
    GithubOAuthConfig oauth2.Config
    JWTSecret         []byte
    DB                *sql.DB
    EmailSender       EmailSender
}

type AuthService struct {
    config Config
}

func NewAuthService(config Config) *AuthService {
    return &AuthService{config: config}
}

func (s *AuthService) RegisterRoutes(e *echo.Echo) {
    e.GET("/auth/login/:provider", s.HandleOAuthLogin)
    e.GET("/auth/callback/:provider", s.HandleOAuthCallback)
    e.POST("/auth/login", s.HandleLocalLogin)
    e.POST("/auth/logout", s.HandleLogout)
    e.POST("/auth/register", s.HandleRegister)
    e.POST("/auth/verify-email", s.HandleVerifyEmail)
    e.POST("/auth/request-password-reset", s.HandleRequestPasswordReset)
    e.POST("/auth/reset-password", s.HandleResetPassword)
    e.POST("/auth/enable-2fa", s.HandleEnable2FA)
    e.POST("/auth/verify-2fa", s.HandleVerify2FA)
}

func (s *AuthService) HandleOAuthLogin(c echo.Context) error {
    provider := c.Param("provider")
    state, err := generateState()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
    }

    var authURL string
    switch OAuthProvider(provider) {
    case Google:
        authURL = s.config.GoogleOAuthConfig.AuthCodeURL(state)
    case GitHub:
        authURL = s.config.GithubOAuthConfig.AuthCodeURL(state)
    default:
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid OAuth provider")
    }

    if err := s.storeState(c.Response(), state); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to store state")
    }

    return c.Redirect(http.StatusFound, authURL)
}

// Adapt other handler functions similarly...

func (s *AuthService) HandleLocalLogin(c echo.Context) error {
    var loginRequest struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        TOTPToken string `json:"totp_token"`
    }

    if err := c.Bind(&loginRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    user, err := s.authenticateUser(loginRequest.Email, loginRequest.Password)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
    }

    if user.TwoFactorSecret != "" {
        if loginRequest.TOTPToken == "" {
            return echo.NewHTTPError(http.StatusUnauthorized, "2FA token required")
        }
        if !totp.Validate(loginRequest.TOTPToken, user.TwoFactorSecret) {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid 2FA token")
        }
    }

    token, err := s.createJWT(user)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
    }

    c.SetCookie(&http.Cookie{
        Name:     "auth_token",
        Value:    token,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    return c.JSON(http.StatusOK, user)
}

func (s *AuthService) HandleOAuthCallback(c echo.Context) error {
    provider := c.Param("provider")
    state := c.QueryParam("state")
    code := c.QueryParam("code")

    if err := s.verifyState(c.Request(), state); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid state")
    }

    var token *oauth2.Token
    var err error
    var userInfo map[string]interface{}

    switch OAuthProvider(provider) {
    case Google:
        token, err = s.config.GoogleOAuthConfig.Exchange(c.Request().Context(), code)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to exchange token")
        }
        userInfo, err = s.getUserInfoFromGoogle(token)
    case GitHub:
        token, err = s.config.GithubOAuthConfig.Exchange(c.Request().Context(), code)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to exchange token")
        }
        userInfo, err = s.getUserInfoFromGitHub(token)
    default:
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid OAuth provider")
    }

    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user info")
    }

    user, err := s.createOrUpdateUser(userInfo, string(provider))
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create or update user")
    }

    // Create JWT for user
    jwtToken, err := s.createJWT(user)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create token")
    }

    c.SetCookie(&http.Cookie{
        Name:     "auth_token",
        Value:    jwtToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    return c.Redirect(http.StatusFound, "/dashboard")
}

func (s *AuthService) HandleLogout(c echo.Context) error {
    c.SetCookie(&http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Expires:  time.Now().Add(-1 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    return c.NoContent(http.StatusOK)
}

func (s *AuthService) HandleRegister(c echo.Context) error {
    var registerRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.Bind(&registerRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
    }

    _, err = s.config.DB.Exec("INSERT INTO users (email, hashed_password, email_verified) VALUES ($1, $2, $3)",
        registerRequest.Email, hashedPassword, false)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
    }

    token, err := s.generateToken()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate verification token")
    }

    err = s.config.EmailSender.SendVerificationEmail(registerRequest.Email, token)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send verification email")
    }

    return c.NoContent(http.StatusCreated)
}

func (s *AuthService) HandleVerifyEmail(c echo.Context) error {
    var verifyRequest struct {
        Email string `json:"email"`
        Token string `json:"token"`
    }

    if err := c.Bind(&verifyRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    // In a real implementation, you'd verify the token against a stored value
    _, err := s.config.DB.Exec("UPDATE users SET email_verified = true WHERE email = $1", verifyRequest.Email)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify email")
    }

    return c.NoContent(http.StatusOK)
}

func (s *AuthService) HandleRequestPasswordReset(c echo.Context) error {
    var resetRequest struct {
        Email string `json:"email"`
    }

    if err := c.Bind(&resetRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    token, err := s.generateToken()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate reset token")
    }

    err = s.config.EmailSender.SendPasswordResetEmail(resetRequest.Email, token)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send reset email")
    }

    return c.NoContent(http.StatusOK)
}

func (s *AuthService) HandleResetPassword(c echo.Context) error {
    var resetRequest struct {
        Email       string `json:"email"`
        Token       string `json:"token"`
        NewPassword string `json:"new_password"`
    }

    if err := c.Bind(&resetRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    // In a real implementation, you'd verify the token against a stored value
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetRequest.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
    }

    _, err = s.config.DB.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, resetRequest.Email)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reset password")
    }

    return c.NoContent(http.StatusOK)
}

func (s *AuthService) HandleEnable2FA(c echo.Context) error {
    user, err := s.getUserFromContext(c)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
    }

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "YourApp",
        AccountName: user.Email,
    })
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate 2FA secret")
    }

    _, err = s.config.DB.Exec("UPDATE users SET two_factor_secret = $1 WHERE id = $2", key.Secret(), user.ID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enable 2FA")
    }

    return c.JSON(http.StatusOK, map[string]string{
        "secret": key.Secret(),
        "qr_url": key.URL(),
    })
}

func (s *AuthService) HandleVerify2FA(c echo.Context) error {
    var verifyRequest struct {
        Token string `json:"token"`
    }

    if err := c.Bind(&verifyRequest); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    user, err := s.getUserFromContext(c)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
    }

    if !totp.Validate(verifyRequest.Token, user.TwoFactorSecret) {
        return echo.NewHTTPError(http.StatusUnauthorized, "Invalid 2FA token")
    }

    return c.NoContent(http.StatusOK)
}

func (s *AuthService) getUserFromContext(c echo.Context) (*User, error) {
    user, ok := c.Get("user").(*User)
    if !ok {
        return nil, errors.New("user not found in context")
    }
    return user, nil
}

func (s *AuthService) verifyState(r *http.Request, state string) error {
    cookie, err := r.Cookie("oauth_state")
    if err != nil {
        return err
    }
    if cookie.Value != state {
        return errors.New("invalid state")
    }
    return nil
}

func (s *AuthService) storeState(c echo.Context, state string) error {
    cookie := new(http.Cookie)
    cookie.Name = "oauth_state"
    cookie.Value = state
    cookie.Expires = time.Now().Add(15 * time.Minute)
    cookie.HttpOnly = true
    cookie.Secure = true
    cookie.SameSite = http.SameSiteStrictMode
    c.SetCookie(cookie)
    return nil
}

func (s *AuthService) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *AuthService) getUserInfoFromGoogle(token *oauth2.Token) (map[string]interface{}, error) {
	client := s.config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *AuthService) getUserInfoFromGitHub(token *oauth2.Token) (map[string]interface{}, error) {
	client := s.config.GithubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *AuthService) createJWT(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.config.JWTSecret)
}