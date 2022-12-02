package transaction

import (
	"TestTask/internal/client"
	"context"
	"log"
	"sync"
)

func RecoverTransactionHistory(ctx context.Context, postgresqlTransaction Repository, postgresqlClient client.Repository) {
	transactions, err := postgresqlTransaction.FindAllNotDone(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	for _, val := range transactions {
		switch val.Status {
		case "send":
			go TransactionsBetweenClients(ctx, val, postgresqlTransaction, postgresqlClient)
		case "receive":
			err := MakeTransaction(ctx, val.Receiver, postgresqlClient, val.Amount, &sync.WaitGroup{})
			if err != nil {
				continue
			}
		}
	}
}
