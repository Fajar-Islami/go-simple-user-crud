package http

import (
	"net/http"
	"strings"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("authorization")
		if token == "" {
			return helper.BuildResponse(c, false, "UNATHORIZED", "FAILED TO GET TOKEN", nil, http.StatusUnauthorized)
		}

		s := strings.Split(token, " ")
		if len(s) < 2 {
			return helper.BuildResponse(c, false, "UNATHORIZED", "FAILED TO GET TOKEN", nil, http.StatusUnauthorized)
		}

		if s[1] != "3cdcnTiBsl" {
			return helper.BuildResponse(c, false, "UNATHORIZED", "TOKEN INVALID", nil, http.StatusUnauthorized)
		}

		// Go to next middleware:
		return next(c)
	}
}
