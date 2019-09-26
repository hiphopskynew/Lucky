package initialize

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"lucky/configs"
	"lucky/general"
	"lucky/services/repository/mysql"
)

func read(name string) string {
	dat, err := ioutil.ReadFile(fmt.Sprintf("initialize/sql/%s", name))
	if err != nil {
		log.Printf("Cannot reading `%s` file", name)
	}
	return string(dat)
}

func exec(s *sql.DB, sql string, tableName string) {
	pre, err := s.Prepare(sql)
	if err != nil {
		log.Printf("Cannot prepare create %s table caused : %s", tableName, err)
	}
	_, err = pre.Exec()
	if err != nil {
		log.Printf("Cannot create %s table caused : %s", tableName, err)
	} else {
		log.Printf("Create %s table successfully", tableName)
	}
}

func Init() {
	// Read configuration file & initialized to the global variable
	bytes, err := ioutil.ReadFile("configs/application.json")
	if err != nil {
		panic(err)
	}
	setting := configs.ConfigurationModel{}
	general.ParseToStruct(bytes, &setting)
	configs.Setting = setting

	tUser := read("create_user_table.sql")
	tUserVerify := read("create_user_verify_table.sql")
	tUserProfile := read("create_user_profile.sql")

	// SQL Init
	session := mysql.New()
	defer session.Close()
	exec(session, tUser, "USER")
	exec(session, tUserVerify, "USER VERIFY")
	exec(session, tUserProfile, "USER")
}
