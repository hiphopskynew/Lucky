package business

import (
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetProfileByUserID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id, first_name, last_name, date_of_birth, address, user_id FROM UserProfile WHERE user_id=?", id)
	if selErr != nil {
		panic(selErr)
	}
	userProfile := new(usermodels.UserProfile)
	if !sel.Next() {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "profile not found"}}, http.StatusNotFound)
		return
	}
	sel.Scan(&userProfile.ID, &userProfile.FirstName, &userProfile.LastName, &userProfile.DateOfBirth, &userProfile.Address, &userProfile.UserIDRef)
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(userProfile)}, http.StatusOK)
}
