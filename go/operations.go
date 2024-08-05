package evm

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"math/big"
)

const thirtyTwoBytes = 32

func pushN(code []byte) ([]*big.Int, int, bool) {
	var stack []*big.Int
	n := int(code[0] - 95)
	ops := 0
	if n == 0 {
		stack = append(stack, big.NewInt(0))
		ops++
	} else {
		x := new(big.Int)
		x.SetBytes(code[1 : n+1])
		stack = append(stack, x)
		ops += n
	}

	return stack, ops, true
}

//func stop(stack []*big.Int) ([]*big.Int, int, bool){
//	return stack, 1, true
//}
//
//func pop(stack []*big.Int) ([]*big.Int, int, bool){
//	return stack[1:], 1, true
//}

func checkForOverflowUnderflow(x *big.Int) *big.Int {
	// check for overflow
	bytes, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	y := new(big.Int)
	y.SetBytes(bytes)
	compare := x.Cmp(y)
	if compare == 1 {
		z := big.NewInt(1)
		x.Sub(x, y)
		x.Sub(x, z)
	}
	zero := big.NewInt(0)
	compare = x.Cmp(zero)
	if compare == -1 {
		z := big.NewInt(1)
		x.Add(x, y)
		x.Add(x, z)
	}
	return x
}

func add(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Add(stack[0], stack[1])
	x = checkForOverflowUnderflow(x)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func multiply(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Mul(stack[0], stack[1])
	x = checkForOverflowUnderflow(x)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func sub(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Sub(stack[0], stack[1])
	x = checkForOverflowUnderflow(x)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func div(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[1].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(big.Int)
		x.Div(stack[0], stack[1])
		tempStack = append(tempStack, x)
	}

	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func sdiv(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[1].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(uint256.Int)
		a := new(uint256.Int)
		a.SetFromBig(stack[0])
		b := new(uint256.Int)
		b.SetFromBig(stack[1])
		x.SDiv(a, b)
		tempStack = append(tempStack, x.ToBig())
	}

	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func mod(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[1].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(big.Int)
		x.Mod(stack[0], stack[1])
		tempStack = append(tempStack, x)
	}

	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func smod(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[1].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(uint256.Int)
		a := new(uint256.Int)
		a.SetFromBig(stack[0])
		b := new(uint256.Int)
		b.SetFromBig(stack[1])
		x.SMod(a, b)
		tempStack = append(tempStack, x.ToBig())
	}

	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func addmod(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 3 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[2].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(big.Int)
		x.Add(stack[0], stack[1])
		x.Mod(x, stack[2])
		x = checkForOverflowUnderflow(x)
		tempStack = append(tempStack, x)
	}

	stack = append(tempStack, stack[3:]...)
	return stack, 1, true
}

func mulmod(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 3 {
		return stack, 1, false
	}
	var tempStack []*big.Int

	zero := big.NewInt(0)
	compare := stack[2].Cmp(zero)
	if compare == 0 {
		tempStack = append(tempStack, big.NewInt(0))
	} else {
		x := new(big.Int)
		x.Mul(stack[0], stack[1])
		x.Mod(x, stack[2])
		x = checkForOverflowUnderflow(x)
		tempStack = append(tempStack, x)
	}

	stack = append(tempStack, stack[3:]...)
	return stack, 1, true
}

func exp(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Exp(stack[0], stack[1], big.NewInt(0))
	x = checkForOverflowUnderflow(x)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func signextend(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(uint256.Int)
	x.SetFromBig(stack[0])
	byteNum := new(uint256.Int)
	byteNum.SetFromBig(stack[1])
	x.ExtendSign(byteNum, x)
	res := x.ToBig()
	res = checkForOverflowUnderflow(res)
	var tempStack []*big.Int
	tempStack = append(tempStack, res)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func lt(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	compare := stack[0].Cmp(stack[1])
	if compare == -1 {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func gt(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	compare := stack[0].Cmp(stack[1])
	if compare == 1 {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func slt(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	//x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	isLt := a.Slt(b)
	if isLt {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func sgt(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	//x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	isGt := a.Sgt(b)
	if isGt {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func eq(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	//x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	isEq := a.Eq(b)
	if isEq {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func isZero(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 1 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	//x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(big.NewInt(0))
	isEq := a.Eq(b)
	if isEq {
		tempStack = append(tempStack, big.NewInt(1))
	} else {
		tempStack = append(tempStack, big.NewInt(0))
	}
	stack = append(tempStack, stack[1:]...)
	return stack, 1, true
}

func and(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	x.And(a, b)

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func or(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	x.Or(a, b)

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func xor(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])
	x.Xor(a, b)

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func not(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 1 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])

	x.Not(a)

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[1:]...)
	return stack, 1, true
}

func shl(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[1])

	x.Lsh(a, uint(stack[0].Uint64()))

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func shr(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[1])

	x.Rsh(a, uint(stack[0].Uint64()))

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func sar(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[1])

	x.SRsh(a, uint(stack[0].Uint64()))

	tempStack = append(tempStack, x.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func getByte(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	var tempStack []*big.Int
	//x := new(uint256.Int)
	a := new(uint256.Int)
	a.SetFromBig(stack[0])
	b := new(uint256.Int)
	b.SetFromBig(stack[1])

	b.Byte(a)

	tempStack = append(tempStack, b.ToBig())
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}

func dupN(code byte, stack []*big.Int) ([]*big.Int, int, bool) {
	var tempStack []*big.Int
	n := int(code - 128)
	if len(stack) < n {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Set(stack[n])

	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[0:]...)

	return stack, 0, true
}

func swapN(code byte, stack []*big.Int) ([]*big.Int, int, bool) {
	n := int(code - 144)
	if len(stack) < n {
		return stack, 1, false
	}
	x := new(big.Int)
	x.Set(stack[n+1])
	stack[n+1] = stack[0]
	stack[0] = x

	return stack, 0, true
}

func getGas() *big.Int {
	bytes, _ := hex.DecodeString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	y := new(big.Int)
	y.SetBytes(bytes)
	return y
}

func mstore(stack []*big.Int, memory []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 2 {
		return stack, memory, false
	}
	offset := int(stack[0].Int64())
	bytes := stack[1].Bytes()
	bytesToStore := make([]byte, thirtyTwoBytes)
	// ! IMPORTANT !
	// Need to prepend '0' byte(s) for numbers that consume less than 32 bytes to be stored 'as-is'
	// For Example, 1 should be represented in bytes as [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1]
	// instead of [1], because it will be misinterpreted while performing subsequent access
	for i, j := len(bytes)-1, 31; j >= 0; i, j = i-1, j-1 {
		if i < 0 {
			break
		}
		bytesToStore[j] = bytes[i]
	}
	// store the bytes to memory
	memory = resizeMemoryIfRequired(memory, offset, thirtyTwoBytes)
	for i, j := 0, offset; i < thirtyTwoBytes; i, j = i+1, j+1 {
		memory[j] = bytesToStore[i]
	}
	stack = stack[2:]
	return stack, memory, true
}

func mstore8(stack []*big.Int, memory []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 2 {
		return stack, memory, false
	}
	offset := int(stack[0].Int64())
	bytes := stack[1].Bytes()
	memory = resizeMemoryIfRequired(memory, offset, 0)
	memory[offset] = bytes[0]
	stack = stack[2:]
	return stack, memory, true
}

func mload(stack []*big.Int, memory []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 1 {
		return stack, memory, false
	}
	offset := int(stack[0].Int64())
	memory = resizeMemoryIfRequired(memory, offset, 31)
	x := new(big.Int)
	bytes := make([]byte, thirtyTwoBytes)
	for i, j := 0, offset; i < thirtyTwoBytes; i, j = i+1, j+1 {
		if j == len(memory) { //
			break
		}
		bytes[i] = memory[j]
	}
	x.SetBytes(bytes)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[1:]...)
	return stack, memory, true
}

func msize(stack []*big.Int, memory []byte) ([]*big.Int, bool) {
	x := big.NewInt(int64(len(memory)))
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func keccak256(stack []*big.Int, memory []byte) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return stack, false
	}
	offset := int(stack[0].Int64())
	bytes := make([]byte, stack[1].Int64())
	for i, j := 0, offset; i < len(bytes); i, j = i+1, j+1 {
		if j == len(memory) { //
			break
		}
		bytes[i] = memory[j]
	}
	x := new(big.Int)
	hash := crypto.Keccak256(bytes)
	x.SetBytes(hash)
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[2:]...)
	return stack, true
}

func address(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(tx.To))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	addr := new(big.Int)
	addr.SetBytes(bytes)
	tempStack := []*big.Int{addr}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func caller(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(tx.From))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	sender := new(big.Int)
	sender.SetBytes(bytes)
	tempStack := []*big.Int{sender}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func origin(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(tx.Origin))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	txOrigin := new(big.Int)
	txOrigin.SetBytes(bytes)
	tempStack := []*big.Int{txOrigin}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func gasprice(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(tx.GasPrice))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	gasPrice := new(big.Int)
	gasPrice.SetBytes(bytes)
	tempStack := []*big.Int{gasPrice}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func basefee(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.BaseFee))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	baseFee := new(big.Int)
	baseFee.SetBytes(bytes)
	tempStack := []*big.Int{baseFee}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func coinbase(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.CoinBase))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func timestamp(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.Timestamp))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func number(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.Number))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func difficulty(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.Difficulty))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func gaslimit(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.GasLimit))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func chainid(stack []*big.Int, block Block) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(block.ChainID))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	coinBase := new(big.Int)
	coinBase.SetBytes(bytes)
	tempStack := []*big.Int{coinBase}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func blockhash(stack []*big.Int) ([]*big.Int, bool) {
	x := big.NewInt(0)
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[1:]...)
	return stack, true
}

func balance(stack []*big.Int, state State) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return stack, false
	}
	addr := "0x" + stack[0].Text(16)
	bal := state[addr].Balance
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(bal))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	x := new(big.Int)
	x.SetBytes(bytes)
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[1:]...)
	return stack, true
}

