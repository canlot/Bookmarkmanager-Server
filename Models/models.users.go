package Models

import (
	"Bookmarkmanager-Server/Helpers"
	"errors"
)

func GetUser(username string) (user User, success bool) {
	if result := Database.Take(&user, "name = ?", username); result.Error == nil {
		return user, true
	} else {
		return user, false
	}
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
func AddUser(creatorId uint, user User, password string) (User, error) {
	var err error
	var owner User
	if db := Database.Take(&owner, creatorId); db.Error != nil {
		return user, db.Error
	}
	if owner.Administrator != true {
		return user, errors.New("User is not administrator")
	}
	var hashedPassword string
	hashedPassword, err = Helpers.CreateHashFromPassword(password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword
	if db := Database.Create(&user); db.Error != nil {
		return user, nil
	}
	return user, nil
}
