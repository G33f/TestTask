package transactionDataHandler

import "TestTask/internal/client"

type clientDBData struct {
	id        int
	firstName string
	lastName  string
	balance   int
}

type ClientDataHandler interface {
}

func (cdb *clientDBData) SaveInDataBase(c client.Client) {
}
