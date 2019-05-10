package block

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

var curTransaction []Transaction

func NewTransaction(sender, recipient string, amount int) *Transaction {
	return &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}

type Block struct {
	Index        int64
	TimeStamp    int64
	Proof        int
	PreviousHash string
	Transactions []Transaction
}

func NewBlock(idx int64, proof int, preHash string) *Block {
	return &Block{
		Index:        idx,
		TimeStamp:    time.Now().Unix(),
		Proof:        proof,
		PreviousHash: preHash,
		Transactions: nil,
	}
}

func calcHash(b Block) string {
	jsonByte, err := json.Marshal(&b)
	if err != nil {
		return ""
	}

	hash := sha256.New()
	hash.Write(jsonByte)
	return hex.EncodeToString(hash.Sum(nil))
}

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock(0, 100, "1")
	blockChain := BlockChain{}
	blockChain.Blocks = append(blockChain.Blocks, genesisBlock)

	return &blockChain
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func validProof(lastProof, proof int) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(proof)

	h := sha256.New()
	h.Write([]byte(guess))

	guessHash := hex.EncodeToString(h.Sum(nil))
	return guessHash[:4] == "0000"
}

func (bc *BlockChain) proofOfWork(lastProof int) int {
	proof := 0

	for {
		if validProof(lastProof, proof) {
			break
		}

		proof += 1
	}

	return proof
}

func (bc *BlockChain) NewTransaction(sender, recipient string, amount int) int64 {
	newTrans := NewTransaction(sender, recipient, amount)
	curTransaction = append(curTransaction, *newTrans)

	return bc.LastBlock().Index + 1
}

func (bc *BlockChain) Mine(identifier string) {
	lastBlock := bc.LastBlock()
	hash := calcHash(*lastBlock)
	proof := bc.proofOfWork(lastBlock.Proof)

	bc.NewTransaction("0", identifier, 1)

	block := NewBlock(lastBlock.Index+1, proof, hash)
	block.Transactions = curTransaction

	curTransaction = curTransaction[0:0]

	bc.Blocks = append(bc.Blocks, block)
}

func ValidChain(blocks []*Block) bool {
	lastBlock := blocks[0]
	curIndex := 1
	end := len(blocks)

	for {
		if curIndex >= end {
			break
		}

		curBlock := blocks[curIndex]
		fmt.Println("curblock:", curBlock)
		fmt.Println("lastblock:", lastBlock)

		if curBlock.PreviousHash != calcHash(*lastBlock) {
			return false
		}

		if !validProof(lastBlock.Proof, curBlock.Proof) {
			return false
		}

		lastBlock = curBlock
		curIndex += 1
	}

	return true
}
