package handler

import (
	"fmt"

	"github.com/gomodule/oauth1/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

const (
	// session name
	sessionName = "test_session"
)

func newSession(c echo.Context) (*sessions.Session, bool, error) {
	session, err := session.Get(sessionName, c)
	if err != nil {
		return nil, false, fmt.Errorf("Error create new session, %s", err.Error())
	}

	if !session.IsNew {
		return nil, true, nil
	}

	return session, false, nil
}

func getSession(c echo.Context) (*sessions.Session, error) {
	session, err := session.Get(sessionName, c)
	if err != nil {
		return nil, err
	}

	if session.IsNew {
		return nil, fmt.Errorf("Error get session, invalid session")
	}

	return session, nil
}

func getCredentialsFromSession(session *sessions.Session) (*oauth.Credentials, error) {
	token, ok := session.Values[tokenKey]
	if !ok {
		return nil, fmt.Errorf("%s", "Error getting token.")
	}

	secret, ok := session.Values[secretKey]
	if !ok {
		return nil, fmt.Errorf("%s", "Error getting secret.")
	}

	credentials := &oauth.Credentials{
		Token:  token.(string),
		Secret: secret.(string),
	}

	return credentials, nil
}

func setValuesToSession(session *sessions.Session, data map[string]interface{}) error {
	for key, value := range data {
		session.Values[key] = value
	}

	return nil
}

func saveSession(c echo.Context, session *sessions.Session) error {
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return nil
}
