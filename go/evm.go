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
	1:  add,
	2:  multiply,
	3:  sub,
	4:  div,
	5:  sdiv,
	6:  mod,
	7:  smod,
	8:  addmod,
	9:  mulmod,
	10: exp,
	11: signextend,
	16: lt,
	17: gt,
	18: slt,
	19: sgt,
	20: eq,
	21: isZero,
	22: and,
	23: or,
	24: xor,
	25: not,
	26: getByte,
	27: shl,
	28: shr,
	29: sar,
}

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0
	var success = true
	var memory map[int64][]*byte
	_ = memory

	for pc < len(code) {
		var tempStack []*big.Int
		var successFlag bool
		op := code[pc]
		n := 0

		// TODO: Implement the EVM here!

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
		if op >= 1 && op <= 29 {
			stack, n, successFlag = functionMap[op](stack)
		} else if op >= 95 && op <= 127 {
			// opcodes from 95 to 197 are for pushing to stack
			tempStack, n, successFlag = pushN(code[pc:])
			stack = append(tempStack, stack...)
		} else if op >= 128 && op <= 143 {
			stack, n, successFlag = dupN(code[pc], stack)
		} else if op >= 144 && op <= 159 {
			stack, n, successFlag = swapN(code[pc], stack)
		} else if op == 88 {
			x := big.NewInt(int64(pc))
			stack = append([]*big.Int{x}, stack...)
			successFlag = true
		} else if op == 90 {
			y := getGas()
			stack = append([]*big.Int{y}, stack...)
			successFlag = true
		} else if op == 86 {
			// check for valid JUMPDEST
			var validJumpDest []int
			for i := pc + 1; i < len(code); i++ {
				currOp := code[i]
				if currOp >= 95 && currOp <= 127 {
					n := int(currOp - 95)
					i += n
				} else if currOp == 91 {
					validJumpDest = append(validJumpDest, i)
				}
			}
			nJumpDest := int(stack[0].Int64())
			isValidJumpDest := false
			for i := 0; i < len(validJumpDest); i++ {
				if validJumpDest[i] == nJumpDest {
					isValidJumpDest = true
				}
			}
			if nJumpDest > len(code) || !isValidJumpDest {
				success = false
				stack = stack[1:]
				break
			} else {
				pc = nJumpDest
				stack = stack[1:]
				continue
			}

		} else if op == 87 {

			// check for valid JUMPDEST
			var validJumpDest []int
			for i := pc + 1; i < len(code); i++ {
				currOp := code[i]
				if currOp >= 95 && currOp <= 127 {
					n := int(currOp - 95)
					i += n
				} else if currOp == 91 {
					validJumpDest = append(validJumpDest, i)
				}
			}
			nJumpDest := int(stack[0].Int64())
			fmt.Println(nJumpDest)
			fmt.Println(code[nJumpDest])
			isValidJumpDest := false
			for i := 0; i < len(validJumpDest); i++ {
				if validJumpDest[i] == nJumpDest {
					isValidJumpDest = true
				}
			}
			if nJumpDest > len(code) || !isValidJumpDest {
				success = false
				stack = stack[2:]
				break
			} else if int(stack[1].Int64()) == 0 {
				stack = stack[2:]
				successFlag = true
			} else {
				pc = nJumpDest
				stack = stack[2:]
				continue
			}
		} else if op == 91 {
			successFlag = true
		} else {
			fmt.Println("******** op *******", op)
			fmt.Println("********* code ******", code)
			fmt.Println("Stack:", stack)
		}
		pc += n + 1
		success = successFlag && success
	}
	return stack, success
}
