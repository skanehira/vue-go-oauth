package handler

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/skanehira/vue-go-oauth2/api/common"
	"github.com/skanehira/vue-go-oauth2/api/model"
)

// UserHandler user handler strcut
type UserHandler struct {
	model *model.Model
}

// NewUserHandler new handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{model: model.New(db)}
}

// DeleteUser delete user info
func (u *UserHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("userID")

		// delete user
		if err := u.model.DeleteUser(userID); err != nil {
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetUser get user info
func (u *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("userID")

		// get user info
		u, err := u.model.GetUser(userID)

		if err != nil {
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, u)
	}
}
