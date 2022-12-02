package serverHandler

import (
	"TestTask/internal/client"
	"TestTask/internal/transaction"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"log"
	"net"
	"strconv"
	"strings"
)

func HandleRequest(ctx context.Context, c net.Conn, postgresqlTransaction transaction.Repository, postgresqlClient client.Repository) {
	defer c.Close()
	buf := make([]byte, 256)
	_, err := c.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	tmp := cast.ToString(buf)
	request := strings.Split(tmp, " ")
	if len(request) != 3 {
		log.Println("Wrong transaction!")
		return
	}
	tmp1 := strings.Split(request[2], "\n")
	amount, err := strconv.Atoi(tmp1[0])
	if amount < 0 {
		log.Println("Wrong transaction!")
		return
	}
	tra := transaction.NewTransaction(request[0], request[1], amount)
	err = postgresqlTransaction.Create(ctx, &tra)
	fmt.Println(tra)
	if err != nil {
		log.Println(err)
		return
	}
	go transaction.TransactionsBetweenClients(ctx, tra, postgresqlTransaction, postgresqlClient)
}
