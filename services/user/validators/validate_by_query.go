package validators

import (
	"database/sql"
	"lucky/constants"
	profmodels "lucky/services/profile/models"
	usermodels "lucky/services/user/models"
)

func HasUserExist(session *sql.DB, id string) bool {
	userModelID := usermodels.User{}
	err := session.QueryRow("SELECT id FROM User WHERE id=?", id).Scan(&userModelID.ID)
	if err != nil {
		return false
	}
	return true
}

func HasProfileExist(session *sql.DB, id string, pid string) bool {
	userProfileModelID := profmodels.UserProfile{}
	err := session.QueryRow("SELECT id FROM UserProfile WHERE user_id=? AND id=?", id, pid).Scan(&userProfileModelID.ID)
	if err != nil {
		return false
	}
	return true
}

func HasProfileByUserIDExist(session *sql.DB, id string) bool {
	userProfileModelID := profmodels.UserProfile{}
	err := session.QueryRow("SELECT id FROM UserProfile WHERE user_id=?", id).Scan(&userProfileModelID.ID)
	if err != nil {
		return false
	}
	return true
}

func HasEmailExist(session *sql.DB, email string) bool {
	user := usermodels.User{}
	err := session.QueryRow("SELECT id FROM User WHERE email=?", email).Scan(&user.ID)
	if err != nil {
		return false
	}
	return true
}

func IsStatusNew(session *sql.DB, email string) bool {
	user := usermodels.User{}
	err := session.QueryRow("SELECT id FROM User WHERE email=? AND status=?", email, constants.StatusNew).Scan(&user.ID)
	if err != nil {
		return false
	}
	return true
}
