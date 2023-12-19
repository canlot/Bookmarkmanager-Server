package Models

import (
	"Bookmarkmanager-Server/Configuration"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func DatabaseConfig() {
	var err error
	if Configuration.Environment == Configuration.Production {
		if Configuration.AppConfiguration.DatabaseConfig.DBProvider == Configuration.Sqlite {
			Database, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		} else if Configuration.AppConfiguration.DatabaseConfig.DBProvider == Configuration.Mysql {
			var dsn string
			dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", Configuration.AppConfiguration.DatabaseConfig.Username,
				Configuration.AppConfiguration.DatabaseConfig.Password, Configuration.AppConfiguration.DatabaseConfig.Host, Configuration.AppConfiguration.DatabaseConfig.Port,
				Configuration.AppConfiguration.DatabaseConfig.Database)
			Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		}
	} else if Configuration.Environment == Configuration.Test || Configuration.Environment == Configuration.Debug {
		Database, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	}

	if err != nil {
		panic("Failed to connect to database")
	}
	Database.AutoMigrate(&User{}, &Category{}, &Bookmark{})

}
