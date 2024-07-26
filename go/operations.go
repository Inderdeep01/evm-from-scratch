package evm

import (
	"fmt"
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

func add(stack []*big.Int) ([]*big.Int, int, bool) {
	if len(stack) < 2 {
		return stack, 1, false
	}
	x := new(big.Int)
	fmt.Println(stack[1])
	x.Add(stack[0], stack[1])
	var tempStack []*big.Int
	tempStack = append(tempStack, x)
	stack = append(tempStack, stack[2:]...)
	return stack, 1, true
}
