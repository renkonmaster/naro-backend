package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"naro-backend/handler"
)

func main() {
	// .envファイルから環境変数を読み込み
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// データーベースの設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	conf := mysql.Config{
		User:      os.Getenv("DB_USERNAME"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	// データベースに接続
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	//usersテーブルが存在しなかったら、usersテーブルを作成する
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (Username VARCHAR(255) PRIMARY KEY, HashedPass VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}

	// セッションの情報を記憶するための場所をデータベース上に設定
	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil{
		log.Fatal(err)
	}

	h := handler.NewHandler(db)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	
	e.POST("/signup", h.SignUpHandler)
	e.POST("/login", h.LoginHandler)
	e.POST("/logout", h.LogoutHandler)
	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong")})

	withAuth := e.Group("") 
	withAuth.Use(handler.UserAuthMiddleware) 
	withAuth.GET("/me", handler.GetMeHandler)
	withAuth.GET("/cities/:cityName", h.GetCityInfoHandler)
	withAuth.POST("/cities", h.PostCityHandler)
	withAuth.GET("/cities", h.GetCitiesByCountryCodeHandler)
	withAuth.GET("/countries", h.GetCountryInfoHandler)

	err = e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func sumPopulationByCountryCode(cities []handler.City) map[string]int64 {
	result := make(map[string]int64)
	for _, city := range cities {
		if city.CountryCode.Valid {
			// まだmapに存在しなかった場合、初期化する
			if _, ok := result[city.CountryCode.String]; !ok {
				result[city.CountryCode.String] = 0
			}
			result[city.CountryCode.String] += city.Population.Int64
		}
	}
	return result
}