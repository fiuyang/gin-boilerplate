package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	expire, err := time.ParseDuration(os.Getenv("JWT_EXPIRE"))
	if err != nil {
		expire = time.Hour // default value if parsing fails
	}

	maxAge, err := strconv.Atoi(os.Getenv("JWT_MAX_AGE"))
	if err != nil {
		maxAge = 0 // default value if parsing fails
	}

	emailPort, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		emailPort = 25 // default value if parsing fails
	}
	return &Config{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
		Swagger: Swagger{
			Host: os.Getenv("SWAGGER_HOST"),
			Url:  os.Getenv("SWAGGER_URL"),
			Mode: os.Getenv("SWAGGER_MODE"),
		},
		Kong: Kong{
			Url: os.Getenv("KONG_URL"),
		},
		Obs: ObsHuawei{
			Ak:       os.Getenv("OBS_HUAWEI_AK"),
			Sk:       os.Getenv("OBS_HUAWEI_SK"),
			Endpoint: os.Getenv("OBS_HUAWEI_ENDPOINT"),
			Bucket:   os.Getenv("OBS_HUAWEI_BUCKET"),
		},
		Jwt: Jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Expire: expire,
			MaxAge: maxAge,
		},
		Email: Email{
			From: os.Getenv("EMAIL_FROM"),
			Host: os.Getenv("EMAIL_HOST"),
			Pass: os.Getenv("EMAIL_PASS"),
			Port: emailPort,
			User: os.Getenv("EMAIL_USER"),
		},
	}
}
