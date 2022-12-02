package main

import (
	"TestTask/handler/serverHandler"
	"TestTask/internal/client"
	clientDB "TestTask/internal/client/db"
	"TestTask/internal/server"
	"TestTask/internal/transaction"
	transactionDB "TestTask/internal/transaction/db"
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

	Server := server.Server{}
	if err = Server.NewConnection(); err != nil {
		log.Println(err)
		return
	}

	for {
		conn, err := Server.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		} else {
			serverHandler.HandleRequest(ctx, conn, postgresqlTransaction, postgresqlClient)
		}
	}
}
