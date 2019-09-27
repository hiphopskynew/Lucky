package business

import (
	"lucky/constants"
	"lucky/general"
	profilemodels "lucky/services/profile/models"
	"lucky/services/repository/mysql"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProfileByUserID(w http.ResponseWriter, r *http.Request) {
	if message, isAuthen := general.IsInvalidToken(r); !isAuthen {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: message}}, http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id, first_name, last_name, date_of_birth, address, user_id FROM UserProfile WHERE user_id=?", id)
	if selErr != nil {
		panic(selErr)
	}
	userProfile := new(profilemodels.UserProfile)
	if !sel.Next() {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "profile not found"}}, http.StatusNotFound)
		return
	}
	sel.Scan(&userProfile.ID, &userProfile.FirstName, &userProfile.LastName, &userProfile.DateOfBirth, &userProfile.Address, &userProfile.UserIDRef)
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(userProfile)}, http.StatusOK)
}
