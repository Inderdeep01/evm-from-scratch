package evm

import (
	"encoding/hex"
	"math/big"
)

//const maxIntValue Int =

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
