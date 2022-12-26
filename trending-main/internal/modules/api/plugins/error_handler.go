package plugins

import (
	"encoding/json"
	"moonbite/trending/internal/modules/api/schemas"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func NewErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {

		// Logging
		logrus.WithField("err", err).Errorf("server error")

		if c.Response().Committed {
			return
		}

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		// Issue #1426
		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			if e.Debug {
				//message = echo.Map{"message": m, "error": err.Error()}
				message = err.Error()
			} else {
				//message = echo.Map{"message": m}
				message = m
			}
		}

		// Reformat error response
		ce := &schemas.Error{Message: message}

		// Validator error
		if ve, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string][]string)
			for _, err := range ve {
				key := err.Field()
				val := strings.ToUpper(err.Tag())
				if _, ok := errors[key]; !ok {
					errors[key] = make([]string, 0)
				}
				errors[key] = append(errors[key], val)
			}

			if len(errors) > 0 {
				code = http.StatusBadRequest
				ce.Message = nil
				ce.ValidationErrors = errors
			}
		}

		// Process unmarhsal error and convert as validation error
		if code == http.StatusBadRequest && he.Internal != nil {
			if ve, ok := he.Internal.(*json.UnmarshalTypeError); ok {
				ce.Message = nil
				ce.ValidationErrors = map[string][]string{
					ve.Field: {"INVALID_TYPE"},
				}
			}
		}

		// Send response
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, ce)
		}
		if err != nil {
			e.Logger.Error(err)
		}
	}
}
