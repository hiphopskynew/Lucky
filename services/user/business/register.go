package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	usermodels "lucky/services/user/models"
	"net/http"
	"time"

	"bitbucket.org/sparkmaker/gohelper/logger/stdout"
)

func Register(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	bytes, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	request := new(Request)
	general.ParseToStruct(bytes, request)
	user := usermodels.User{ID: general.GenerateID(constants.PrefixUser), Email: request.Email, Password: request.Password, Status: constants.StatusNew, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	session := mysql.New()
	defer session.Close()
	if _, err := session.Query("INSERT INTO User(id, email, password, status, created_at, updated_at) VALUES(?,?,?,?,?,?)", user.ID, user.Email, user.Password, user.Status, user.CreatedAt, user.UpdatedAt); err != nil {
		panic(err)
	}
	token := general.GenerateToken()
	if _, err := session.Query("INSERT INTO UserVerify(id, email, token, created_at) VALUES(?,?,?,?)", general.GenerateID(constants.PrefixUserVerify), user.Email, token, time.Now()); err != nil {
		panic(err)
	}
	stdout.Info("verify token: ", token)
	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(user)}, http.StatusCreated)
}
