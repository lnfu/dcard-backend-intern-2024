package config

import (
	"fmt"
	"os"
)

type Config struct {
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

func Init(env string) *Config {
	conf := Config{}

	switch env {
	case "dev":
		conf.Address = ":8080"

		conf.Database.Driver = "mysql"
		conf.Database.Source = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			"localhost",
			os.Getenv("MYSQL_DATABASE"),
		)

		conf.Redis.Addr = "localhost:6379"
		conf.Redis.Password = ""
		conf.Redis.DB = 0
	case "prod":
		conf.Address = ":8080"

		conf.Database.Driver = "mysql"
		conf.Database.Source = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			"mysql",
			os.Getenv("MYSQL_DATABASE"),
		)

		conf.Redis.Addr = "redis:6379"
		conf.Redis.Password = ""
		conf.Redis.DB = 0
	}

	return &conf
}
