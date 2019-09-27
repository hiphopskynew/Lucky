package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	profilemodels "lucky/services/profile/models"
	"lucky/services/repository/mysql"
	"lucky/services/user/validators"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
		Address     string `json:"address"`
	}
	if message, isAuthen := general.IsInvalidToken(r); !isAuthen {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: message}}, http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]
	pid := mux.Vars(r)["pid"]
	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)

	if check, _ := regexp.MatchString(`^([0-2][0-9]|(3)[0-1])(\/)(((0)[0-9])|((1)[0-2]))(\/)\d{4}$`, request.DateOfBirth); !check {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "date of birth does not match"}}, http.StatusBadRequest)
		return
	}

	userProfile := profilemodels.UserProfile{ID: pid, FirstName: request.FirstName, LastName: request.LastName, DateOfBirth: request.DateOfBirth, Address: request.Address, UserIDRef: id}
	session := mysql.New()
	defer session.Close()

	if bUser := validators.HasUserExist(session, id); !bUser {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user does not exist"}}, http.StatusBadRequest)
		return
	}

	if bProfile := validators.HasProfileExist(session, id, pid); !bProfile {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "profile does not exist"}}, http.StatusBadRequest)
		return
	}

	if _, err := session.Query("UPDATE UserProfile SET first_name=?, last_name=?, date_of_birth=?, address=? WHERE user_id=? AND id=?", userProfile.FirstName, userProfile.LastName, userProfile.DateOfBirth, userProfile.Address, id, pid); err != nil {
		panic(err)
	}

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(userProfile)}, http.StatusOK)
}
