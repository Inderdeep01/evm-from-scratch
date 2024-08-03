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

func resizeMemoryIfRequired(memory []byte, startingOffset int, length int) []byte {
	timesToExpand := int((startingOffset + length) / 32)
	newMemorySize := 32 * (timesToExpand + 1)
	if newMemorySize > len(memory) {
		var tempMemory = make([]byte, 32*(timesToExpand+1))
		for i, element := range memory {
			tempMemory[i] = element
		}
		memory = tempMemory
	}

	return memory
}

func checkAndConvertToValidHexString(hexString string) string {
	if hexString != "" && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}
	if len(hexString)%2 != 0 {
		// Strip the "0x" from the input string and add a preceding "0"
		hexString = "0" + hexString
	}

	return hexString
}
