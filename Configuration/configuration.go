package Configuration

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type AppConfiguration struct {
	ListenPort     int
	DatabaseConfig DatabaseConfig
	DeploymentMode string
}

type DatabaseConfig struct {
	DBProvider string
	Host       string
	Port       string
	Username   string
	Password   string
}

var Configuration AppConfiguration

func init() {
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

	err = viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("couln't parse config file: %w", err)
	}

}
