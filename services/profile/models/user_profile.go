package models

type UserProfile struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Address     string `json:"address"`
	UserIDRef   string `json:"user_id"`
}
