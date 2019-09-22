package business

import (
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id, email, password, status, created_at, updated_at FROM User WHERE id=?", id)
	if selErr != nil {
		panic(selErr)
	}
	user := new(usermodels.User)
	if !sel.Next() {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "user not found"}})
		return
	}
	sel.Scan(&user.ID, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(user)})
}
