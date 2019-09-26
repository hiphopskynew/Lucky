package routes

import (
	hcbiz "lucky/services/healthcheck/business"
	userbiz "lucky/services/user/business"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	// New Router
	router := mux.NewRouter()

	// Service Health Check
	router.HandleFunc("/ping", hcbiz.HealthCheck).Methods(http.MethodGet)

	// User Management
	router.HandleFunc("/api/users", userbiz.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}", userbiz.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/api/users/login", userbiz.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/users/register", userbiz.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/users/verify", userbiz.Verify).Methods(http.MethodPost)

	router.HandleFunc("/api/users/{id}/profile", userbiz.GetProfileByUserID).Methods(http.MethodGet)
	router.HandleFunc("/api/users/{id}/profile", userbiz.CreateProfile).Methods(http.MethodPost)
	router.HandleFunc("/api/users/{id}/profiles/{pid}", userbiz.UpdateProfile).Methods(http.MethodPut)
	router.HandleFunc("/api/users/{id}/profiles/{pid}", userbiz.DeleteProfile).Methods(http.MethodDelete)

	// Profile Management

	return router
}
