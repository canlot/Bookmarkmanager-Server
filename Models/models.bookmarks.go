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
	if result := Database.Create(&bookmark); result.Error != nil {
		return bookmark, false
	}
	return bookmark, true
}
func EditBookmark(userId uint, categoryId uint, bookmark Bookmark) error {
	var category Category

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return result.Error
	}
	if category.OwnerID != userId {
		return errors.New("Owner of bookmark not the same as loged on user")
	}

	if result := Database.Save(&bookmark); result.Error != nil {
		return result.Error
	}
	return nil
}
func DeleteBookmark(userId uint, categoryId uint, bookmarkId uint) error {
	var category Category

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return result.Error
	}
	if category.OwnerID != userId {
		return errors.New("Owner of bookmark not the same as loged on user")
	}

	if result := Database.Delete(&Bookmark{}, bookmarkId); result.Error != nil {
		return result.Error
	}
	return nil
}
