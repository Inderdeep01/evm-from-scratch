package evm

func checkForValidJumpDestination(jumpDest int, validJumpDestinations []int) bool {
	for i := 0; i < len(validJumpDestinations); i++ {
		if validJumpDestinations[i] == jumpDest {
			return true
		}
	}
	return false
}

func getValidJumpDestinations(byteCode []byte, pc int) []int {
	var validJumpDestinations []int
	for i := pc + 1; i < len(byteCode); i++ {
		currOp := byteCode[i]
		if currOp >= 95 && currOp <= 127 {
			n := int(currOp - 95)
			i += n
		} else if currOp == 91 {
			validJumpDestinations = append(validJumpDestinations, i)
		}
	}
	return validJumpDestinations
}
