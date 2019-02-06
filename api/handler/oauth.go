package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gomodule/oauth1/oauth"
	"github.com/labstack/echo"
	"github.com/skanehira/vue-go-oauth2/api/common"
	"github.com/skanehira/vue-go-oauth2/api/config"
	"github.com/skanehira/vue-go-oauth2/api/model"

	"github.com/jinzhu/gorm"
)

const (
	// session token key
	tokenKey  = "tokenKey"
	secretKey = "secretKey"
	// user info url
	userInfoURI = "https://api.twitter.com/1.1/account/verify_credentials.json"
	// twitter base uri
	baseURI = "https://twitter.com/"
)

// OAuth oauth
type OAuth struct {
	client oauth.Client
	db     *gorm.DB
	config *config.Config
}

// Response singin response
type Response struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
}

// NewOAuthHandler new oauth handler
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
		db:     db,
		config: config,
	}
}

// Signin signin twitter
func (o *OAuth) Signin() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get twitter access token
		credentials, err := o.client.RequestTemporaryCredentials(nil, o.config.Twitter.CallbackURI, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrGetCredentials, err))
		}

		// get session
		sess, isSignined, err := newSession(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrInvalidSession, err))
		}

		// if session is exsist return mypage
		if isSignined {
			return c.JSON(http.StatusOK, Response{http.StatusFound, "/mypage"})
		}

		// set session option
		sess.Options.HttpOnly = true

		// set request token value to session
		setValuesToSession(sess, map[string]interface{}{
			tokenKey:  credentials.Token,
			secretKey: credentials.Secret,
		})

		// save session
		if err := saveSession(c, sess); err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrNotFoundUserInfo, err))
		}

		// return twitter authp page url
		return c.JSON(http.StatusOK, Response{http.StatusOK, o.client.AuthorizationURL(credentials, nil)})
	}
}

// TwitterCallback twitter callback endpoint
func (o *OAuth) TwitterCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		// get session
		sess, err := getSession(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrInvalidSession, err))
		}

		// get request token
		credentials, err := getCredentialsFromSession(sess)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrInvalidCredentials, err))
		}

		// if request token is not equles
		if credentials.Token != c.QueryParam("oauth_token") {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrInvalidCredentials, err))
		}

		// get access token.
		accessCredentials, _, err := o.client.RequestToken(nil, credentials, c.QueryParam("oauth_verifier"))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrGetCredentials, err))
		}

		// get twitter account info
		user, err := o.GetUserInfo(accessCredentials)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrNotFoundUserInfo, err))
		}

		// generate wtitter account url
		user.URL = baseURI + user.ScreenName

		// set access token value to session
		setValuesToSession(sess, map[string]interface{}{
			tokenKey:  credentials.Token,
			secretKey: credentials.Secret,
			"id":      user.ID,
		})

		// save access token to session
		if err := saveSession(c, sess); err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrSaveSession, err))
		}

		// save user info to db
		if err := model.New(o.db).SaveUser(user); err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewError(common.ErrSaveUserInfo, err))
		}

		// redirect mypage
		return c.Redirect(http.StatusFound, "/#/mypage")
	}
}

// GetUserInfo get twitter user info
func (o OAuth) GetUserInfo(credentials *oauth.Credentials) (model.User, error) {
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

// APIGet call get twitter api
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
