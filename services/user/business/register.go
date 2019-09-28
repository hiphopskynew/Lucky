package business

import (
	"io/ioutil"
	"lucky/constants"
	"lucky/general"
	"lucky/services/repository/mysql"
	"lucky/services/user/validators"
	"net/http"
	"time"

	"bitbucket.org/sparkmaker/gohelper/validator"
	"bitbucket.org/sparkmaker/gohelper/validator/rule"
	"golang.org/x/crypto/bcrypt"
)

func validateRegister(data string) []rule.Failure {
	rules := validator.New(data)
	rules.AddRule("email", rule.Required(), rule.IsString(), rule.NonEmpty(), rule.Format(regexEmailFormat), rule.MaxLength(50))
	rules.AddRule("password", rule.Required(), rule.IsString(), rule.NonEmpty(), rule.MinLength(8), rule.MaxLength(140))
	return general.MergeValidates(rules.Validate())
}

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

	failures := validateRegister(string(bytes))
	if len(failures) > 0 {
		general.JsonResponse(w, constants.M{constants.KeyError: failures}, http.StatusBadRequest)
		return
	}

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

	pHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: err.Error()}}, http.StatusInternalServerError)
		return
	}

	if _, err := session.Query("INSERT INTO User(id, email, password, status, created_at, updated_at) VALUES(?,?,?,?,?,?)", user.ID, user.Email, string(pHash), user.Status, user.CreatedAt, user.UpdatedAt); err != nil {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: err.Error()}}, http.StatusInternalServerError)
		return
	}
	token := general.GenerateToken()
	if _, err := session.Query("INSERT INTO UserVerify(id, email, token, created_at) VALUES(?,?,?,?)", general.GenerateID(constants.PrefixUserVerify), user.Email, token, time.Now()); err != nil {
		general.JsonResponse(w, constants.M{constants.KeyError: constants.M{constants.KeyMessage: err.Error()}}, http.StatusInternalServerError)
		return
	}

	user.Token = token

	general.JsonResponse(w, constants.M{constants.KeyData: general.InterfaceToM(user)}, http.StatusCreated)
}
