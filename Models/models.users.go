package Models

import "gorm.io/gorm"

func GetUser(username string) (user User, success bool) {
	showPassword = true
	if result := Database.Take(&user, "name = ?", username); result.Error == nil {
		return user, true
	} else {
		return user, false
	}
}
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if showPassword == false {
		u.Password = ""
	}
	showPassword = false
	return nil
}

func GetUsersForCategoryFull(categoryID uint) (users []User, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id = ?", categoryID); result.Error == nil {
		if result1 := Database.Model(&category).Association("UsersFullAccess"); result1.Error == nil {
			result2 := result1.Find(&category.UsersFullAccess)
			if result2 == nil {
				return category.UsersFullAccess, true
			}
		}
	}
	return users, false
}
func GetUsersForCategoryInherit(categoryID uint) (users []User, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id = ?", categoryID); result.Error == nil {
		if result1 := Database.Model(&category).Association("UsersInheritAccess"); result1.Error == nil {
			result2 := result1.Find(&category.UsersFullAccess)
			if result2 == nil {
				return category.UsersFullAccess, true
			}
		}
	}
	return users, false
}
func AddUsersForCategory(userID uint, categoryID uint, users *[]User) (errorString string, success bool) {
	category := Category{}
	if result := Database.Find(&category, "id=?", categoryID); result.Error == nil {
		if category.OwnerID == userID {
			if error, success := category.SetPermissionForAll(users); success == true {
				return "", true
			} else {
				return error, false
			}
		} else {
			return "No permission to add a permission", false
		}
	} else {
		return "No such category", false
	}
}