func callvalue(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(tx.Value))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	gasPrice := new(big.Int)
	gasPrice.SetBytes(bytes)
	tempStack := []*big.Int{gasPrice}
	stack = append(tempStack, stack[:]...)
	return stack, true
}

func calldataload(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return stack, false
	}
	offset := int(stack[0].Int64())
	x := new(big.Int)
	bytes := make([]byte, thirtyTwoBytes)
	data, err := hex.DecodeString(checkAndConvertToValidHexString(tx.Data))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	for i, j := 0, offset; i < len(bytes); i, j = i+1, j+1 {
		if j == len(data) { //
			break
		}
		bytes[i] = data[j]
	}
	x.SetBytes(bytes)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[1:]...)
	return stack, true
}

func calldatasize(stack []*big.Int, tx Tx) ([]*big.Int, bool) {
	data, err := hex.DecodeString(checkAndConvertToValidHexString(tx.Data))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	size := big.NewInt(int64(len(data)))
	var tempStack []*big.Int
	tempStack = append(tempStack, size)
	stack = append(tempStack, stack...)
	return stack, true
}

func calldatacopy(stack []*big.Int, tx Tx, memory []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 3 {
		return stack, memory, false
	}
	destOffset := int(stack[0].Int64())
	calldataOffset := int(stack[1].Int64())
	size := int(stack[2].Int64())
	// loading data from calldata
	bytes := make([]byte, thirtyTwoBytes)
	data, err := hex.DecodeString(checkAndConvertToValidHexString(tx.Data))
	if err != nil {
		fmt.Println(err)
		return stack, memory, false
	}
	for i, j := 0, calldataOffset; i < size; i, j = i+1, j+1 {
		if j == len(data) { //
			break
		}
		bytes[i] = data[j]
	}
	memory = resizeMemoryIfRequired(memory, destOffset, len(bytes))
	for i, j := 0, destOffset; i < len(bytes); i, j = i+1, j+1 {
		memory[j] = bytes[i]
	}
	return stack[3:], memory, true
}

