package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

func LoadConfig() *Config {
	// .envファイルを作成する(.env.exampleに例があるのでコピーして.envを作成する)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// 本当ならdefault値を設定するがこのサンプルでは実装しない

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: could not read .env file (%v). Using environment variables.", err)
	} else {
		log.Printf("Loaded config from: %s", viper.ConfigFileUsed())
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
	}

	if cfg.Database.Password == "" {
		log.Fatal("DB_PASSWORD is required")
	}

	return cfg
}
