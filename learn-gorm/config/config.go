package config

import (
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	DB ConfigDatabase
}

type ConfigDatabase struct {
	Endpoint string
	Username string
	Password string
	Database string
}

func (c ConfigDatabase) GetMysqlConfig() *mysql.Config {
	return &mysql.Config{
		User:   c.Username,
		Passwd: c.Password,
		Net:    "tcp",
		Addr:   c.Endpoint,
		DBName: c.Database,
		Params: map[string]string{
			"parseTime": "true",
		},
		AllowNativePasswords: true,
	}
}
