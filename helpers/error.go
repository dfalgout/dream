package helpers

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func BindAndValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}

func NoResultsError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
