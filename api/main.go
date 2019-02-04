package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/comail/colog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/skanehira/vue-go-oauth2/api/config"
	"github.com/skanehira/vue-go-oauth2/api/model"
	"github.com/skanehira/vue-go-oauth2/api/server"
)

func main() {
	// get config
	config := config.New()

	// connect db
	// user:password@tcp(localhost:3306)/dbname?parseTime=true&charaset=utf8mb4,utf8
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4,utf8", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(config.DBLog)

	// db migrate
	flag.Parse()
	if len(flag.Args()) > 0 {
		if "create" == flag.Args()[0] {
			fmt.Println("create tables...")
			db.AutoMigrate(model.User{})
		} else if "drop" == flag.Args()[0] {
			fmt.Println("drop tables...")
			db.AutoMigrate(model.User{})
		}
		os.Exit(0)
	}

	// log setting
	// colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	// ログの設定方法
	// log.Printf("trace: this is a trace log.")
	// log.Printf("debug: this is a debug log.")
	// log.Printf("info: this is an info log.")
	// log.Printf("warn: this is a warning log.")
	// log.Printf("error: this is an error log.")
	// log.Printf("alert: this is an alert log.")
	// log.Printf("this is a default level log.")

	// start server
	server.New(config, db, echo.New()).Start()
}
