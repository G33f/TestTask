package main

import (
	"TestTask/internal/client"
	client2 "TestTask/internal/client/db"
	"TestTask/internal/transaction"
	transaction2 "TestTask/internal/transaction/db"
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
	err, dBClient := PostgreSQL.NewClient(ctx, cfg)
	if err != nil {
		log.Panic(err)
	}
	logger := logging.GetLogger()
	postgresql := client2.NewRepository(dBClient, logger)
	postgresql2 := transaction2.NewRepository(dBClient, logger)
	err = client.GetClientFromDB(ctx, postgresql)
	if err != nil {
		log.Println(err)
	}

	transaction.TransactionsBetweenClients(ctx, transaction.NewTransaction("ccd04378-17dc-4857-8d04-9c3c939a319f", "eca548a8-8833-4152-8629-9e0fa69d8ca7", 500), postgresql2, postgresql)
}
