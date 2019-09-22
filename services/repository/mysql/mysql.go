package mysql

import (
	"database/sql"
	"fmt"
	"lucky/configs"

	_ "github.com/go-sql-driver/mysql"
)

func New() *sql.DB {
	host := configs.Setting.Repository.Mysql.Host
	port := configs.Setting.Repository.Mysql.Port
	database := configs.Setting.Repository.Mysql.Database
	username := configs.Setting.Repository.Mysql.Credentials.Username
	password := configs.Setting.Repository.Mysql.Credentials.Password

	session, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database))
	if err != nil {
		panic(err)
	}

	return session
}
