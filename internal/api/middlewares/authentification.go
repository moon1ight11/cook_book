package middlewares

import (
	"net/http"

	"cook_book/internal/api/jwt"
	"cook_book/pkg/logger"

	"github.com/labstack/echo/v4"
)

func Auth(jwtService jwt.JWTInterface, logger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil || cookie.Value == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
			}

			claims, err := jwtService.ValidateAccessToken(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			for _, role := range roles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "access violation")
		}
	}
}
