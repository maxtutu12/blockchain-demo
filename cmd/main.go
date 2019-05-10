package main

import (
	block "blockchain-demo"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/kataras/iris"
)

var (
	port       = flag.String("port", "5000", "server listen to the port")
	identifier string
	blockChain *block.BlockChain
	netNodes   = make(map[string]string)
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
	app.Post("/nodes/register", registerNodesHandle)
	app.Get("/nodes/resolve", consensusHandle)

	addr := ":" + *port
	fmt.Println("The Sever Listen on ", addr)
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

type node struct {
	Name []string `json:"node"`
}

func registerNodesHandle(ctx iris.Context) {
	var n node
	err := ctx.ReadJSON(&n)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	for _, node := range n.Name {
		fmt.Println("register node:", node)
		u, err := url.Parse(node)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
		netNodes[u.Host] = u.Host
	}

	response := iris.Map{
		"message":     "New nodes have been added",
		"total_nodes": netNodes,
	}

	ctx.JSON(response)
}

type ChainRespMsg struct {
	Blocks []*block.Block `json:"chain"`
	Length int
}

func resolveConflict() bool {
	maxLength := len(blockChain.Blocks)

	var tempChain ChainRespMsg

	for node := range netNodes {
		reqUrl := "http://" + node + "/chain"
		resp, err := http.Get(reqUrl)
		if err != nil {
			fmt.Println("get chain request err:", err)
			continue
		}

		if resp.StatusCode == 200 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("read chain data err:", err)
				continue
			}

			err = json.Unmarshal(data, &tempChain)
			if err != nil {
				fmt.Println("unmarshal chain response err:", err)
				continue
			}

			if tempChain.Length > maxLength && block.ValidChain(tempChain.Blocks) {
				maxLength = tempChain.Length
				blockChain.Blocks = tempChain.Blocks
			}
		}
	}

	return true
}

func consensusHandle(ctx iris.Context) {
	replaced := resolveConflict()

	var msg string
	if replaced {
		msg = "Our chain was replaced"
	} else {
		msg = "Our chain is authoritative"
	}

	response := iris.Map{
		"message": msg,
		"chain":   blockChain.Blocks,
	}

	ctx.JSON(response)
}
