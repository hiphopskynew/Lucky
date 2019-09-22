package initialize

import (
	"io/ioutil"
	"lucky/configs"
	"lucky/general"
)

func Init() {
	// Read configuration file & initialized to the global variable
	bytes, err := ioutil.ReadFile("configs/application.json")
	if err != nil {
		panic(err)
	}
	setting := configs.ConfigurationModel{}
	general.ParseToStruct(bytes, &setting)
	configs.Setting = setting
}
