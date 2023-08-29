package Models

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
