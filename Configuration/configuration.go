package Configuration

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Configuration struct {
	ListenPort     int
	TokenLifetime  string
	IconFolderPath string
	DatabaseConfig DatabaseConfig
	SslEncryption  SslEncryption
	SetUpUser      SetUpUser
}

type DatabaseConfig struct {
	DBProvider DatabaseType
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
}

type SslEncryption struct {
	Enabled  bool
	CertPath string
	KeyPath  string
}
type SetUpUser struct {
	Email    string
	Name     string
	Password string
}

type DatabaseType string

const (
	Sqlite DatabaseType = "Sqlite"
	Mysql               = "Mysql"
)

type EnvironmentType int

const (
	Production EnvironmentType = iota
	Test
	Debug
)

var Environment EnvironmentType

var AppConfiguration Configuration

func init() {
	Environment = Test
}

func GetConfig() {
	//viper.SetDefault("ListenPort", 8080)
	if Environment == Test || Environment == Debug {
		AppConfiguration.ListenPort = 8080
		if AppConfiguration.TokenLifetime == "" {
			AppConfiguration.TokenLifetime = "1h"
		}
		if AppConfiguration.IconFolderPath == "" {
			AppConfiguration.IconFolderPath = "./icons"
		}
		setUpThings()
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	path, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	viper.AddConfigPath(path)

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Fatalf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(&AppConfiguration)
	log.Println(AppConfiguration.ListenPort)
	if err != nil {
		log.Fatalf("couln't parse config file: %w", err)
	}

	err = checkConfig()
	if err != nil {
		log.Fatalf("config did not pass %w", err)
	}
	if AppConfiguration.TokenLifetime == "" {
		AppConfiguration.TokenLifetime = "1h"
	}
	if AppConfiguration.IconFolderPath == "" {
		AppConfiguration.IconFolderPath = "./icons"
	}

	setUpThings()
}
func checkConfig() error {
	if AppConfiguration.ListenPort <= 0 || AppConfiguration.ListenPort >= 65536 {
		return errors.New("Listen port is not in the port range")
	}
	return nil
}
func setUpThings() error {
	createIconsFolder()
	return nil
}
func createIconsFolder() {
	err := os.Mkdir(AppConfiguration.IconFolderPath, 0777)
	if err != nil {
	}
}
