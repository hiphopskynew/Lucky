package business

import (
	"lucky/general"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Data string `json:"data"`
	}
	general.JsonResponse(w, Response{Data: "OK"})
}
