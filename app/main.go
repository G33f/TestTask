package main

import (
	"TestTask/internal/client"
	clientDB "TestTask/internal/client/db"
	"TestTask/internal/transaction"
	transactionDB "TestTask/internal/transaction/db"
	"TestTask/pkg/logging"
	"TestTask/pkg/repository/PostgreSQL"
	"TestTask/pkg/repository/PostgreSQL/utils/config"
	"context"
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

	ctx := context.Background()
	err, dBClient := PostgreSQL.NewClient(ctx, cfg)
	if err != nil {
		log.Panic(err)
	}
	logger := logging.GetLogger()

	postgresqlClient := clientDB.NewRepository(dBClient, logger)
	postgresqlTransaction := transactionDB.NewRepository(dBClient, logger)

	err, transaction.Queue = client.GetClientFromDB(ctx, postgresqlClient)
	if err != nil {
		log.Println(err)
	}

	transaction.RecoverTransactionHistory(ctx, postgresqlTransaction, postgresqlClient)

}
