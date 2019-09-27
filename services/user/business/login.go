package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"lucky/services/user/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userModel := models.User{}
	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)
	session := mysql.New()
	defer session.Close()
	err := session.QueryRow("SELECT id, email, password FROM User WHERE email=? AND status=?", request.Email, constants.StatusVerified).Scan(&userModel.ID, &userModel.Email, &userModel.Password)
	if err != nil {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user not found"}}, http.StatusNotFound)
		return
	}
	if match := checkPasswordHash(request.Password, userModel.Password); !match {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "email or password is invalid"}}, http.StatusBadRequest)
		return
	}
	general.JsonResponse(w, constants.M{constants.KeyData: constants.M{constants.KeyToken: general.GenerateJWTToken(userModel.ID)}}, http.StatusOK)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
