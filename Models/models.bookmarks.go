package Models

import ()

func GetBookmarksWithCategoryID(userId uint, categoryId uint) (bookmarks []Bookmark, success bool) {
	var user User
	var category Category
	if result := Database.Take(&user, userId); result.Error != nil {
		return bookmarks, false
	}
	if result := Database.Take(&category, categoryId); result.Error != nil {
		return bookmarks, false
	}
	if result := Database.Model(&user).Association("BookmarksFullAccess").Find(&bookmarks, "category_id = ?", categoryId); result == nil {
		return bookmarks, true
	}
	return bookmarks, false
}
func AddBookmark(userId uint, categoryId uint, bookmark Bookmark) (bookmarkret Bookmark, success bool) {
	var user User
	var category Category
	if result := Database.Take(&user, userId); result.Error != nil {
		return bookmark, false
	}
	if result := Database.Take(&category, categoryId); result.Error != nil {
		return bookmark, false
	}
	if category.OwnerID != userId {
		return bookmark, false
	}
	bookmark.CategoryID = category.ID
	if result := Database.Create(&bookmark); result.Error == nil {
		if result := Database.Model(&bookmark).Association("UsersFullAccess").Append(&user); result == nil {
			return bookmark, true
		}
	}
	return bookmark, false
}
