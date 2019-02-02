package handler

import (
	"github.com/gomodule/oauth1/oauth"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

const (
	// Session state keys.
	tempCredKey  = "tempCred"
	tokenCredKey = "tokenCred"

	// twitter callback url
	twitterCallBack = "http://localhost:8080/twitter/callback"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	Credentials: oauth.Credentials{
		// TODO move to config
		Token:  "",
		Secret: "",
	},
}

// OAuth oauth
type OAuth struct {
	db *gorm.DB
}

// NewOAuthHandler new oauth
func NewOAuthHandler(db *gorm.DB) *OAuth {
	return &OAuth{
		db: db,
	}
}

// Signin signin
func (o *OAuth) Signin() echo.HandlerFunc {
	return func(c echo.Context) error {
		credentials, err := oauthClient.RequestTemporaryCredentials(nil, twitterCallBack, nil)
		if err != nil {
			return c.String(500, "Error getting credentials, "+err.Error())
		}
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			//Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		sess.Values[tempCredKey] = credentials
		sess.Save(c.Request(), c.Response())

		return c.Redirect(302, oauthClient.AuthorizationURL(credentials, nil))
	}
}

// TwitterCallback twitter callback
func (o *OAuth) TwitterCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get credentials from session
		sess, _ := session.Get("session", c)
		credentials := sess.Values[tempCredKey].(*oauth.Credentials)

		// if credentials is nil and toeken is not equles
		if credentials == nil && credentials.Token != c.FormValue("oauth_token") {
			return c.String(500, "Unknown oauth_token.")
		}

		// get access token.
		accessCredentials, _, err := oauthClient.RequestToken(nil, credentials, c.FormValue("oauth_verifier"))
		if err != nil {
			return c.String(500, "Error getting request token, "+err.Error())
		}

		sess.Values[tokenCredKey] = accessCredentials
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return c.String(500, "Error saving session, "+err.Error())
		}

		return c.Redirect(302, "/")
	}
}
