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

// TODO 讀取 dotenv (https://blog.wu-boy.com/2019/04/how-to-load-env-file-in-go/)

func Init(env string) *Config {
	conf := Config{}

	switch env {
	case "dev":
		conf.Address = ":8080"

		conf.Database.Driver = "mysql"
		conf.Database.Source = fmt.Sprintf(
			"%s:%s@/%s?parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASSWORD"),
			os.Getenv("MYSQL_DATABASE"),
		)

		conf.Redis.Addr = "localhost:6379"
		conf.Redis.Password = ""
		conf.Redis.DB = 0
	}

	return &conf
}
