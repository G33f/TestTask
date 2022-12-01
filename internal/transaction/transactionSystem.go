package transaction

import (
	"TestTask/internal/client"
	"context"
	"log"
)

var Queue map[client.Client]chan int

func transaction(ctx context.Context, c client.Client, amount int) {
	if _, err := Queue[c]; err {
		Queue[c] = make(chan int)
	}
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
		//TODO need save data in DB "client valet = c.amount" context
	}(c, Queue[c])
	Queue[c] <- amount
}
