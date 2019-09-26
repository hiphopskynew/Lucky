package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"lucky/services/user/validators"
	"net/http"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token string `json:"token"`
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)

	session := mysql.New()
	defer session.Close()
	sel, selErr := session.Query("SELECT id, email, token, created_at FROM UserVerify WHERE token=?", request.Token)
	if selErr != nil {
		panic(selErr)
	}
	uv := new(usermodels.UserVerify)
	if !sel.Next() {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "token is invalid"}}, http.StatusBadRequest)
		return
	}

	sel.Scan(&uv.ID, &uv.Email, &uv.Token, &uv.CreatedAt)

	if bEmail := validators.IsStatusNew(session, uv.Email); !bEmail {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "email not found"}}, http.StatusNotFound)
		return
	}

	if _, err := session.Query("UPDATE User SET status=? WHERE email=?", constants.StatusVerified, uv.Email); err != nil {
		panic(err)
	}
	general.JsonResponse(w, constants.M{constants.KeyData: constants.M{constants.KeyEmail: uv.Email, constants.KeyStatus: constants.StatusVerified}}, http.StatusOK)
}
