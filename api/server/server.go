package server

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/skanehira/vue-go-oauth2/api/config"
	"github.com/skanehira/vue-go-oauth2/api/handler"
)

// Server server
type Server struct {
	config *config.Config
	db     *gorm.DB
	e      *echo.Echo
}

// New new server
func New(config *config.Config, db *gorm.DB, e *echo.Echo) *Server {
	return &Server{
		config: config,
		db:     db,
		e:      e,
	}
}

// Start start server
func (s *Server) Start() {
	s.InitHandler()
	// Start server
	s.e.Logger.Fatal(s.e.Start(":" + s.config.Port))
}

// InitHandler add handler
func (s *Server) InitHandler() {

	// TODO chnage dir
	s.e.Static("/", "../front/dist/")

	// session setting
	// TODO set keyPairs from env
	s.e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	//s.e.Use(middleware.Logger())

	// init handlers
	// u := handler.NewUserHandler(s.db)
	oauth := handler.NewOAuthHandler(s.config, s.db)
	//	s.e.GET("/users/:userID", u.GetUser())
	//	s.e.PUT("/users/:userID", u.UpdateUser())
	//	s.e.DELETE("/users/:userID", u.DeleteUser())
	s.e.POST("/users/signin", oauth.Signin())
	s.e.GET("/twitter/callback", oauth.TwitterCallback())
}
