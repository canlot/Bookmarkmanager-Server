package Models

import (
	"gorm.io/gorm"
)

func (category Category) RemoveUsersFromCategoryInherit(dbContext *gorm.DB, users *[]User) (success bool) {
	if !category.RemoveUsersFromCategory(dbContext, users) {
		return true
	}
	var childCategories []Category
	if result := dbContext.Find(&childCategories, "parent_id = ?", category.ID); result.Error != nil {
		return false
	}
	for i := range childCategories {
		if childCategories[i].RemoveUsersFromCategoryInherit(dbContext, users) != true {
			return false
		}
	}
	return true
}
func (category Category) AddUsersToCategoryInherit(dbContext *gorm.DB, users *[]User) (success bool) {
	if !category.AddUsersToCategory(dbContext, users) {
		return false
	}
	var childCategories []Category
	if result := dbContext.Find(&childCategories, "parent_id = ?", category.ID); result.Error != nil {
		return false
	}
	for i := range childCategories {
		if childCategories[i].AddUsersToCategoryInherit(dbContext, users) != true {
			return false
		}
	}
	return true
}
func (category Category) RemoveUsersFromCategory(dbContext *gorm.DB, users *[]User) (success bool) {
	var uniqueUsers []User
	if uniqueUsers, success = category.GetOnlyExistingUsersInCurrentCategory(dbContext, users); success != true {
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
	if result := dbContext.Model(&category).Association("UsersAccess").Delete(&uniqueUsers); result != nil {
		return false
	}
	if leftusers, success := category.GetUsersForCategory(dbContext); success != true {
		return false
	} else {
		if len(leftusers) == 1 {
			if result := dbContext.Model(&category).Update("shared", "0"); result.Error != nil {
				return false
			}
		}
	}
	return true
}
func (category Category) AddUsersToCategory(dbContext *gorm.DB, users *[]User) (success bool) {
	var uniqueUsers []User
	if uniqueUsers, success = category.GetOnlyNonExistingUsersInCurrentCategory(dbContext, users); success != true {
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
	if result := dbContext.Model(&category).Omit("UsersAccess").Association("UsersAccess").Append(&uniqueUsers); result != nil {
		return false
	}
	if result := dbContext.Model(&category).Update("shared", "1"); result.Error != nil {
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
func (category Category) GetOnlyExistingUsersInCurrentCategory(dbContext *gorm.DB, users *[]User) (differentUsers []User, success bool) {
	if originalUsers, success := category.GetUsersForCategory(dbContext); success == true {
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
func (category Category) GetOnlyNonExistingUsersInCurrentCategory(dbContext *gorm.DB, users *[]User) (differentUsers []User, success bool) {
	if originalUsers, success := category.GetUsersForCategory(dbContext); success == true {
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
func (category Category) GetUsersForCategory(dbContext *gorm.DB) (users []User, success bool) {
	if result := dbContext.Model(&category).Association("UsersAccess").Find(&users); result != nil {
		return users, false
	}
	return users, true
}

func (category Category) DeleteAll(dbContext *gorm.DB) bool {
	var categories []Category
	if result := Database.Find(&categories, "parent_id = ?", category.ID); result.Error != nil {
		return false
	}
	for i := range categories {
		if categories[i].DeleteAll(dbContext) != true {
			return false
		}
	}
	if result := dbContext.Model(&category).Association("UsersAccess").Clear(); result != nil {
		return false
	}
	if result := dbContext.Delete(&Bookmark{}, "category_id = ?", category.ID); result.Error != nil {
		return false
	}
	if result := dbContext.Delete(&category); result.Error != nil {
		return false
	}
	return true
}
