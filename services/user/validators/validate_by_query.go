package validators

import (
	"database/sql"
	"lucky/services/user/models"
)

func HasUserExist(session *sql.DB, id string) bool {
	userModelID := models.User{}
	err := session.QueryRow("SELECT id FROM User WHERE id=?", id).Scan(&userModelID.ID)
	if err != nil {
		return false
	}
	return true
}

func HasProfileExist(session *sql.DB, id string, pid string) bool {
	userProfileModelID := models.UserProfile{}
	err := session.QueryRow("SELECT id FROM UserProfile WHERE user_id=? AND id=?", id, pid).Scan(&userProfileModelID.ID)
	if err != nil {
		return false
	}
	return true
}

func HasProfileByUserIDExist(session *sql.DB, id string) bool {
	userProfileModelID := models.UserProfile{}
	err := session.QueryRow("SELECT id FROM UserProfile WHERE user_id=?", id).Scan(&userProfileModelID.ID)
	if err != nil {
		return false
	}
	return true
}
