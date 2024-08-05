// Package evm is an **incomplete** implementation of the Ethereum Virtual
// Machine for the "EVM From Scratch" course:
// https://github.com/w1nt3r-eth/evm-from-scratch
//
// To work on EVM From Scratch In Go:
//
// - Install Golang: https://golang.org/doc/install
// - Go to the `go` directory: `cd go`
// - Edit `evm.go` (this file!), see TODO below
// - Run `go test ./...` to run the tests
package evm

import (
	"github.com/holiman/uint256"
	"math/big"
)

// Account is used to represent an account on EVM
type Account struct {
	Balance string `json:"balance"`
	Code    Code   `json:"code"`
}

type Code struct {
	Bin string `json:"bin"`
}

// State Represents the current state of EVM
// @dev State of EVM is a mapping of an addresses to the Account
type State map[string]Account

// ContractStorage is a mapping of 32 bytes key to 32 bytes value
// @dev Since Golang does not allow []byte to be the key type, we can encode the bytes to string for storing it
type ContractStorage map[string][]byte

var Storage = ContractStorage{}

type EVMMemory struct {
	Memory          []byte
	StartingIndices []int
}

// Tx stores all the data which accessed in smart contracts using 'tx' or 'msg' object
// @notice Although 'tx' and 'msg' are different objects in solidity, for the simplicity sake, here they are stored in
// same object. The opcodes to access those values work just fine
type Tx struct {
	To       string `json:"to"`
	From     string `json:"from"`   // msg.sender
	Origin   string `json:"origin"` // tx.origin
	GasPrice string `json:"gasprice"`
	Value    string `json:"value"` // msg.value
	Data     string `json:"data"`
}

// Block stores all the data of the relevant block; accessed in smart contracts using 'block' object
type Block struct {
	BaseFee    string `json:"basefee"`
	CoinBase   string `json:"coinbase"`
	Timestamp  string `json:"timestamp"`
	Number     string `json:"number"`
	Difficulty string `json:"difficulty"`
	GasLimit   string `json:"gaslimit"`
	ChainID    string `json:"chainid"`
}

type Log struct {
	Address string   `json:"address"`
	Data    string   `json:"data"`
	Topics  []string `json:"topics"`
}

// EVMContext is the combination of all the above structs - to make function calls more convenient
type EVMContext struct {
	State          State
	Tx             Tx
	Block          Block
	Bytecode       []byte
	Stack          []*uint256.Int
	Memory         []byte
	ProgramCounter int
	GasConsumed    uint256.Int
	GasRefund      uint256.Int
	Logs           []Log
	// Flags and other important variables
	isStaticCall  bool
	OpCodeCounter int
	Success       bool
	HaltExecution bool
}

type ReturnContext struct {
	Stack  []*big.Int
	Logs   []Log
	Return string
}

// Evm Run runs the EVM code and returns the stack and a success indicator.
// @param byteCode - Feed the compiled bytecode for Evm to execute
// @param tx       - All the data related to current transaction
// @param block    - All the data of current block
// @param state    - Current state of EVM (mapping of address to Account)
// @return stack   - For testcases to assess the correctness
// @return bool    - Indicating whether the execution was successful(true) or reverted(false)
func Evm(byteCode []byte, tx Tx, block Block, state State) ([]*big.Int, []Log, string, bool) {
	// Invoke the Executor
	return Executor(byteCode, tx, block, state, false)
}
