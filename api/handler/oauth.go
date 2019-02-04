package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gomodule/oauth1/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/skanehira/vue-go-oauth2/api/common"
	"github.com/skanehira/vue-go-oauth2/api/config"
	"github.com/skanehira/vue-go-oauth2/api/model"

	"github.com/jinzhu/gorm"
)

const (
	// twitter callback url
	twitterCallBack = "http://localhost:8080/twitter/callback"
	// session token key
	tokenKey  = "tokenKey"
	secretKey = "secretKey"
	// session name
	sessionName = "test_session"
	// user info url
	userInfoURI = "https://api.twitter.com/1.1/account/verify_credentials.json"
	// twitter base uri
	baseURI = "https://twitter.com/"
)

// OAuth oauth
type OAuth struct {
	client oauth.Client
	db     *gorm.DB
}

// NewOAuthHandler new oauth
func NewOAuthHandler(config *config.Config, db *gorm.DB) *OAuth {
	// new oauth setting from cnofig.yaml
	return &OAuth{
		client: oauth.Client{
			TemporaryCredentialRequestURI: config.Twitter.RequestURI,
			ResourceOwnerAuthorizationURI: config.Twitter.AuthorizationURI,
			TokenRequestURI:               config.Twitter.TokenRequestURI,
			Credentials: oauth.Credentials{
				Token:  config.Twitter.Token,
				Secret: config.Twitter.Secret,
			},
		},
		db: db,
	}
}

// Signin signin
func (o *OAuth) Signin() echo.HandlerFunc {
	return func(c echo.Context) error {
		credentials, err := o.client.RequestTemporaryCredentials(nil, twitterCallBack, nil)
		if err != nil {
			err := fmt.Errorf("Error getting credentials, %s", err.Error())
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}
		sess, err := session.Get(sessionName, c)
		if err != nil {
			err := fmt.Errorf("Error getting credentials, %s", err.Error())
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		sess.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}

		sess.Values[tokenKey] = credentials.Token
		sess.Values[secretKey] = credentials.Secret
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			err := fmt.Errorf("Error saving session, %s", err.Error())
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		return c.JSON(http.StatusOK, struct{ URL string }{o.client.AuthorizationURL(credentials, nil)})
	}
}

// TwitterCallback twitter callback
func (o *OAuth) TwitterCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get credentials from session
		sess, _ := session.Get(sessionName, c)

		if sess.IsNew {
			err := fmt.Errorf("%s", "Error invalid session.")
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		token, ok := sess.Values[tokenKey]
		if !ok {
			err := fmt.Errorf("%s", "Error getting token.")
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		secret, ok := sess.Values[secretKey]
		if !ok {
			err := fmt.Errorf("%s", "Error getting secret.")
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		credentials := &oauth.Credentials{
			Token:  token.(string),
			Secret: secret.(string),
		}

		// if credentials is nil and toeken is not equles
		if credentials.Token != c.QueryParam("oauth_token") {
			err := fmt.Errorf("%s", "Unknown oauth_token.")
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		// get access token.
		accessCredentials, _, err := o.client.RequestToken(nil, credentials, c.QueryParam("oauth_verifier"))

		if err != nil {
			err := fmt.Errorf("Error getting request token, %s", err)
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		// save access token to session
		sess.Values[tokenKey] = accessCredentials.Token
		sess.Values[secretKey] = accessCredentials.Secret

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			err := fmt.Errorf("Error saving session, %s", err)
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		// get twitter account info to save db
		user, err := o.GetAccountInfo(accessCredentials)
		if err != nil {
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		user.URL = baseURI + user.ScreenName

		if err := model.New(o.db).SaveUser(user); err != nil {
			return c.JSON(common.GetErrorCode(err), common.NewError(err.Error()))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (o OAuth) GetAccountInfo(credentials *oauth.Credentials) (model.User, error) {
	var user model.User

	if err := o.APIGet(
		credentials,
		userInfoURI,
		url.Values{"include_entities": {"true"}},
		&user); err != nil {

		return user, fmt.Errorf("Error getting timeline, %s", err.Error())
	}

	return user, nil
}

func (o OAuth) APIGet(cred *oauth.Credentials, urlStr string, form url.Values, data interface{}) error {
	resp, err := o.client.Get(nil, cred, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeResponse(resp, data)
}

func decodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != http.StatusOK {
		p, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("get %s returned status %d, %s", resp.Request.URL, resp.StatusCode, p)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
