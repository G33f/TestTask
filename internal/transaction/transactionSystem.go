package transaction

import (
	"TestTask/internal/client"
	"context"
	"fmt"
	"log"
	"sync"
)

var Queue map[string]chan int

func MakeTransaction(ctx context.Context, cli string, postgresql client.Repository, amount int, wg *sync.WaitGroup) error {
	var e error
	c, err := postgresql.FindOne(ctx, cli)
	if err != nil {
		return err
	}
	if _, err := Queue[cli]; err {
		Queue[cli] = make(chan int)
	}
	go func(c client.Client, a <-chan int) {
		defer wg.Done()
		select {
		case am := <-a:
			if am >= 0 {
				c.AddMoney(am)
			} else {
				if err := c.MinusMoney(am); err != nil {
					log.Println(err)
					e = err
					return
				}
			}
		}
		if err := postgresql.Update(ctx, c); err != nil {
			log.Println(err)
		}
	}(c, Queue[cli])
	Queue[cli] <- amount
	return e
}

func TransactionsBetweenClients(ctx context.Context, transaction Transaction, postgresql Repository, clientRep client.Repository) {
	wg := sync.WaitGroup{}
	fmt.Println(transaction)
	wg.Add(1)
	err := MakeTransaction(ctx, transaction.Sender, clientRep, -transaction.Amount, &wg)
	if err != nil {
		transaction.Status = "closest"
		err = postgresql.Update(ctx, transaction)
		return
	}
	wg.Wait()
	transaction.Status = "receive"
	err = postgresql.Update(ctx, transaction)
	if err != nil {
		return
	}
	wg.Add(1)
	err = MakeTransaction(ctx, transaction.Receiver, clientRep, transaction.Amount, &wg)
	if err != nil {
		return
	}
	transaction.Status = "closest"
	err = postgresql.Update(ctx, transaction)
	wg.Wait()
}
