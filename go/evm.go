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
	"fmt"
	"math/big"
)

var functionMap = map[byte]func([]*big.Int) ([]*big.Int, int, bool){
	1: add,
	2: multiply,
	3: sub,
	4: div,
	6: mod,
	8: addmod,
}

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0
	var success = true

	for pc < len(code) {
		var tempStack []*big.Int
		var successFlag bool
		op := code[pc]
		n := 0

		// TODO: Implement the EVM here!
		fmt.Println("******** op *******", op)
		fmt.Println("********* code ******", code[pc:])
		fmt.Println("Stack:", stack)
		// STOP
		if op == 0 {
			successFlag = true
			break
		}
		switch op {

		// POP
		case 80:
			stack = stack[1:]
			successFlag = true
		}

		// Arithmetic Operations
		if op >= 1 && op <= 11 {
			stack, n, successFlag = functionMap[op](stack)
		}

		// opcodes from 95 to 197 are for pushing to stack
		if op >= 95 && op <= 197 {
			tempStack, n, successFlag = pushN(code[pc:])
			stack = append(tempStack, stack...)

		}
		pc += n + 1
		success = successFlag && success
		fmt.Println("this is executing")
	}

	return stack, success
}
