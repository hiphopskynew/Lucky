package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	profilemodels "lucky/services/profile/models"
	"lucky/services/repository/mysql"
	"lucky/services/user/validators"
	"net/http"

	"bitbucket.org/sparkmaker/gohelper/validator"
	"bitbucket.org/sparkmaker/gohelper/validator/rule"
	"github.com/gorilla/mux"
)

const (
	regexDOBFormat = `^([0-2][0-9]|(3)[0-1])(\/)(((0)[0-9])|((1)[0-2]))(\/)\d{4}$`
)

func validateCreateProfile(data string) []rule.Failure {
	rules := validator.New(data)
	rules.AddRule("first_name", rule.Required(), rule.IsString(), rule.NonEmpty())
	rules.AddRule("last_name", rule.Required(), rule.IsString(), rule.NonEmpty())
	rules.AddRule("date_of_birth", rule.Required(), rule.IsString(), rule.NonEmpty(), rule.Format(regexDOBFormat))
	rules.AddRule("address", rule.Required(), rule.IsString(), rule.NonEmpty())
	return general.MergeValidates(rules.Validate())
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
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
	bytes, _ := ioutil.ReadAll(r.Body)

	failures := validateCreateProfile(string(bytes))
	if len(failures) > 0 {
		general.JsonResponse(w, constants.M{constants.KeyError: failures}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)
	userProfile := profilemodels.UserProfile{ID: general.GenerateID(constants.PrefixProfile), FirstName: request.FirstName, LastName: request.LastName, DateOfBirth: request.DateOfBirth, Address: request.Address}
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
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: err.Error()}}, http.StatusInternalServerError)
		return
	}

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(userProfile)}, http.StatusCreated)
}
