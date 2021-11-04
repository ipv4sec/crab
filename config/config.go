package config

type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
type Redis struct {
	Host     string
	Port     int
	Password string
}
type Config struct {
	Mysql Mysql
	Redis Redis
}