func codesize(stack []*big.Int, byteCode []byte) ([]*big.Int, bool) {
	x := big.NewInt(int64(len(byteCode)))
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack...)
	return stack, true
}

func codecopy(stack []*big.Int, memory []byte, byteCode []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 3 {
		return stack, memory, false
	}
	destOffset := int(stack[0].Int64())
	offset := int(stack[1].Int64())
	size := int(stack[2].Int64())
	// loading data from calldata
	bytes := make([]byte, size)
	for i, j := 0, offset; i < size; i, j = i+1, j+1 {
		if j == len(byteCode) { //
			break
		}
		bytes[i] = byteCode[j]
	}
	memory = resizeMemoryIfRequired(memory, destOffset, len(bytes))
	for i, j := 0, destOffset; i < len(bytes); i, j = i+1, j+1 {
		memory[j] = bytes[i]
	}
	return stack[3:], memory, true
}

func extcodesize(stack []*big.Int, state State) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return stack, false
	}
	addr := "0x" + stack[0].Text(16)
	code := state[addr].Code
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(code.Bin))
	if err != nil {
		fmt.Println(err)
		return stack, false
	}
	x := big.NewInt(int64(len(bytes)))
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[1:]...)
	return stack, true
}

func extcodecopy(stack []*big.Int, memory []byte, state State) ([]*big.Int, []byte, bool) {
	if len(stack) < 4 {
		return stack, memory, false
	}
	addr := "0x" + stack[0].Text(16)
	code := state[addr].Code
	bytes, err := hex.DecodeString(checkAndConvertToValidHexString(code.Bin))
	if err != nil {
		fmt.Println(err)
		return stack, memory, false
	}
	destOffset := int(stack[1].Int64())
	offset := int(stack[2].Int64())
	size := int(stack[3].Int64())
	// loading data from state
	codeBytes := make([]byte, size)
	for i, j := 0, offset; i < size; i, j = i+1, j+1 {
		if j >= len(bytes) { //
			codeBytes[i] = 0
		} else {
			codeBytes[i] = bytes[j]
		}
	}
	memory = resizeMemoryIfRequired(memory, destOffset, len(codeBytes))
	for i, j := 0, destOffset; i < len(codeBytes); i, j = i+1, j+1 {
		memory[j] = codeBytes[i]
	}
	return stack[4:], memory, true
}

