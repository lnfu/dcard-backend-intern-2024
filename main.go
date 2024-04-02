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

	app := &application{
		errorLogger:     errorLogger,
		infoLoggger:     infoLoggger,
		databaseQueries: db.New(dbConnection),
	}

	router := newRouter(app)

	router.Run(*addr)
}

func newRouter(app *application) *gin.Engine {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	apiV1 := router.Group("api/v1/")
	apiV1.POST("ad", app.createAdvertisementHandler)
	apiV1.GET("ad", app.getAdvertisementHandler)

	return router
}
