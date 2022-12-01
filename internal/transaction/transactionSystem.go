package transaction

import (
	"TestTask/internal/client"
	"context"
	"log"
	"sync"
)

var Queue map[client.Client]chan int

func Transaction(ctx context.Context, c client.Client, postgresql client.Repository, amount int) {
	if _, err := Queue[c]; err {
		Queue[c] = make(chan int)
	}
	ls := sync.WaitGroup{}
	ls.Add(1)
	go func(c client.Client, a <-chan int) {
		select {
		case am := <-a:
			if am >= 0 {
				c.AddMoney(am)
			} else {
				if err := c.MinusMoney(am); err != nil {
					log.Println(err)
				}
			}
		}
		if err := postgresql.Update(ctx, c); err != nil {
			log.Println(err)
		}
		ls.Done()
	}(c, Queue[c])
	Queue[c] <- amount
	ls.Wait()
}

//func TransactionsBetweenClients(ctx context.Context, sender client.Client, receiver client.Client, postgresql client.Repository, amount int) {
//
//}
