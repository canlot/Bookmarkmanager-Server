package Models

import "errors"

func GetBookmarksWithCategoryID(userId uint, categoryId uint) (bookmarks []Bookmark, success bool) {
	var user User
	var category Category

	if result := Database.Take(&user, userId); result.Error != nil {
		return bookmarks, false
	}

	if result := Database.Take(&category, categoryId); result.Error != nil {
		return bookmarks, false
	}

	var users []User
	if result := Database.Model(&category).Association("UsersAccess").Find(&users); result != nil {
		return bookmarks, false
	}

	var userExist = false
	for i := range users {
		if users[i].ID == user.ID {
			userExist = true
			break
		}
	}

	if userExist != true {
		return bookmarks, false
	}

	if result := Database.Find(&bookmarks, "category_id = ?", categoryId); result.Error != nil {
		return bookmarks, false
	}
	return bookmarks, true
}

func AddBookmark(userId uint, categoryId uint, bookmark Bookmark) (bookmarkret Bookmark, success bool) {
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
func EditBookmark(userId uint, categoryId uint, bookmarkId uint, bookmark Bookmark) error {
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
