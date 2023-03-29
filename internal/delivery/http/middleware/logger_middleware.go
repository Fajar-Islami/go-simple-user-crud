package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func isLoggingSkip(c echo.Context) bool {
	requestPath := c.Request().URL.String()
	skipPath := map[string]bool{
		"/":                      true,
		"/api/v1/auth/login/":    true,
		"/api/v1/auth/register/": true,
	}

	return skipPath[requestPath]
}

func LoggerMiddleware(log *zerolog.Logger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			reqBody, err := hookRequest(c)
			if err != nil {
				log.Error().Err(err)
			}

			var jsonCompact bytes.Buffer
			if !isLoggingSkip(c) {
				err = json.Compact(&jsonCompact, reqBody)
				if err != nil {
					log.Error().Err(err)
				}
			}

			loggerRequest := log.With().
				Str("header", stringifyHeader(c)).
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				RawJSON("raw body", jsonCompact.Bytes()).
				Str("remote_ip", c.RealIP()).
				Logger()

			// Add logger to context
			c.Set("loggerRequest", &loggerRequest)

			// // Log pre-request message
			loggerRequest.Info().Msg("Pre-request")

			var resBod []byte
			err = middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
				Skipper: nil,
				Handler: func(c echo.Context, reqBody, resBody []byte) {
					resBod = resBody

				},
			})(next)(c)

			if err != nil {
				return err
			}

			loggerResponse := log.With().
				Str("header", stringifyHeader(c)).
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				Str("raw body", string(resBod)).
				Str("remote_ip", c.RealIP()).
				Logger()

			// Add logger to context
			c.Set("loggerResponse", &loggerResponse)

			// // Log pre-request message
			loggerResponse.Info().Msg("Post-request")

			return nil
		}
	}
}

func hookRequest(c echo.Context) (body []byte, err error) {
	if c.Request().Body != nil { // Read
		body, err = ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return body, err
		}
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}

func stringifyHeader(c echo.Context) (h string) {
	if headers := c.Request().Header; headers != nil {
		var temp []string
		for k, v := range headers {
			temp = append(temp, fmt.Sprintf("%s: %s", k, strings.Join(v, " ")))
		}
		return strings.Join(temp, ", ")
	}
	return
}
