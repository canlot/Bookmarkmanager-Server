package Models

import (
	"gorm.io/gorm"
)

func GetUser(username string, showPassword bool) (user User, success bool) {
	ShowPassword = showPassword
	if result := Database.Take(&user, "name = ?", username); result.Error == nil {
		return user, true
	} else {
		return user, false
	}
}
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if ShowPassword == false {
		u.Password = ""
	}
	ShowPassword = false
	return nil
}

func GetAllUsers() (users []User, success bool) {
	if result := Database.Find(&users); result.Error == nil {
		return users, true
	}
	return users, false
}
func UsersExist(users []User) bool {
	for i := range users {
		if result := Database.Take(&(users[i])); result.Error != nil {
			return false
		}
	}
	return true
}
