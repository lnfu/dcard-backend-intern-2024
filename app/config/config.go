package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Mode     string
	Address  string
	Database Database
	Redis    Redis
}

type Database struct {
	Driver string
	Source string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func Init(mode string) *Config {

	conf := Config{}
	conf.Database.Driver = "mysql"
	conf.Redis.Password = ""
	conf.Redis.DB = 0

	switch mode {
	case "prod":
		gin.SetMode(gin.ReleaseMode)

		conf.Database.Source = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			"mysql",
			os.Getenv("MYSQL_DATABASE"),
		)
		conf.Redis.Addr = "redis:6379"
	case "dev":
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		conf.Database.Source = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			"localhost",
			os.Getenv("MYSQL_DATABASE"),
		)
		conf.Redis.Addr = "localhost:6379"
	default: // dev
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		conf.Database.Source = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			"localhost",
			os.Getenv("MYSQL_DATABASE"),
		)
		conf.Redis.Addr = "localhost:6379"
	}

	return &conf
}
