package Models

func (category Category) SetPermissionForAll(users *[]User) (errortext string, success bool) {
	if result := Database.Model(&category).Association("UsersFullAccess").Append(users); result == nil {
		Database.Model(&category).Update("shared", "1")
		if error, success := category.SetPermissionForParent(users); success == true {
			if error, success := category.SetPermissionForChildren(users); success == true {
				return "", true
			} else {
				return error, false
			}
		} else {
			return error, false
		}

	} else {
		return result.Error(), false
	}
}
func (category Category) SetPermissionForParent(users *[]User) (errortext string, success bool) {
	if category.ParentID != 0 {
		var categoryParent Category
		if result := Database.First(&categoryParent, category.ParentID); result.Error == nil {
			if result := Database.Model(&categoryParent).Association("UsersInheritAccess").Append(users); result == nil {
				if categoryParent.Shared != true {
					Database.Model(&categoryParent).Update("shared", "1")
				}
				if error, success := categoryParent.SetPermissionForParent(users); success == true {
					return "", true
				} else {
					return error, false
				}
			} else {
				return "Could not set permission for category", false
			}
		} else {
			return "Could not find parent category", false
		}
	} else {
		return "", true
	}
}
func (category Category) SetPermissionForChildren(users *[]User) (errortext string, success bool) {
	var categoryChildren []Category
	var bookmarks []Bookmark
	if result := Database.Model(&category).Association("Bookmarks").Find(&bookmarks); result == nil {
		if len(bookmarks) != 0 {
			if result := Database.Model(&bookmarks).Association("UsersFullAccess").Append(users); result != nil {
				return "Could not add permissions to bookmarks", false
			}
		}
	} else {
		return "Could not find any bookmarks", false
	}
	if result := Database.Where("parent_id", category.ID).Find(&categoryChildren); result.Error == nil {
		if len(categoryChildren) != 0 {
			if result := Database.Model(&categoryChildren).Association("UsersFullAccess").Append(users); result == nil {
				for i := 0; i < len(categoryChildren); i++ {
					if err, success := categoryChildren[i].SetPermissionForChildren(users); success != true {
						return err, false
					}
				}
				return "", true
			} else {
				return "Could not add permissions to categories", false
			}
		} else {
			return "", true
		}
	} else {
		return "could not get categories", false
	}
}
