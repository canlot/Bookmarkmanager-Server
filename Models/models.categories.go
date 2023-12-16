package Models

import "errors"

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
	searchText = "%" + searchText + "%"
	var user User
	if db := Database.Take(&user, UserId); db.Error != nil {
		return nil, db.Error
	}
	Database.Model(&user).Association("CategoriesAccess")
	categories := Database.Where("owner_id = ?", UserId).
		Where(Database.Where("name LIKE ?", searchText).Or("description LIKE ?", searchText))
}
func AddCategory(UserId uint, category Category) (Category, error) {
	var user User
	if category.OwnerID == 0 || category.OwnerID != UserId {
		category.OwnerID = UserId
	}
	if result := Database.First(&user, UserId); result.Error == nil {
		if category.ParentID == 0 {
			if result := Database.Create(&category); result.Error == nil {
				if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
					return category, nil
				}
			}
		} else {
			var parentCategory Category
			if result := Database.Take(&parentCategory, category.ParentID); result.Error == nil {
				if parentCategory.OwnerID == category.OwnerID {
					if result := Database.Create(&category); result.Error == nil {
						if result := Database.Model(&category).Association("UsersAccess").Append(&user); result == nil {
							return category, nil
						}
					}
				}
			}
		}

	}
	return category, errors.New("Error occured")
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
	Database.Save(&category)
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
	if category.DeleteAll() != true {
		return false
	}
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

	if len(usersToAdd) != 0 {
		if err := AddUsersForCategory(userID, categoryID, &usersToAdd); err != nil {
			return err
		}
	}
	if len(usersToRemove) != 0 {
		if err := RemoveUsersFromCategory(userID, categoryID, &usersToRemove); err != nil {
			return err
		}
	}
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
	if !UsersExist(*users) {
		return errors.New("Not all users exist in database")
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
	if !UsersExist(*users) {
		return errors.New("Not all users exist in database")
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
