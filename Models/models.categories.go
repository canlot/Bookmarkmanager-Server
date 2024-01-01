package Models

import (
	"Bookmarkmanager-Server/Helpers"
	"errors"
	"gorm.io/gorm"
)

// get all categories for permitted user with parent_id
func GetCategories(UserId uint, categoryID uint) (categories []Category, err error) {
	var user User
	if result := Database.Take(&user, UserId); result.Error != nil {
		return nil, errors.New("Database error")
	}

	if result := Database.Model(&user).Association("CategoriesAccess").Find(&categories, "parent_id = ?", categoryID); result == nil {
		return categories, nil
	} else {
		return nil, errors.New(result.Error())
	}

}
func SearchCategories(UserId uint, searchText string) ([]Category, error) {
	var user User
	if db := Database.Take(&user, UserId); db.Error != nil {
		return nil, db.Error
	}
	var allCategories []Category
	if err := Database.Model(&user).Association("CategoriesAccess").Find(&allCategories); err != nil {
		return nil, err
	}
	var categories []Category
	for i := range allCategories {
		if Helpers.ContainsInAnyString(searchText, allCategories[i].Name, allCategories[i].Description) {
			categories = append(categories, allCategories[i])
		}
	}
	return categories, nil
}
func AddCategory(UserId uint, category Category) (Category, error) {
	var user User
	if category.OwnerID == 0 || category.OwnerID != UserId {
		category.OwnerID = UserId
	}
	if result := Database.First(&user, UserId); result.Error != nil {
		return category, result.Error
	}
	dbContext := Database.Begin()
	if category.ParentID == 0 {
		if result := dbContext.Create(&category); result.Error != nil {
			dbContext.Rollback()
			return category, result.Error
		}
		if err := dbContext.Model(&category).Association("UsersAccess").Append(&user); err != nil {
			dbContext.Rollback()
			return category, err
		}
		dbContext.Commit()
		return category, nil
	} else {
		var parentCategory Category
		if result := Database.Take(&parentCategory, category.ParentID); result.Error != nil {
			return category, result.Error
		}
		if parentCategory.OwnerID != category.OwnerID {
			return category, errors.New("Category owner is not the same as parent category owner")
		}
		category.Shared = parentCategory.Shared
		if result := dbContext.Create(&category); result.Error != nil {
			dbContext.Rollback()
			return category, result.Error
		}
		var parentCategoryUsers []User
		if err := dbContext.Model(&parentCategory).Association("UsersAccess").Find(&parentCategoryUsers); err != nil {
			dbContext.Rollback()
			return category, err
		}
		if err := dbContext.Model(&category).Association("UsersAccess").Append(&parentCategoryUsers); err != nil {
			dbContext.Rollback()
			return category, err
		}
		dbContext.Commit()
		return category, nil
	}
}

func EditCategory(userId uint, category Category) (Category, error) {
	if category.OwnerID != userId {
		return category, errors.New("Category OwnerID is not UserId")
	}
	var curCategory Category
	if result := Database.Take(&curCategory, category.ID); result.Error != nil {
		return category, errors.New("Database error or category not found")
	}
	if category.ParentID != curCategory.ParentID {
		return category, errors.New("ParentId cannot be changed")
	}
	dbContext := Database.Begin()
	if db := dbContext.Model(&category).Updates(&category); db.Error != nil {
		dbContext.Rollback()
		return category, db.Error
	}
	dbContext.Commit()
	return category, nil
}

func DeleteCategory(categoryId uint, UserId uint) bool {
	var category Category
	if result := Database.Take(&category, categoryId); result.Error != nil {
		return false
	}
	if category.OwnerID != UserId {
		return false
	}
	dbContext := Database.Begin()
	if category.DeleteAll(dbContext) != true {
		dbContext.Rollback()
		return false
	}
	dbContext.Commit()
	return true
}
func EditUsersForCategory(userID uint, categoryID uint, users *[]User) error {

	var currentUsers []User
	var success bool
	if currentUsers, success = GetUsersForCategory(categoryID); success != true {
		return errors.New("Could not get users for this category")
	}

	usersToAdd := getUsersThatAreOnlyOnLeftSide(users, &currentUsers)
	usersToRemove := getUsersThatAreOnlyOnLeftSide(&currentUsers, users)
	dbContext := Database.Begin()
	if len(usersToAdd) != 0 {
		if err := AddUsersForCategory(dbContext, userID, categoryID, &usersToAdd); err != nil {
			dbContext.Rollback()
			return err
		}
	}
	if len(usersToRemove) != 0 {
		if err := RemoveUsersFromCategory(dbContext, userID, categoryID, &usersToRemove); err != nil {
			dbContext.Rollback()
			return err
		}
	}
	dbContext.Commit()
	return nil
}

func getUsersThatAreOnlyOnLeftSide(usersLeft *[]User, usersRight *[]User) []User {
	var users []User
	for i := range *usersLeft {
		found := false
		for j := range *usersRight {
			if (*usersLeft)[i].ID == (*usersRight)[j].ID {
				found = true
				break
			}
		}
		if !found {
			users = append(users, (*usersLeft)[i])
		}
	}
	return users
}

func AddUsersForCategory(dbContext *gorm.DB, userID uint, categoryID uint, users *[]User) error {
	category := Category{}
	if result := dbContext.Find(&category, "id=?", categoryID); result.Error != nil {
		return result.Error
	}
	if category.OwnerID != userID {
		return errors.New("User != Owner")
	}
	if category.ParentID != 0 { //only to level categories can be shared
		return errors.New("Cannot grant permission because it is not a top level category")
	}
	if !UsersExist(dbContext, *users) {
		return errors.New("Not all users exist in database")
	}
	if success := category.AddUsersToCategoryInherit(dbContext, users); success != true {
		return errors.New("Could not add users to category and child categories")
	}
	return nil
}
func RemoveUsersFromCategory(dbContext *gorm.DB, userID uint, categoryID uint, users *[]User) error {
	category := Category{}
	if result := dbContext.Find(&category, "id=?", categoryID); result.Error != nil {
		return errors.New("could not find category")
	}
	if category.OwnerID != userID {
		return errors.New("OwnerID != userID")
	}
	if category.ParentID != 0 { //only to level categories can be shared
		return errors.New("Cannot remove permission because it is not a top level category")
	}
	if !UsersExist(dbContext, *users) {
		return errors.New("Not all users exist in database")
	}
	if success := category.RemoveUsersFromCategoryInherit(dbContext, users); success != true {
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
