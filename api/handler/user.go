package handler

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/skanehira/pgw/api/common"
	"github.com/skanehira/pgw/api/model"
)

// UserHandler user handler strcut
type UserHandler struct {
	model *model.Model
}

// NewUserHandler new handler
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{model: model.New(db)}
}

// UpdateUser update user info
func (u *UserHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		newUser := model.User{}

		// get post from data
		if err := c.Bind(&newUser); err != nil {
			return c.JSON(common.GetErrorCode(common.ErrInvalidPostData), common.NewError(common.ErrInvalidPostData.Error()))
		}
		newUser.ID = c.Param("userID")

		//  update user
		newUser, err := u.model.UpdateUser(newUser)

		if err != nil {
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, newUser)
	}
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
