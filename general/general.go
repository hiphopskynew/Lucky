package general

import (
	"encoding/json"
	"fmt"
	"lucky/configs"
	"lucky/constants"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/sparkmaker/gohelper/validator/rule"

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

func InterfaceToString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

func GenerateJWTToken(payload interface{}) string {
	secret := configs.Setting.Jwt.Secret
	expired := configs.Setting.Jwt.Expired
	codec := jwt.ConfigWithExpire(secret, time.Duration(expired)*time.Second)
	return codec.Encode(payload)
}

func IsInvalidToken(req *http.Request) (msg string, is bool) {
	token := req.Header.Get("authorization")
	if token == "" {
		msg = "unauthorized access"
		is = false
		return
	}
	secret := configs.Setting.Jwt.Secret
	codec := jwt.Config(secret)
	if _, err := codec.Decode(token); err != nil {
		msg = "token is invalid"
		is = false
		return
	}
	is = true
	return
}

func aggregate(fs []rule.Failure) []rule.Failure {
	malformed := "malformed json"
	result := []rule.Failure{}
	unique := func(l []string) []string {
		m := make(map[string]interface{})
		r := []string{}
		for _, v := range l {
			m[v] = nil
		}
		for k := range m {
			r = append(r, k)
		}
		return r
	}
	find := func(e rule.Failure) int {
		for i, v := range result {
			if v.Key == e.Key {
				return i
			}
		}
		return -1
	}
	for _, f := range fs {
		if len(f.Key) == 0 && len(f.Messages) != 0 && f.Messages[0] == malformed {
			return []rule.Failure{f}
		}
		rIndex := find(f)
		if rIndex != -1 {
			result[rIndex].Messages = unique(append(result[rIndex].Messages, f.Messages...))
			continue
		} else {
			f.Messages = unique(f.Messages)
		}
		result = append(result, f)
	}
	return result
}

func MergeValidates(failures ...[]rule.Failure) []rule.Failure {
	fsr := []rule.Failure{}
	for _, fs := range failures {
		fsr = append(fsr, fs...)
	}
	return aggregate(fsr)
}

func JsonResponse(w http.ResponseWriter, result interface{}, status int) {
	bytes, _ := json.Marshal(result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}
