package general

import (
	"encoding/json"
	"fmt"
	"lucky/configs"
	"lucky/constants"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/sparkmaker/gohelper/jwt"
	"github.com/google/uuid"
)

func GenerateID(prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, uuid.New().String())
}

func GenerateToken() string {
	token := ""
	for i := 0; i < 5; i++ {
		token += uuid.New().String()
	}
	return strings.Replace(token, "-", "", -1)
}

func ParseToStruct(bytes []byte, model interface{}) {
	json.Unmarshal(bytes, &model)
}

func InterfaceToSliceM(v interface{}) []constants.M {
	ms := []constants.M{}
	bytes, _ := json.Marshal(v)
	json.Unmarshal(bytes, &ms)
	return ms
}

func InterfaceToM(v interface{}) constants.M {
	m := constants.M{}
	bytes, _ := json.Marshal(v)
	json.Unmarshal(bytes, &m)
	return m
}

func GenerateJWTToken(payload interface{}) string {
	secret := configs.Setting.Jwt.Secret
	expired := configs.Setting.Jwt.Expired
	codec := jwt.ConfigWithExpire(secret, time.Duration(expired)*time.Second)
	return codec.Encode(payload)
}

func JsonResponse(w http.ResponseWriter, result interface{}, status int) {
	bytes, _ := json.Marshal(result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}
