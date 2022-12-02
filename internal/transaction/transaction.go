package transaction

type Transaction struct {
	Id       string
	Sender   string
	Receiver string
	Status   string
	Amount   int
}

func NewTransaction(sender, receiver string, amount int) Transaction {
	return Transaction{
		Id:       "",
		Sender:   sender,
		Receiver: receiver,
		Status:   "send",
		Amount:   amount,
	}
}
