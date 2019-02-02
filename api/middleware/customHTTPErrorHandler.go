package middleware

import (
	"github.com/labstack/echo"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	if _, ok := err.(*echo.HTTPError); ok {

	}
}
