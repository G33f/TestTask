package main

import (
	"TestTask/internal/client"
	client2 "TestTask/internal/client/db"
	"TestTask/internal/transaction"
	"TestTask/pkg/logging"
	"TestTask/pkg/repository/PostgreSQL"
	"TestTask/pkg/repository/PostgreSQL/utils/config"
	"context"
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
	err, clientDB := PostgreSQL.NewClient(ctx, cfg)
	if err != nil {
		log.Panic(err)
	}
	logger := logging.GetLogger()
	postgresql := client2.NewRepository(clientDB, logger)
	err = client.GetClientFromDB(ctx, postgresql)
	if err != nil {
		log.Println(err)
	}

	transaction.MakeTransaction(ctx, client.ClientMaps["eca548a8-8833-4152-8629-9e0fa69d8ca7"], postgresql, -100)
}
