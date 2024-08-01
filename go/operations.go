package evm

import (
	"encoding/hex"
	"fmt"
	"github.com/holiman/uint256"
	"math/big"
)

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

func mstore(stack []*big.Int, memory map[int64]byte) ([]*big.Int, map[int64]byte, bool) {
	if len(stack) < 2 {
		return stack, memory, false
	}
	offset := stack[0].Int64()
	bytes := stack[1].Bytes()
	for i, j := 0, offset; i < len(bytes); i, j = i+1, j+1 {
		memory[j] = bytes[i]
	}
	stack = stack[2:]
	return stack, memory, true
}

func mstore8(stack []*big.Int, memory map[int64]byte) ([]*big.Int, map[int64]byte, bool) {
	if len(stack) < 2 {
		return stack, memory, false
	}
	offset := stack[0].Int64()
	bytes := stack[1].Bytes()

	memory[offset] = bytes[0]
	stack = stack[2:]
	return stack, memory, true
}

func mload(stack []*big.Int, memory map[int64]byte) ([]*big.Int, map[int64]byte, bool) {
	if len(stack) < 1 {
		return stack, memory, false
	}
	offset := stack[0].Int64()
	x := new(big.Int)
	bytes := make([]byte, 32)
	for i, j := 0, offset; i < len(bytes); i, j = i+1, j+1 {
		if j == 32 { //
			fmt.Println("breaking at", j)
			break
		}
		bytes[i] = memory[j]
	}
	fmt.Println("memory:", memory)
	fmt.Println("bytes:", bytes)
	x.SetBytes(bytes)
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[1:]...)
	return stack, memory, true
}
