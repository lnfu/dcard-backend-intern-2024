package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/lnfu/dcard-intern/db/sqlc"
	docs "github.com/lnfu/dcard-intern/docs"
	"github.com/lnfu/dcard-intern/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type application struct {
	errorLogger     *log.Logger
	infoLoggger     *log.Logger
	databaseQueries *db.Queries
}

const (
	dbDriver = "mysql"
	dbSource = "web:pass@/dcard?parseTime=true"
)

// @title Dcard Backend Intern 2024
// @version 1.0
// @Description 請⽤ Golang 設計並且實作⼀個簡化的廣告投放服務，該服務應該有兩個 API，⼀個⽤於產⽣廣告，⼀個⽤於列出廣告。每個廣告都有它出現的條件(例如跟據使⽤者的年齡)，產⽣廣告的 API ⽤來產⽣與設定條件。投放廣告的 API 就要跟據條件列出符合使⽤條件的廣告
// @Host localhost:8080
func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	flag.Parse()

	errorLogger := log.New(os.Stderr, color.RedString("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)
	infoLoggger := log.New(os.Stdout, color.BlueString("INFO\t"), log.Ldate|log.Ltime|log.Lshortfile)

	// MySQL Database
	dbConnection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer dbConnection.Close()

	// TODO app refactor (甚至拿掉?)
	app := &application{
		errorLogger:     errorLogger,
		infoLoggger:     infoLoggger,
		databaseQueries: db.New(dbConnection),
	}
	_ = app

	router := newRouter()

	handler := handlers.NewHandler(db.New(dbConnection))
	apiV1 := router.Group("api/v1/")
	apiV1.POST("ad", handler.CreateAdvertisementHandler)
	apiV1.GET("ad", handler.GetAdvertisementHandler)

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	apiV1.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(*addr)
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})
	return router
}
