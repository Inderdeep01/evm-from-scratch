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
	"math/big"
)

type EVM struct {
	Tx Tx
}

// Tx stores all the data which accessed in smart contracts using 'tx' object
type Tx struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Origin   string `json:"origin"`
	GasPrice string `json:"gasprice"`
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

// Evm Run runs the EVM code and returns the stack and a success indicator.
// @param byteCode - Feed the compiled bytecode for Evm to execute
// @param tx       - All the data related to current transaction
// @param block    - All the data of current block
// @return stack   - For testcases to assess the correctness
// @return bool    - Indicating whether the execution was successful(true) or reverted(false)
func Evm(byteCode []byte, tx Tx, block Block) ([]*big.Int, bool) {
	// Invoke the Executor
	return Executor(byteCode, tx, block)
}
