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

func GetUsersForCategory(categoryID uint) (users []User, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id = ?", categoryID); result.Error == nil {
		if result1 := Database.Model(&category).Association("UsersAccess"); result1.Error == nil {
			result2 := result1.Find(&category.UsersAccess)
			if result2 == nil {
				return category.UsersAccess, true
			}
		}
	}
	return users, false
}

func RemoveUsersFromCategory(userID uint, categoryID uint, users *[]User) (error string, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id=?", categoryID); result.Error != nil {
		return "could not find category", false
	}
	if category.OwnerID != userID {
		return "OwnerID != userID", false
	}
	if success := category.RemoveUsersFromCategoryInherit(users); success != true {
		return "Could not remove permissions from categories inherit", false
	}
	return "", true
}
func AddUsersForCategory(userID uint, categoryID uint, users *[]User, inherit bool) (error string, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id=?", categoryID); result.Error != nil {
		return "", false
	}
	if category.OwnerID != userID {
		return "OwnerID != userID", false
	}
	if category.ParentID == 0 {
		if inherit {
			if success := category.AddUsersToCategory(users); success != true {
				return "Could not add users to category", false
			}
		} else {
			if success := category.AddUsersToCategoryInherit(users); success != true {
				return "Could not add users to category and child categories", false
			}
		}
	} else {
		if exists, successret := UsersExistsInParentCategory(category.ParentID, users); successret == false {
			return "Database Failure", false
		} else {
			if exists == false {
				return "Not same users in parent category", false
			}
		}
		if inherit {
			if success := category.AddUsersToCategoryInherit(users); success != true {
				return "", false
			}
		} else {
			if success := category.AddUsersToCategory(users); success != true {
				return "", false
			}
		}
	}
	return "", true
}

func UsersExistsInParentCategory(id uint, users *[]User) (exists bool, success bool) {
	parentCategory := Category{}
	var usersFromParentCategory []User
	if result := Database.Find(&parentCategory, "id=?", id); result.Error != nil {
		return false, false
	}
	if result := Database.Model(&parentCategory).Association("UsersAccess").Find(&usersFromParentCategory); result != nil {
		return false, false
	}
	for i := range *users {
		var exists = false
		for j := range usersFromParentCategory {
			if (*users)[i].ID == usersFromParentCategory[j].ID {
				exists = true
			}
		}
		if !exists {
			return false, true
		}
	}
	return true, true
}

func GetAllUsers() (users []User, success bool) {
	if result := Database.Find(&users); result.Error == nil {
		return users, true
	}
	return users, false
}
