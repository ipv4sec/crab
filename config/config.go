package config

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type Config struct {
	Mysql Mysql
}
