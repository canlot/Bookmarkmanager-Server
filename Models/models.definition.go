package Models

import (
	"gorm.io/gorm"
)

var ShowPassword bool = false

type Category struct {
	gorm.Model
	ParentID    uint       `json:"parentid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Shared      bool       `json:"shared"`
	OwnerID     uint       `json:"ownerid"`
	Bookmarks   []Bookmark `json:"-"`
	UsersAccess []User     `json:"-" gorm:"many2many:user_categories;"`
}

type Bookmark struct {
	gorm.Model
	CategoryID  uint   `json:"categoryid"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	gorm.Model
	Name             string     `json:"name"`
	Password         string     `json:"-"`
	Administrator    bool       `json:"administrator"`
	CategoriesAccess []Category `json:"-" gorm:"many2many:user_categories;"`
}

type JsonError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
