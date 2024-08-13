package config

import "time"

type Config struct {
	Server   Server
	Database Database
	Swagger  Swagger
	Kong     Kong
	Obs      ObsHuawei
	Jwt      Jwt
	Email    Email
}

type Server struct {
	Port string
}

type Database struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Swagger struct {
	Host string
	Url  string
	Mode string
}

type Kong struct {
	Url string
}

type ObsHuawei struct {
	Ak       string
	Sk       string
	Endpoint string
	Bucket   string
}

type Jwt struct {
	Secret string
	Expire time.Duration
	MaxAge int
}

type Email struct {
	From string
	Host string
	Pass string
	Port int
	User string
}
