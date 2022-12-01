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

var ClientMaps map[string]client.Client

func GetClientFromDB(ctx context.Context, postgresql client.Repository) error {
	clientArray, err := postgresql.FindAll(ctx)
	if err != nil {
		return err
	}
	clientMap := map[string]client.Client{}
	clientTransferQueue := map[client.Client]chan int{}
	for _, val := range clientArray {
		c := client.NewClient(val.FirstName, val.LastName, val.Balance)
		clientMap[val.Id] = c
		clientTransferQueue[c] = make(chan int)
	}
	transaction.Queue = clientTransferQueue
	ClientMaps = clientMap
	return nil
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
	err = GetClientFromDB(ctx, postgresql)
	if err != nil {
		log.Println(err)
	}
	for k, val := range ClientMaps {
		fmt.Println("Client ID is ", k, " and his valet data is", val)
	}
}
