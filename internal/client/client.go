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

var ClientMaps map[string]Client

func GetClientFromDB(ctx context.Context, postgresql Repository) (error, map[Client]chan int) {
	clientArray, err := postgresql.FindAll(ctx)
	if err != nil {
		return err, nil
	}
	clientMap := map[string]Client{}
	clientTransferQueue := map[Client]chan int{}
	for _, val := range clientArray {
		c := NewClient(val.FirstName, val.LastName, val.Balance)
		c.Id = val.Id
		clientMap[val.Id] = c
		clientTransferQueue[c] = make(chan int)
	}
	ClientMaps = clientMap
	return nil, clientTransferQueue
}
