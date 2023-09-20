package Models

// Get all permitted categories for the current user if no id is provided, if id is provided, then get all child categories
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

func AddCategory(UserId uint, category Category) (Category, bool) {
	var user User
	if category.OwnerID == 0 || category.OwnerID != UserId {
		category.OwnerID = UserId
	}
	if result := Database.First(&user, UserId); result.Error == nil {
		if category.ParentID == 0 {
			if result := Database.Create(&category); result.Error == nil {
				if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
					return category, true
				}
			}
		} else {
			var parentCategory Category
			if result := Database.Take(&parentCategory, category.ParentID); result.Error == nil {
				if parentCategory.OwnerID == category.OwnerID {
					if result := Database.Create(&category); result.Error == nil {
						if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
							return category, true
						}
					}
				}
			}
		}

	}
	return Category{}, false
}

func EditCategory(userId uint, category Category) (Category, bool) {
	if category.OwnerID != userId {
		return Category{}, false
	}
	Database.Save(&category)
	return category, true
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
