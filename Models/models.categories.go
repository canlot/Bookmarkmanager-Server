package Models

func GetCategories(UserId uint, categoryID uint) (categories []Category, success bool) {
	var user User
	if result := Database.Take(&user, UserId); result.Error == nil {
		if result := Database.Model(&user).Association("CategoriesFullAccess").Find(&categories, "parent_id = ?", categoryID); result == nil {
			var categoriesInherit []Category
			if result := Database.Model(&user).Association("CategoriesInheritAccess").Find(&categoriesInherit, "parent_id = ?", categoryID); result == nil {
				categories := append(categories, categoriesInherit...)
				return categories, true
			} else {
				return categories, false
			}
		} else {
			return categories, false
		}
	} else {
		return categories, false
	}
}

func AddCategory(UserId uint, category Category) (categoryret Category, success bool) {
	var user User
	if category.OwnerID == 0 || category.OwnerID != UserId {
		category.OwnerID = UserId
	}
	if result := Database.First(&user, UserId); result.Error == nil {
		if category.ParentID == 0 {
			if result := Database.Create(&category); result.Error == nil {
				if result := Database.Model(&category).Association("UsersFullAccess").Append(&user); result == nil {
					return category, true
				}
			}
		} else {
			var parentCategory Category
			if result := Database.Take(&parentCategory, category.ParentID); result.Error == nil {
				if parentCategory.OwnerID != category.OwnerID {
					return category, false
				}
				if result := Database.Create(&category); result.Error == nil {
					if result := Database.Model(&category).Association("UsersFullAccess").Append(&user); result == nil {
						return category, true
					}
				}
			}
		}

	}
	return category, false
}
