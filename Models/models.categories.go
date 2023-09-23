package Models

import "errors"

// get all categories for permitted user with parent_id
func GetCategories(UserId uint, categoryID uint) (categories []Category, success bool) {
	var user User
	if result := Database.Take(&user, UserId); result.Error == nil {
		if result := Database.Model(&user).Association("CategoriesAccess").Find(&categories, "parent_id = ?", categoryID); result == nil {
			return categories, true
		} else {
			return []Category{}, false
		}
	} else {
		return []Category{}, false
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
