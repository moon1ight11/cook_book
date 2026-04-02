package middlewares

import (
	"cook_book/internal/api/jwt"
	"cook_book/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

// аус-мидлвар
func Auth(jwtService jwt.JWTInterface, logger logger.Logger) echo.MiddlewareFunc {
	// возвращаем как мидлвар
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// получаем аксес токен в куки
			cookie, err := c.Cookie("access_token")
			if err != nil || cookie.Value == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
			}

			tokenString := cookie.Value

			// достаем и валидируем access token
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// устанавливаем в контекст id и роль
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

// role-check мидлвар
func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// получаем роль из контекста
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			// сравниваем с разрешенной
			if userRole == role {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden, "access violation")
		}
	}
}
