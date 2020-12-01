package Models

import (
	"gorm.io/gorm"
)

var showPassword bool = false

type Category struct {
	gorm.Model
	ParentID           uint       `json:"parent_id"`
	Name               string     `json:"name"`
	Shared             bool       `json:"shared"`
	OwnerID            uint       `json:"owner_id"`
	Bookmarks          []Bookmark `json:"-"`
	UsersFullAccess    []User     `json:"-" gorm:"many2many:user_categories_full;"`
	UsersInheritAccess []User     `json:"-" gorm:"many2many:user_categories_inherit"`
}

type Bookmark struct {
	gorm.Model
	CategoryID      uint   `json:"category_id"`
	Url             string `json:"url"`
	UsersFullAccess []User `json:"-" gorm:"many2many:user_bookmarks_full;"`
}

type User struct {
	gorm.Model
	Name                    string     `json:"name"`
	Password                string     `json:"password"`
	Administrator           bool       `json:"administrator"`
	CategoriesFullAccess    []Category `json:"-" gorm:"many2many:user_categories_full;"`
	CategoriesInheritAccess []Category `json:"-" gorm:"many2many:user_categories_inherit;"`
	BookmarksFullAccess     []Bookmark `json:"-" gorm:"many2many:user_bookmarks_full;"`
}

type JsonError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
