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

// GetUser get user info
func (u *UserHandler) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get session
		sess, err := getSession(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrInvalidSession, err))
		}

		// get user id from session
		id, ok := sess.Values["id"]
		if !ok {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrGetUserID, nil))
		}

		// get user info with user id
		user, err := u.model.GetUser(id.(string))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrNotFoundUserInfo, err))
		}

		// return user info
		return c.JSON(http.StatusOK, user)
	}
}