func extcodehash(stack []*big.Int, state State) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return stack, false
	}
	addr := "0x" + stack[0].Text(16)
	account, exists := state[addr]
	if !exists {
		x := big.NewInt(0)
		tempStack := []*big.Int{x}
		stack = append(tempStack, stack[1:]...)
	} else {
		bytes, err := hex.DecodeString(checkAndConvertToValidHexString(account.Code.Bin))
		if err != nil {
			fmt.Println(err)
			return stack, false
		}
		x := new(big.Int)
		hash := crypto.Keccak256(bytes)
		x.SetBytes(hash)
		tempStack := []*big.Int{x}
		stack = append(tempStack, stack[1:]...)
	}
	return stack, true
}

func selfbalance(stack []*big.Int, state State) ([]*big.Int, bool) {
	x := new(big.Int)
	for _, val := range state {
		bytes, _ := hex.DecodeString(checkAndConvertToValidHexString(val.Balance))
		x.SetBytes(bytes)
	}
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack...)
	return stack, true
}

func sstore(stack []*big.Int, canNotWriteToStorage bool) ([]*big.Int, bool) {
	if len(stack) < 2 || canNotWriteToStorage {
		return stack, false
	}
	key := stack[0].Text(16)
	val := stack[1].Bytes()
	Storage[key] = val
	stack = stack[2:]
	return stack, true
}

func sload(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return stack, false
	}
	key := stack[0].Text(16)
	val := Storage[key]
	x := new(big.Int)
	x.SetBytes(val)
	tempStack := []*big.Int{x}
	stack = append(tempStack, stack[1:]...)
	return stack, true
}

