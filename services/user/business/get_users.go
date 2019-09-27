package business

import (
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if message, isAuthen := general.IsInvalidToken(r); !isAuthen {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: message}}, http.StatusUnauthorized)
		return
	}
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id, email, password, status, created_at, updated_at FROM User")
	if selErr != nil {
		panic(selErr)
	}
	users := []usermodels.User{}
	for sel.Next() {
		user := usermodels.User{}
		sel.Scan(&user.ID, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}
	if len(users) == 0 {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user does not exist"}}, http.StatusNotFound)
		return
	}
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToSliceM(users)}, http.StatusOK)
}
