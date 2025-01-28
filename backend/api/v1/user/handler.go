package user

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	auth "backend/auth"
)

type UserHandler struct {
	userStore Store
	jwtSecret []byte
}

func NewHandler(userS Store) (*UserHandler, error) {
	secret, ok := os.LookupEnv("Signing_Key")
	if !ok {
		return nil, errors.New("no secret key ")
	}

	return &UserHandler{
		userStore: userS,
		jwtSecret: []byte(secret),
	}, nil
}

func (h *UserHandler) Register(group *echo.Group) {
	skipper := func(c echo.Context) bool {
		// Skip middleware if path is equal 'login'
		if c.Request().URL.Path == "/login" || c.Request().URL.Path == "/signup" {
			return true
		}
		return false
	}

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(h.jwtSecret),
		Skipper:    skipper,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
	})

	group.POST("/signup", h.SignUp)
	group.POST("/login", h.Login)

	user := group.Group("/user", jwtMiddleware)
	user.GET("/me", h.GetMe)
	user.GET("/:id", h.GetUser)
	user.PUT("/:id", h.UpdateUser)
	user.PUT("/me", h.UpdateMe)
	user.DELETE("/:id", h.DeleteUser)
}
