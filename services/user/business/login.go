package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id FROM User WHERE email=? AND password=?", request.Email, request.Password)
	if selErr != nil {
		panic(selErr)
	}
	if !sel.Next() {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "email or password is invalid"}})
		return
	}
	var id string
	sel.Scan(&id)
	general.JsonResponse(w, constants.M{constants.KeyData: constants.M{constants.KeyToken: general.GenerateJWTToken(id)}})
}
