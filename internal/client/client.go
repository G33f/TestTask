package client

import (
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
