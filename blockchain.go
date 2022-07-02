package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp	%d\n", b.timestamp)
	fmt.Printf("previousHash	%x\n", b.previousHash)
	fmt.Printf("nonce	%d\n", b.nonce)
	for _, t := range b.transactions {
		t.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

type BlockChain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockChainAddress string
}

func NewBlockChain(blockChainAddress string) *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.blockChainAddress = blockChainAddress
	bc.CreatBlock(0, b.Hash())
	return bc
}

func (bc *BlockChain) CreatBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *BlockChain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// Poolにあるトランザクション値をコピーする
func (bc *BlockChain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value))
	}
	return transactions
}

// nonce, previousHash, transactionsからハッシュ値の先頭が0 x difficultyであればtrue
func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

// 次のブロックを作成する際に使用するnonceを返す
func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	// ValidProofがtrueになるまでfor文を回す
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockChainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreatBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *BlockChain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockchainAddress:    sender,
		recipientBlockchainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address  %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address  %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value  %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	// json用の構造体宣言はCapital、型が同様であれば別名で良い
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	myBlockchainAddress := "my_blockchain_address"

	blockChain := NewBlockChain(myBlockchainAddress)
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("C", "D", 2.0)
	blockChain.AddTransaction("X", "Y", 3.0)
	blockChain.Mining()
	blockChain.Print()

	fmt.Printf("my %.1f\n", blockChain.CalculateTotalAmount("my_blockchain_address"))
	fmt.Printf("C %.1f\n", blockChain.CalculateTotalAmount("C"))
	fmt.Printf("D %.1f\n", blockChain.CalculateTotalAmount("D"))
}
