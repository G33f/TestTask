package transaction

import (
	"TestTask/internal/client"
	"context"
	"log"
	"sync"
)

func recoverTransactionHistory(ctx context.Context, postgresql Repository, clientRep client.Repository) {
	transactions, err := postgresql.FindAllNotDone(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	for _, val := range transactions {
		switch val.Status {
		case "send":
			TransactionsBetweenClients(ctx, val, postgresql, clientRep)
		case "receive":
			err := MakeTransaction(ctx, val.Receiver, clientRep, val.Amount, &sync.WaitGroup{})
			if err != nil {
				continue
			}
		}
	}
}
