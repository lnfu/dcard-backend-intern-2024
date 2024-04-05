package config

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

func Init() *Config {
	// TODO 區分 dev/test/prod
	dev := Config{}

	dev.Address = ":8080"

	dev.Database.Driver = "mysql"
	dev.Database.Source = "web:pass@/dcard?parseTime=true"

	dev.Redis.Addr = "localhost:6379"
	dev.Redis.Password = ""
	dev.Redis.DB = 0

	return &dev
}
