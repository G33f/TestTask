package config

import (
	"github.com/spf13/viper"
)

type storageConfig struct {
	Host        string
	Port        string
	Database    string
	Username    string
	Password    string
	MaxAttempts int
}

type StorageConfig interface {
	GetStorageConfig() (Host, Port, Database, Username, Password string)
	GetMaxAttempt() int
}

func SetStorageConfig() StorageConfig {
	var db storageConfig
	db.Host = viper.GetString("storage.host")
	db.Port = viper.GetString("storage.port")
	db.Database = viper.GetString("storage.database")
	db.Username = viper.GetString("storage.username")
	db.Password = viper.GetString("storage.password")
	db.MaxAttempts = viper.GetInt("storage.maxAttempts")
	return &db
}

func (sc *storageConfig) GetStorageConfig() (Host, Port, Database, Username, Password string) {
	return sc.Host, sc.Port, sc.Database, sc.Username, sc.Password
}

func (sc *storageConfig) GetMaxAttempt() int {
	return sc.MaxAttempts
}
