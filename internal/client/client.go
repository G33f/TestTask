package client

import (
	"context"
	"fmt"
)

type Client struct {
	Id        string
	FirstName string
	LastName  string
	Balance   int
}

func (c *Client) MinusMoney(amount int) error {
	if c.Balance+amount < 0 {
		return fmt.Errorf("you dont have enogh money on your accaunt")
	}
	c.Balance += amount
	return nil
}

func (c *Client) AddMoney(amount int) {
	c.Balance += amount
}

func NewClient(firstName, secondName string, valet int) Client {
	return Client{FirstName: firstName, LastName: secondName, Balance: valet}
}

func GetClientFromDB(ctx context.Context, postgresql Repository) (error, map[string]chan int) {
	clientArray, err := postgresql.FindAll(ctx)
	if err != nil {
		return err, nil
	}
	clientTransferQueue := map[string]chan int{}
	for _, val := range clientArray {
		c := NewClient(val.FirstName, val.LastName, val.Balance)
		c.Id = val.Id
		clientTransferQueue[val.Id] = make(chan int)
	}
	return nil, clientTransferQueue
}