func logN(stack []*big.Int, memory []byte, logs []Log, tx Tx, opcode byte) ([]*big.Int, []Log, bool) {
	// get the number of logs to append
	n := int(opcode - 160)
	// validate the stack length
	if len(stack) < 2+n {
		return stack, logs, false
	}
	offset := int(stack[0].Int64())
	size := int(stack[1].Int64())
	memory = resizeMemoryIfRequired(memory, offset, size-1)
	// x is the value to push to stack
	x := new(big.Int)
	// get the data to return with logs from memory and store in 'bytes'
	bytes := make([]byte, size)
	for i, j := 0, offset; i < len(bytes); i, j = i+1, j+1 {
		if j == len(memory) {
			break
		}
		bytes[i] = memory[j]
	}
	x.SetBytes(bytes)
	// List of topics to emit
	var topics []string

	// check how many topics are to be added
	if n > 0 {
		for i := 1; i <= n; i++ {
			// we need to skip 0th and 1st element of stack as they are the 'offset' and 'size' values respectively
			topic := getFormattedHexString(stack[i+1].Text(16))
			topics = append(topics, topic)
		}
	}
	logs = append(logs, Log{tx.To, checkAndConvertToValidHexString(x.Text(16)), topics})
	stack = stack[n+2:] // popping the stack
	return stack, logs, true
}

func returnvalue(stack []*big.Int, memory []byte) ([]*big.Int, []byte, string, bool) {
	if len(stack) < 2 {
		return stack, memory, "", false
	}
	offset := int(stack[0].Int64())
	size := int(stack[1].Int64())
	bytes := make([]byte, size)
	memory = resizeMemoryIfRequired(memory, offset, size)
	for i, j := 0, offset; i < size; i, j = i+1, j+1 {
		bytes[i] = memory[j]
	}
	str := checkAndConvertToValidHexString(hex.EncodeToString(bytes))
	stack = stack[2:]
	return stack, memory, str, true
}

func revert(stack []*big.Int, memory []byte) ([]*big.Int, []byte, string, bool) {
	if len(stack) < 2 {
		return stack, memory, "", false
	}
	offset := int(stack[0].Int64())
	size := int(stack[1].Int64())
	bytes := make([]byte, size)
	memory = resizeMemoryIfRequired(memory, offset, size)
	for i, j := 0, offset; i < size; i, j = i+1, j+1 {
		bytes[i] = memory[j]
	}
	str := checkAndConvertToValidHexString(hex.EncodeToString(bytes))
	stack = stack[2:]
	return stack, memory, str, false
}

func returnDataCopy(stack []*big.Int, memory []byte, lastCallData []byte) ([]*big.Int, []byte, bool) {
	if len(stack) < 3 {
		return stack, memory, false
	}
	destOffset := int(stack[0].Int64())
	offset := int(stack[1].Int64())
	size := int(stack[2].Int64())
	if (offset + size) > len(lastCallData) {
		return stack, memory, false
	}
	// copy the desired data from previous return value to memory
	memory = resizeMemoryIfRequired(memory, destOffset, size)
	for i, j := offset, destOffset; i < len(lastCallData); i, j = i+1, j+1 {
		memory[j] = lastCallData[i]
	}
	return stack[3:], memory, true
}

func delegateCall(stack []*big.Int, memory []byte, tx Tx, block Block, state State) ([]*big.Int, []byte, []byte, bool) {
	if len(stack) < 6 {
		return stack[6:], memory, nil, false
	}
	addr := getFormattedHexString(stack[1].Text(16))
	retOffset := int(stack[4].Int64())
	retSize := int(stack[5].Int64())
	acc, exists := state[addr]
	if !exists {
		return stack[6:], memory, nil, false
	}
	subByteCode, err := hex.DecodeString(acc.Code.Bin)
	if err != nil {
		fmt.Println(err)
		return stack[6:], memory, nil, false
	}
	_, _, subValue, subSuccess := Executor(subByteCode, tx, block, state, false)
	bytes, err := hex.DecodeString(subValue)
	if err != nil {
		fmt.Println(err)
		return stack[6:], memory, nil, false
	}
	memory = resizeMemoryIfRequired(memory, retOffset, retSize)
	for i, j := 0, retOffset; i < len(bytes); i, j = i+1, j+1 {
		memory[j] = bytes[i]
	}
	if subSuccess {
		stack = append([]*big.Int{big.NewInt(1)}, stack[6:]...)
	} else {
		stack = append([]*big.Int{big.NewInt(0)}, stack[6:]...)
	}
	return stack, memory, bytes, true
}

