package spi

// setBit is used to set a bit at a given position to 1.
func setBit(n uint8, pos uint8) uint8 {
	n |= (1 << pos)
	return n
}

// clearBit is used to set a bit at a given position to 0.
func clearBit(n uint8, pos uint8) uint8 {
	mask := ^uint8(1 << pos)
	n &= mask
	return n
}

func twosComplement16Bit(uValue uint16) int16 {
	result := int32(uValue)
	if result&0x8000 != 0 {
		result -= 1 << 16
	}
	return int16(result)
}

func swapBytes(value uint16) uint16 {
	return (value << 8) | (value >> 8)
}
