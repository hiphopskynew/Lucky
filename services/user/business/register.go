package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"lucky/services/user/validators"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type UserRegister struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Token     string    `json:"token"`
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)
	user := UserRegister{ID: general.GenerateID(constants.PrefixUser), Email: request.Email, Status: constants.StatusNew, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	session := mysql.New()
	defer session.Close()

	if bEmail := validators.HasEmailExist(session, user.Email); bEmail {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: "email already exist"}}, http.StatusBadRequest)
		return
	}

	if _, err := session.Query("INSERT INTO User(id, email, password, status, created_at, updated_at) VALUES(?,?,?,?,?,?)", user.ID, user.Email, request.Password, user.Status, user.CreatedAt, user.UpdatedAt); err != nil {
		panic(err)
	}
	token := general.GenerateToken()
	if _, err := session.Query("INSERT INTO UserVerify(id, email, token, created_at) VALUES(?,?,?,?)", general.GenerateID(constants.PrefixUserVerify), user.Email, token, time.Now()); err != nil {
		panic(err)
	}

	user.Token = token

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(user)}, http.StatusCreated)
}
