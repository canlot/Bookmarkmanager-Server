package Models

import (
	"Bookmarkmanager-Server/Helpers"
	"errors"
)

func GetUser(username string) (user User, success bool) {
	if result := Database.Take(&user, "email = ?", username); result.Error == nil {
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
		return user, errors.New("user is not administrator")
	}
	var similarUser User
	if db := Database.Where("email = ?", user.Email).First(&similarUser); db.Error == nil {
		return User{}, errors.New("user with same email exists")
	}
	var hashedPassword string
	hashedPassword, err = Helpers.CreateHashFromPassword(password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword
	if db := Database.Create(&user); db.Error != nil {
		return user, db.Error
	}
	return user, nil
}
func EditUser(editorId uint, user User, password string) error {
	var err error
	var editor User
	if db := Database.Take(&editor, editorId); db.Error != nil {
		return db.Error
	}
	if editor.Administrator != true {
		return errors.New("user is not administrator")
	}
	if password != "" {
		var hashedPassword string
		hashedPassword, err = Helpers.CreateHashFromPassword(password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	if db := Database.Model(&user).Updates(&user); db.Error != nil {
		return db.Error
	}
	return nil
}
func SearchUsers(searchString string) ([]User, error) {
	searchString = "%" + searchString + "%"
	var users []User
	if db := Database.Where("email like ?", searchString).Or("name like ?", searchString).Find(&users); db.Error != nil {
		return nil, db.Error
	}
	return users, nil
}
