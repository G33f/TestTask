package transaction

type Transaction struct {
	Id       string
	Sender   string
	Receiver string
	Status   string
	Amount   int
}
