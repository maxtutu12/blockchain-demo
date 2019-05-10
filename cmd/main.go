package main

import (
	block "blockchain-demo"
	"flag"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/kataras/iris"
)

var (
	port       = flag.String("port", "5000", "server listen to the port")
	identifier string
	blockChain *block.BlockChain
)

func main() {
	flag.Parse()

	identifier = strings.ReplaceAll(fmt.Sprintf("%s", uuid.Must(uuid.NewV4())), "-", "")
	fmt.Println("the identifier:", identifier)

	blockChain = block.NewBlockChain()

	app := iris.New()

	app.Get("/chain", getChainHandle)
	app.Post("/transactions/new", newTransactionHandle)
	app.Get("mine", mineHandle)

	addr := ":" + *port
	app.Run(iris.Addr(addr))
}

func getChainHandle(ctx iris.Context) {
	response := iris.Map{
		"chain":  blockChain.Blocks,
		"length": len(blockChain.Blocks),
	}

	ctx.JSON(response)
}

func newTransactionHandle(ctx iris.Context) {
	var ta block.Transaction

	err := ctx.ReadJSON(&ta)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	index := blockChain.NewTransaction(ta.Sender, ta.Recipient, ta.Amount)
	response := iris.Map{
		"message": "Transaction will be add to Block " + strconv.Itoa(int(index)),
	}

	ctx.JSON(response)
}

func mineHandle(ctx iris.Context) {
	blockChain.Mine(identifier)

	block := blockChain.LastBlock()
	fmt.Println("the last block:", block)

	response := iris.Map{
		"message":       "New Block Forged",
		"index":         block.Index,
		"proof":         block.Proof,
		"previous_hash": block.PreviousHash,
		"transactions":  block.Transactions,
	}

	ctx.JSON(response)
}
