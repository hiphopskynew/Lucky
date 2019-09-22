package configs

var Setting ConfigurationModel

type ConfigurationModel struct {
	Jwt struct {
		Secret  string `json:"secret"`
		Expired int    `json:"expired"`
	} `json:"jwt"`
	Repository struct {
		Mysql struct {
			Host        string `json:"host"`
			Port        string `json:"port"`
			Database    string `json:"database"`
			Credentials struct {
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"credentials"`
		} `json:"mysql"`
	} `json:"repository"`
}
