package Models

import "errors"

// get all categories for permitted user with parent_id
func GetCategories(UserId uint, categoryID uint) (categories []Category, err error) {
	var user User
	if result := Database.Take(&user, UserId); result.Error != nil {
		return []Category{}, errors.New("Database error")
	}

	if result := Database.Model(&user).Association("CategoriesAccess").Find(&categories, "parent_id = ?", categoryID); result == nil {
		return categories, nil
	} else {
		return []Category{}, errors.New(result.Error())
	}

}

func AddCategory(UserId uint, category Category) error {
	var user User
	if category.OwnerID == 0 || category.OwnerID != UserId {
		category.OwnerID = UserId
	}
	if result := Database.First(&user, UserId); result.Error == nil {
		if category.ParentID == 0 {
			if result := Database.Create(&category); result.Error == nil {
				if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
					return nil
				}
			}
		} else {
			var parentCategory Category
			if result := Database.Take(&parentCategory, category.ParentID); result.Error == nil {
				if parentCategory.OwnerID == category.OwnerID {
					if result := Database.Create(&category); result.Error == nil {
						if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
							return nil
						}
					}
				}
			}
		}

	}
	return errors.New("Error occured")
}

func EditCategory(userId uint, category Category) error {
	if category.OwnerID != userId {
		return errors.New("Category OwnerID is not UserId")
	}
	var curCategory Category
	if result := Database.Take(&curCategory, category.ID); result.Error != nil {
		return errors.New("Database error or category not found")
	}
	if category.ParentID != curCategory.ParentID {
		return errors.New("ParentId cannot be changed")
	}
	Database.Save(&category)
	return nil
}

func DeleteCategory(categoryId uint, UserId uint) bool {
	var category Category
	if result := Database.Take(&category, categoryId); result.Error != nil {
		return false
	}
	if category.OwnerID != UserId {
		return false
	}
	if category.DeleteAll() != true {
		return false
	}
	return true
}

func AddUsersForCategory(userID uint, categoryID uint, users *[]User) error {
	category := Category{}
	if result := Database.Find(&category, "id=?", categoryID); result.Error != nil {
		return result.Error
	}
	if category.OwnerID != userID {
		return errors.New("User != Owner")
	}
	if category.ParentID != 0 { //only to level categories can be shared
		return errors.New("Cannot grant permission because it is not a top level category")
	}

	if success := category.AddUsersToCategoryInherit(users); success != true {
		return errors.New("Could not add users to category and child categories")
	}

	return nil
}
func RemoveUsersFromCategory(userID uint, categoryID uint, users *[]User) error {
	category := Category{}
	if result := Database.Find(&category, "id=?", categoryID); result.Error != nil {
		return errors.New("could not find category")
	}
	if category.OwnerID != userID {
		return errors.New("OwnerID != userID")
	}
	if category.ParentID != 0 { //only to level categories can be shared
		return errors.New("Cannot remove permission because it is not a top level category")
	}
	if success := category.RemoveUsersFromCategoryInherit(users); success != true {
		return errors.New("Could not remove permissions from categories inherit")
	}
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
