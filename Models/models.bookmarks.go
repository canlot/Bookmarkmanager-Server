package Models

import (
	"Bookmarkmanager-Server/Helpers"
	"errors"
)

func GetBookmarksWithCategoryID(userId uint, categoryId uint) (bookmarks []Bookmark, success bool) {
	var user User
	var category Category

	if result := Database.Take(&user, userId); result.Error != nil {
		return nil, false
	}

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return nil, false
	}

	var users []User
	if result := Database.Model(&category).Association("UsersAccess").Find(&users); result != nil {
		return nil, false
	}

	var userExist = false
	for i := range users {
		if users[i].ID == user.ID {
			userExist = true
			break
		}
	}

	if userExist != true {
		return nil, false
	}

	if result := Database.Find(&bookmarks, "category_id = ?", categoryId); result.Error != nil {
		return nil, false
	}
	return bookmarks, true
}

func SearchBookmarks(UserId uint, searchText string) ([]Bookmark, error) {
	var user User
	if db := Database.Take(&user, UserId); db.Error != nil {
		return nil, db.Error
	}
	var allCategories []Category
	if err := Database.Model(&user).Association("CategoriesAccess").Find(&allCategories); err != nil {
		return nil, err
	}
	var bookmarks []Bookmark
	for i := range allCategories {
		var categoryBookmarks []Bookmark
		if result := Database.Find(&categoryBookmarks, "category_id = ?", allCategories[i].ID); result.Error != nil {
			return nil, result.Error
		}
		for j := range categoryBookmarks {
			if Helpers.ContainsInAnyString(searchText, categoryBookmarks[j].Title, categoryBookmarks[j].Url, categoryBookmarks[j].Description) {
				bookmarks = append(bookmarks, categoryBookmarks[j])
			}
		}
	}
	return bookmarks, nil
}
func AddBookmark(userId uint, categoryId uint, bookmark Bookmark) (Bookmark, bool) {
	var category Category

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return bookmark, false
	}
	if category.OwnerID != userId {
		return bookmark, false
	}
	bookmark.CategoryID = category.ID
	dbContext := Database.Begin()
	if result := dbContext.Create(&bookmark); result.Error != nil {
		dbContext.Rollback()
		return bookmark, false
	}
	dbContext.Commit()
	return bookmark, true
}
func EditBookmark(userId uint, categoryId uint, bookmark Bookmark) error {
	var category Category
	var existingBookmark Bookmark

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return result.Error
	}
	if result := Database.Take(&existingBookmark, bookmark.ID); result.Error != nil {
		return result.Error
	}
	if existingBookmark.CategoryID != categoryId || bookmark.CategoryID != existingBookmark.CategoryID {
		return errors.New("Category id of bookmark is not the same as given")
	}
	if existingBookmark.CreatedAt.IsZero() { // a wild hack to fix createdat zero value in db
		bookmark.CreatedAt = existingBookmark.UpdatedAt
	}
	if category.OwnerID != userId {
		return errors.New("owner of bookmark not the same as logged on user")
	}
	dbContext := Database.Begin()
	if result := dbContext.Save(&bookmark); result.Error != nil {
		dbContext.Rollback()
		return result.Error
	}
	dbContext.Commit()
	return nil
}
func DeleteBookmark(userId uint, categoryId uint, bookmarkId uint) error {
	var category Category
	var bookmark Bookmark

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return result.Error
	}
	if result := Database.Take(&bookmark, bookmarkId); result.Error != nil {
		return result.Error
	}
	if bookmark.CategoryID != categoryId {
		return errors.New("Category id of bookmark is not the same as given")
	}
	if category.OwnerID != userId {
		return errors.New("Owner of bookmark not the same as loged on user")
	}
	dbContext := Database.Begin()
	if result := dbContext.Delete(&Bookmark{}, bookmarkId); result.Error != nil {
		dbContext.Rollback()
		return result.Error
	}
	dbContext.Commit()
	return nil
}
func GetBookmark(userId uint, categoryId uint, bookmarkId uint) (Bookmark, error) {
	var category Category
	var bookmark Bookmark

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return bookmark, result.Error
	}
	if result := Database.Take(&bookmark, bookmarkId); result.Error != nil {
		return bookmark, result.Error
	}
	if bookmark.CategoryID != categoryId {
		return bookmark, errors.New("Category id of bookmark is not the same as given")
	}
	if category.OwnerID != userId {
		return bookmark, errors.New("Owner of bookmark not the same as loged on user")
	}
	return bookmark, nil

}
func MoveBookmarkToCategory(userId, bookmarkId, sourceCategoryId, destinationCategoryId uint) error {
	var sourceCategory, destinationCategory Category
	var bookmark Bookmark
	if result := Database.Take(&sourceCategory, sourceCategoryId); result.Error != nil {
		return result.Error
	}
	if result := Database.Take(&destinationCategory, destinationCategoryId); result.Error != nil {
		return result.Error
	}
	if sourceCategory.OwnerID != userId {
		return errors.New("Owner of source category not the same as logged on user")
	}
	if destinationCategory.OwnerID != userId {
		return errors.New("Owner of destination category not the same as logged on user")
	}
	if result := Database.Take(&bookmark, bookmarkId); result.Error != nil {
		return result.Error
	}
	if bookmark.CategoryID != sourceCategoryId {
		errors.New("Category id of bookmark is not the same as given")
	}
	bookmark.CategoryID = destinationCategoryId

	dbContext := Database.Begin()
	if result := dbContext.Save(&bookmark); result.Error != nil {
		dbContext.Rollback()
		return result.Error
	}
	dbContext.Commit()
	return nil
}
