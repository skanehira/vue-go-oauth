package common

import (
	"github.com/pkg/errors"
)

var (
	// ErrGetUserID cannot get user id error
	ErrGetUserID = errors.New("cannot get user id from session")
	// ErrGetCredentials cannot get user credentials form twitter
	ErrGetCredentials = errors.New("cannot get credentials")
	// ErrInvalidSession invalid session
	ErrInvalidSession = errors.New("invalid session")
	// ErrSaveSession cannot save sesison
	ErrSaveSession = errors.New("cannot save session")
	// ErrInvalidCredentials invalid credentials
	ErrInvalidCredentials = errors.New("invalid credentails")
	// ErrNotFoundUserInfo not found user info
	ErrNotFoundUserInfo = errors.New("not found user info")
	// ErrSaveUserInfo cannot save user info
	ErrSaveUserInfo = errors.New("cannnot save user info")
)

// ErrorMessage error struct
type ErrorMessage struct {
	Message string `json:"message"`
}

// NewError make new error
func NewError(err error, cause error) ErrorMessage {
	return ErrorMessage{errors.Wrap(err, cause.Error()).Error()}
}