func staticCall(stack []*big.Int, memory []byte, tx Tx, block Block, state State) ([]*big.Int, []byte, []byte, bool) {
	if len(stack) < 6 {
		return stack[6:], memory, nil, false
	}
	addr := getFormattedHexString(stack[1].Text(16))
	retOffset := int(stack[4].Int64())
	retSize := int(stack[5].Int64())
	acc, exists := state[addr]
	if !exists {
		return stack[6:], memory, nil, false
	}
	subByteCode, err := hex.DecodeString(acc.Code.Bin)
	if err != nil {
		fmt.Println(err)
		return stack[6:], memory, nil, false
	}
	_, _, subValue, subSuccess := Executor(subByteCode, tx, block, state, true)
	bytes, err := hex.DecodeString(subValue)
	if err != nil {
		fmt.Println(err)
		return stack[6:], memory, nil, false
	}
	memory = resizeMemoryIfRequired(memory, retOffset, retSize)
	for i, j := 0, retOffset; i < len(bytes); i, j = i+1, j+1 {
		memory[j] = bytes[i]
	}
	if subSuccess {
		stack = append([]*big.Int{big.NewInt(1)}, stack[6:]...)
	} else {
		stack = append([]*big.Int{big.NewInt(0)}, stack[6:]...)
	}
	return stack, memory, bytes, true
}

func create(stack []*big.Int, memory []byte, tx Tx, state State, isStaticCall bool) ([]*big.Int, []byte, State, bool) {
	if len(stack) < 3 || isStaticCall {
		return stack, memory, state, false
	}
	value := stack[0]
	offset := int(stack[1].Int64())
	size := int(stack[2].Int64())
	bytes := make([]byte, size)
	memory = resizeMemoryIfRequired(memory, offset, size)
	for i, j := 0, offset; i < size; i, j = i+1, j+1 {
		bytes[i] = memory[j]
	}
	addr := tx.To
	_, alreadyExists := state[addr]
	if alreadyExists {
		return stack[3:], memory, state, false
	}
	_, _, subValue, subSuccess := Executor(bytes, Tx{}, Block{}, nil, false)
	if subSuccess {

		acc := Account{
			Balance: getFormattedHexString(value.Text(16)),
			Code:    Code{Bin: checkAndConvertToValidHexString(subValue)},
		}
		state[addr] = acc
		addrBytes, err := hex.DecodeString(checkAndConvertToValidHexString(addr))
		if err != nil {
			fmt.Println(err)
			return stack[3:], memory, state, false
		}
		x := new(big.Int)
		x.SetBytes(addrBytes)
		stack = append([]*big.Int{x}, stack[3:]...)
	} else {
		x := big.NewInt(0)
		stack = append([]*big.Int{x}, stack[3:]...)
	}
	return stack, memory, state, true
}

func selfDestruct(stack []*big.Int, state State) ([]*big.Int, State, bool) {
	if len(stack) < 1 {
		return stack, state, false
	}
	addressToSendValue := getFormattedHexString(stack[0].Text(16))
	var bal string
	var selfAddr string
	for addr, acc := range state {
		bal = acc.Balance
		selfAddr = addr
	}
	account, exists := state[addressToSendValue]
	if !exists {
		account = Account{
			Balance: bal,
			Code:    Code{},
		}
		state[addressToSendValue] = account
	} else {
		existingBalBytes, _ := hex.DecodeString(checkAndConvertToValidHexString(account.Balance))
		newBalanceBytes, _ := hex.DecodeString(checkAndConvertToValidHexString(bal))
		a, b := new(big.Int), new(big.Int)
		a.SetBytes(existingBalBytes)
		b.SetBytes(newBalanceBytes)
		a.Add(a, b)
		account.Balance = getFormattedHexString(a.Text(16))
	}
	delete(state, selfAddr)
	return stack[1:], state, true
}
