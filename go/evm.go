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

// Evm Run runs the EVM code and returns the stack and a success indicator.
// @param byteCode - Feed the compiled bytecode for Evm to execute
func Evm(byteCode []byte) ([]*big.Int, bool) {
	// Invoke the Executor
	return Executor(byteCode)
}
