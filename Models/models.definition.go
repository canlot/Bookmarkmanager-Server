package Models

import (
	"gorm.io/gorm"
	"time"
)

var ShowPassword bool = false

type Category struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ParentID    uint           `json:"parentid"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Shared      bool           `json:"shared"`
	OwnerID     uint           `json:"ownerid"`
	Bookmarks   []Bookmark     `json:"-"`
	UsersAccess []User         `json:"-" gorm:"many2many:user_categories;"`
}

type Bookmark struct {
	ID uint `gorm:"primaryKey"`
	//CreatedAt   time.Time      `gorm:"<-:create" json:",omitempty"`
	CreatedAt   time.Time      `json:",omitempty"`
	UpdatedAt   time.Time      `json:",omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CategoryID  uint           `json:"categoryid"`
	Url         string         `json:"url"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	IconName    string         `json:"iconname"`
}

type User struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Email            string         `json:"email"`
	Name             string         `json:"name"`
	Password         string         `json:"-"`
	Administrator    bool           `json:"administrator"`
	CategoriesAccess []Category     `json:"-" gorm:"many2many:user_categories;"`
}

type JsonError struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
