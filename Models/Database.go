package Models

import (
	"Bookmarkmanager-Server/Configuration"
	"Bookmarkmanager-Server/Helpers"
	"errors"
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

	if Configuration.Environment != Configuration.Production {
		return
	}
	if err = UserSetUp(); err != nil {
		panic(err)
	}
}

func UserSetUp() error {
	var db *gorm.DB
	if db = Database.Take(&User{}, "administrator = ?", true); db.Error != nil {
		return db.Error
	}
	if db.RowsAffected > 0 { // if one user with admin found, return
		return nil
	}
	if Configuration.AppConfiguration.SetUpUser.Email == "" ||
		Configuration.AppConfiguration.SetUpUser.Password == "" ||
		Configuration.AppConfiguration.SetUpUser.Name == "" {
		return errors.New("configuration not complete")
	}
	password, err := Helpers.CreateHashFromPassword(Configuration.AppConfiguration.SetUpUser.Password)
	if err != nil {
		return err
	}
	var user = User{
		Email:         Configuration.AppConfiguration.SetUpUser.Email,
		Name:          Configuration.AppConfiguration.SetUpUser.Name,
		Password:      password,
		Administrator: true,
	}
	if db = Database.Create(&user); db.Error != nil {
		return db.Error
	}
	return nil
}
