package Models

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseType int

const (
	Sqlite DatabaseType = iota
	Mysql
)

type Environment int

const (
	Production Environment = iota
	Test
)

var Database *gorm.DB

func DatabaseConfig(databaseType DatabaseType, environment Environment) {
	var err error
	if databaseType == Sqlite {
		//Database, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		Database, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	} else if databaseType == Mysql {
		var dsn string
		if environment == Production {
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "gotest", "gotest", "localhost", 3306, "gotest")
		} else if environment == Test {
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "gotest", "gotest", "localhost", 3306, "test")
		}
		Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		panic("Failed to connect to database")
	}
	Database.AutoMigrate(&User{}, &Category{}, &Bookmark{})
}
