package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// HaltExecution is the Flag variable to halt execution
var HaltExecution bool

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
func Executor(byteCode []byte, tx Tx, block Block, state State, isStaticCall bool) ([]*big.Int, []Log, string, bool) {
	// TODO: Disable all the disallowed operations in a STATICCALL
	var stack []*big.Int
	pc := 0 // The Program Counter
	var success = true
	// Valid Jump Destinations
	var validJumpDestinations []int
	// Memory of the EVM - Volatile
	// Need to initialize maps in go; attempts to write to a nil map will cause a runtime panic
	var memory []byte
	// Initialize Logs
	var logs []Log
	// Initialize the return string
	var returnValue string
	// Initialize the empty returnDataSize
	var returnDataSize int
	var lastCallReturnData []byte

	if state == nil {
		state = State{}
	}

	for pc < len(byteCode) && !HaltExecution {
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
		} else if op == 49 {
			stack, successFlag = balance(stack, state)
		} else if op == 52 {
			stack, successFlag = callvalue(stack, tx)
		} else if op == 53 {
			stack, successFlag = calldataload(stack, tx)
		} else if op == 54 {
			stack, successFlag = calldatasize(stack, tx)
		} else if op == 55 {
			stack, memory, successFlag = calldatacopy(stack, tx, memory)
		} else if op == 56 {
			stack, successFlag = codesize(stack, byteCode)
		} else if op == 57 {
			stack, memory, successFlag = codecopy(stack, memory, byteCode)
		} else if op == 59 {
			stack, successFlag = extcodesize(stack, state)
		} else if op == 60 {
			stack, memory, successFlag = extcodecopy(stack, memory, state)
		} else if op == 63 {
			stack, successFlag = extcodehash(stack, state)
		} else if op == 71 {
			stack, successFlag = selfbalance(stack, state)
		} else if op == 84 {
			stack, successFlag = sload(stack)
		} else if op == 85 {
			stack, successFlag = sstore(stack, isStaticCall)
		} else if op >= 160 && op <= 164 {
			stack, logs, successFlag = logN(stack, memory, logs, tx, op)
		} else if op == 243 {
			stack, memory, returnValue, successFlag = returnvalue(stack, memory)
		} else if op == 253 {
			stack, memory, returnValue, successFlag = revert(stack, memory)
		} else if op == 241 {
			if len(stack) >= 7 {
				var acc Account
				for _, val := range state {
					acc = val
				}
				retOffset := int(stack[5].Int64())
				retSize := int(stack[6].Int64())
				subByteCode, err := hex.DecodeString(acc.Code.Bin)
				if err != nil {
					successFlag = false
					fmt.Println(err)
				}
				_, _, subValue, subSuccessFlag := Executor(subByteCode, Tx{From: tx.To}, Block{}, state, false)
				bytes, err := hex.DecodeString(subValue)
				if err != nil {
					successFlag = false
					fmt.Println(err)
				}
				returnDataSize = len(bytes)
				lastCallReturnData = bytes
				memory = resizeMemoryIfRequired(memory, retOffset, retSize)
				for i, j := 0, retOffset; i < len(bytes); i, j = i+1, j+1 {
					memory[j] = bytes[i]
				}
				if subSuccessFlag {
					stack = append([]*big.Int{big.NewInt(1)}, stack[7:]...)
				} else {
					stack = append([]*big.Int{big.NewInt(0)}, stack[7:]...)
				}
			}
			successFlag = true
		} else if op == 61 {
			successFlag = true
			stack = append([]*big.Int{big.NewInt(int64(returnDataSize))}, stack...)
		} else if op == 62 {
			stack, memory, successFlag = returnDataCopy(stack, memory, lastCallReturnData)
		} else if op == 244 {
			stack, memory, lastCallReturnData, successFlag = delegateCall(stack, memory, tx, block, state)
		} else if op == 250 {
			stack, memory, lastCallReturnData, successFlag = staticCall(stack, memory, tx, block, state)
		} else if op == 240 {
			stack, memory, state, successFlag = create(stack, memory, tx, state, isStaticCall)
		} else if op == 255 {
			stack, state, successFlag = selfDestruct(stack, state)
		} else {
			// Do Nothing
		}
		pc += n + 1
		success = successFlag && success
	}
	return stack, logs, returnValue, success
}
