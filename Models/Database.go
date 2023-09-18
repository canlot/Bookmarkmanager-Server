package Models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseType int

const (
	Sqlite DatabaseType = iota
	Mysql
)

var Database *gorm.DB

func DatabaseConfig() {
	var err error
	databaseType := Sqlite

	if databaseType == Sqlite {
		Database, err = gorm.Open(sqlite.Open("C:\\Users\\Jakob\\go\\src\\github.com\\canlot\\Bookmarkmanager-Server\\gorm.db"), &gorm.Config{})
	} else if databaseType == Mysql {
		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "gotest", "gotest", "localhost", 3306, "gotest")
		Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		panic("Failed to connect to database")
	}
	Database.AutoMigrate(&User{}, &Category{}, &Bookmark{})
}
