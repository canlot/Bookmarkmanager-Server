package Configuration

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Configuration struct {
	ListenPort     int
	DatabaseConfig DatabaseConfig
}

type DatabaseConfig struct {
	DBProvider DatabaseType
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
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
	if Environment == Test {
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
	if err != nil {
		log.Fatalf("couln't parse config file: %w", err)
	}

}
