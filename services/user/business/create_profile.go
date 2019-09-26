package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"lucky/services/user/validators"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DateOfBirth string `json:"date_of_birth"`
		Address     string `json:"address"`
	}
	id := mux.Vars(r)["id"]
	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)

	if check, _ := regexp.MatchString(`^([0-2][0-9]|(3)[0-1])(\/)(((0)[0-9])|((1)[0-2]))(\/)\d{4}$`, request.DateOfBirth); !check {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "date of birth does not match"}}, http.StatusBadRequest)
		return
	}

	general.ParseToStruct(bytes, request)
	userProfile := usermodels.UserProfile{ID: general.GenerateID(constants.PrefixProfile), FirstName: request.FirstName, LastName: request.LastName, DateOfBirth: request.DateOfBirth, Address: request.Address}
	session := mysql.New()
	defer session.Close()

	if bUser := validators.HasUserExist(session, id); !bUser {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user does not exist"}}, http.StatusBadRequest)
		return
	}

	if bProfile := validators.HasProfileByUserIDExist(session, id); bProfile {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "profile is already exist"}}, http.StatusBadRequest)
		return
	}

	if _, err := session.Query("INSERT INTO UserProfile(id, first_name, last_name, date_of_birth, address, user_id) VALUES(?,?,?,?,?,?)", userProfile.ID, userProfile.FirstName, userProfile.LastName, userProfile.DateOfBirth, userProfile.Address, id); err != nil {
		panic(err)
	}

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(userProfile)}, http.StatusCreated)
}
