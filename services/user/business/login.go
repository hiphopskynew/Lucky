package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"lucky/services/user/models"
	"net/http"

	"bitbucket.org/sparkmaker/gohelper/validator"
	"bitbucket.org/sparkmaker/gohelper/validator/rule"
	"golang.org/x/crypto/bcrypt"
)

const (
	regexEmailFormat = "^(?:[a-z0-9+_-]+(?:\\.[a-z0-9+_-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\]).$"
)

func validateLogin(data string) []rule.Failure {
	rules := validator.New(data)
	rules.AddRule("email", rule.Required(), rule.IsString(), rule.NonEmpty(), rule.Format(regexEmailFormat), rule.MaxLength(50))
	rules.AddRule("password", rule.Required(), rule.IsString(), rule.NonEmpty(), rule.MinLength(8), rule.MaxLength(140))
	return general.MergeValidates(rules.Validate())
}

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userModel := models.User{}
	bytes, _ := ioutil.ReadAll(r.Body)

	failures := validateLogin(string(bytes))
	if len(failures) > 0 {
		general.JsonResponse(w, constants.M{constants.KeyError: failures}, http.StatusBadRequest)
		return
	}

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
