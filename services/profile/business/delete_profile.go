package business

import (
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"net/http"

	"github.com/gorilla/mux"
)

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Deleted bool `json:"deleted"`
	}
	if message, isAuthen := general.IsInvalidToken(r); !isAuthen {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: message}}, http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]
	pid := mux.Vars(r)["pid"]

	session := mysql.New()
	defer session.Close()
	del, err := session.Exec("DELETE FROM UserProfile where id = ? AND user_id = ?", pid, id)
	if err != nil {
		panic(err)
	}

	resp := Response{}
	if c, err := del.RowsAffected(); err != nil {
		panic(err)
	} else if c > 0 {
		resp = Response{Deleted: true}
	} else {
		resp = Response{Deleted: false}
	}

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(resp)}, http.StatusOK)
}
