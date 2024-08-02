package evm

import (
	"fmt"
	"math/big"
)

// functionMap Stores a mapping of opcodes to their respective functions
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

// Executor is responsible for execution of all the valid opcode by parsing the incoming bytecode
// @pram   byteCode - The bytecode compiled from a Smart Contract
// @return stack    - The current state of stack
// @return bool     - Representing whether the execution was successful(true) or reverted(false)
func Executor(byteCode []byte, tx Tx, block Block) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0 // The Program Counter
	var success = true
	// Valid Jump Destinations
	var validJumpDestinations []int
	// Memory of the EVM - Volatile
	// Need to initialize maps in go; attempts to write to a nil map will cause a runtime panic
	var memory []byte

	for pc < len(byteCode) {
		var tempStack []*big.Int
		var successFlag bool
		op := byteCode[pc]
		n := 0

		// TODO: Implement the EVM here!

		if op == 0 { // STOP
			successFlag = true
			break
		} else if op == 80 { // POP
			stack = stack[1:]
			successFlag = true
		} else if op >= 1 && op <= 29 {
			// Arithmetic Operations
			stack, n, successFlag = functionMap[op](stack)
		} else if op >= 95 && op <= 127 { // PUSH N bytes to stack
			// opcodes from 95 to 197 are for pushing to stack
			tempStack, n, successFlag = pushN(byteCode[pc:])
			stack = append(tempStack, stack...)
		} else if op >= 128 && op <= 143 { // DUP N
			stack, n, successFlag = dupN(byteCode[pc], stack)
		} else if op >= 144 && op <= 159 { // SWAP nth element in stack with top
			stack, n, successFlag = swapN(byteCode[pc], stack)
		} else if op == 88 { // PC
			x := big.NewInt(int64(pc))
			stack = append([]*big.Int{x}, stack...)
			successFlag = true
		} else if op == 90 { // GAS
			y := getGas()
			stack = append([]*big.Int{y}, stack...)
			successFlag = true
		} else if op == 86 { // JUMP
			// check for valid JUMPDEST
			if len(validJumpDestinations) == 0 {
				validJumpDestinations = getValidJumpDestinations(byteCode, pc)
			}

			nJumpDest := int(stack[0].Int64())
			isValidJumpDest := checkForValidJumpDestination(nJumpDest, validJumpDestinations)

			if nJumpDest > len(byteCode) || !isValidJumpDest {
				success = false
				stack = stack[1:]
				break
			} else {
				pc = nJumpDest
				stack = stack[1:]
				continue
			}

		} else if op == 87 { // JUMPI
			// check for valid JUMPDEST
			if len(validJumpDestinations) == 0 {
				validJumpDestinations = getValidJumpDestinations(byteCode, pc)
			}
			nJumpDest := int(stack[0].Int64())
			isValidJumpDest := checkForValidJumpDestination(nJumpDest, validJumpDestinations)
			if nJumpDest > len(byteCode) || !isValidJumpDest {
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
		} else if op == 91 { // JUMPDEST
			successFlag = true
		} else if op == 82 { // MSTORE
			stack, memory, successFlag = mstore(stack, memory)
		} else if op == 81 { // MLOAD
			stack, memory, successFlag = mload(stack, memory)
		} else if op == 83 { // MSTORE8
			stack, memory, successFlag = mstore8(stack, memory)
		} else if op == 89 { // MSIZE
			stack, successFlag = msize(stack, memory)
		} else if op == 32 { // KECCAK256
			stack, successFlag = keccak256(stack, memory)
		} else if op == 48 { // ADDRESS
			stack, successFlag = address(stack, tx)
		} else if op == 51 { // msg.sender
			stack, successFlag = caller(stack, tx)
		} else if op == 50 { // tx.origin
			stack, successFlag = origin(stack, tx)
		} else if op == 58 { // GASPRICE
			stack, successFlag = gasprice(stack, tx)
		} else if op == 72 {
			stack, successFlag = basefee(stack, block)
		} else if op == 65 {
			stack, successFlag = coinbase(stack, block)
		} else if op == 66 {
			stack, successFlag = timestamp(stack, block)
		} else if op == 67 {
			stack, successFlag = number(stack, block)
		} else if op == 68 {
			stack, successFlag = difficulty(stack, block)
		} else if op == 69 {
			stack, successFlag = gaslimit(stack, block)
		} else if op == 70 {
			stack, successFlag = chainid(stack, block)
		} else if op == 64 {
			stack, successFlag = blockhash(stack)
		} else {
			fmt.Println("******** op *******", op)
			fmt.Println("********* byteCode ******", byteCode)
			fmt.Println("Stack:", stack)
		}
		pc += n + 1
		success = successFlag && success
	}
	return stack, success
}
