package Models

import "fmt"

func (category Category) RemoveUsersFromCategoryInherit(users *[]User) (success bool) {
	if !category.RemoveUsersFromCategory(users) {
		return true
	}
	var childCategories []Category
	if result := Database.Find(&childCategories, "parent_id = ?", category.ID); result.Error != nil {
		return false
	}
	for i := range childCategories {
		if childCategories[i].RemoveUsersFromCategoryInherit(users) != true {
			return false
		}
	}
	return true
}
func (category Category) AddUsersToCategoryInherit(users *[]User) (success bool) {
	if !category.AddUsersToCategory(users) {
		return false
	}
	var childCategories []Category
	if result := Database.Find(&childCategories, "parent_id = ?", category.ParentID); result.Error != nil {
		return false
	}
	for i := range childCategories {
		if childCategories[i].AddUsersToCategoryInherit(users) != true {
			return false
		}
	}
	return true
}
func (category Category) RemoveUsersFromCategory(users *[]User) (success bool) {
	var uniqueUsers []User
	if uniqueUsers, success = category.GetOnlyExistingUsersInCurrentCategory(users); success != true {
		return false
	}
	if uniqueUsers, success = category.GetUsersWithoutOwner(&uniqueUsers); success != true {
		return false
	}
	if !category.AreUsersUnique(&uniqueUsers) {
		return false
	}
	if len(uniqueUsers) == 0 {
		return false
	}
	if result := Database.Model(&category).Association("UsersAccess").Delete(uniqueUsers); result != nil {
		return false
	}
	if leftusers, success := category.GetUsersForCategory(); success != true {
		return false
	} else {
		if len(leftusers) == 1 {
			if result := Database.Model(&category).Update("shared", "0"); result.Error != nil {
				return false
			}
		}
	}
	return true
}
func (category Category) AddUsersToCategory(users *[]User) (success bool) {
	var uniqueUsers []User
	if uniqueUsers, success = category.GetOnlyNonExistingUsersInCurrentCategory(users); success != true {
		return false
	}
	if !category.AreUsersUnique(&uniqueUsers) {
		return false
	}
	if uniqueUsers, success = category.GetUsersWithoutOwner(&uniqueUsers); success != true {
		return false
	}
	if len(uniqueUsers) == 0 {
		return false
	}
	if result := Database.Model(&category).Association("UsersAccess").Append(uniqueUsers); result != nil {
		return false
	}
	if result := Database.Model(&category).Update("shared", "1"); result.Error != nil {
		return false
	}
	return true
}
func (category Category) AreUsersUnique(users *[]User) (success bool) {
	for i := range *users {
		for j := range *users {
			if i != j {
				if (*users)[i].ID == (*users)[j].ID {
					return false
				}
			}
		}
	}
	return true
}
func (category Category) GetUsersWithoutOwner(users *[]User) (newusers []User, success bool) {
	for i := range *users {
		if (*users)[i].ID == category.OwnerID {

		} else {
			newusers = append(newusers, (*users)[i])
		}
	}
	return newusers, true
}
func (category Category) GetOnlyExistingUsersInCurrentCategory(users *[]User) (differentUsers []User, success bool) {
	if originalUsers, success := category.GetUsersForCategory(); success == true {
		for i := range *users {
			var exists = false
			for j := range originalUsers {
				if (*users)[i].ID == originalUsers[j].ID {
					exists = true
				}
			}
			if exists {
				differentUsers = append(differentUsers, (*users)[i])
			}
		}
	} else {
		return differentUsers, false
	}
	return differentUsers, true
}
func (category Category) GetOnlyNonExistingUsersInCurrentCategory(users *[]User) (differentUsers []User, success bool) {
	if originalUsers, success := category.GetUsersForCategory(); success == true {
		for i := 0; i < len(*users); i++ {
			var exists = false
			for j := 0; j < len(originalUsers); j++ {
				if (*users)[i].ID == originalUsers[j].ID {
					exists = true
				}
			}
			if exists == false {
				differentUsers = append(differentUsers, (*users)[i])
			}
		}
	} else {
		return differentUsers, false
	}
	return differentUsers, true
}
func (category Category) GetUsersForCategory() (users []User, success bool) {
	if result := Database.Model(&category).Association("UsersAccess").Find(&users); result != nil {
		return users, false
	}
	return users, true
}

func (category Category) DeleteAll() bool {

	var categories []Category
	if result := Database.Find(&categories, "parent_id = ?", category.ID); result.Error != nil {
		return false
	}
	for i := range categories {
		if categories[i].DeleteAll() != true {
			return false
		}
	}
	if result := Database.Model(&category).Association("UsersAccess").Clear(); result != nil {
		return false
	}
	if result := Database.Delete(Bookmark{}, "category_id = ?", category.ID); result.Error != nil {
		return false
	}
	if result := Database.Delete(&category); result.Error != nil {
		fmt.Println(result.Error)
		return false
	}
	return true
}
