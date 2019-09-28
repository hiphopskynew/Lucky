package business

import (
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"net/http"
	"time"
)

func GetUsersProfiles(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		ID          string    `json:"id"`
		Email       string    `json:"email"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		ProfileID   string    `json:"profile_id"`
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		DateOfBirth string    `json:"date_of_birth"`
		Address     string    `json:"address"`
	}
	if message, isAuthen := general.IsInvalidToken(r); !isAuthen {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: message}}, http.StatusUnauthorized)
		return
	}
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT User.id, User.email, User.status, User.created_at, User.updated_at, UserProfile.id, UserProfile.first_name, UserProfile.last_name, UserProfile.date_of_birth, UserProfile.address FROM User LEFT JOIN UserProfile ON User.id = UserProfile.user_id")
	if selErr != nil {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: selErr.Error()}}, http.StatusInternalServerError)
		return
	}
	usersProfiles := []Response{}
	for sel.Next() {
		user := Response{}
		sel.Scan(&user.ID, &user.Email, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.ProfileID, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Address)
		usersProfiles = append(usersProfiles, user)
	}
	if len(usersProfiles) == 0 {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user profile does not exist"}}, http.StatusNotFound)
		return
	}
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToSliceM(usersProfiles)}, http.StatusOK)
}
