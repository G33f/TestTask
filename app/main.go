package main

import (
	"context"
	"example.com/m/author/repository/PostgreSQL"
	"example.com/m/author/repository/PostgreSQL/utils/config"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
}

func main() {
	cfg := config.SetStorageConfig()
	fmt.Println(cfg)
	ctx := context.Background()
	err, _ := PostgreSQL.NewClient(ctx, cfg) //TODO change "_" to "storage"
	if err != nil {
		log.Panic(err)
	}
}
