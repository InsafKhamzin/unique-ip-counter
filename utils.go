package main

import (
	"bytes"
	"errors"
)

func IpToIntArray(ip []byte) []uint8 {
	ip = bytes.TrimSpace(ip)
	var result []uint8
	pointer := 0
	for i := 0; i < len(ip); i++ {
		if ip[i] == '.' || i == len(ip)-1 {
			octetSegment := ip[pointer:i]
			if i == len(ip)-1 {
				octetSegment = ip[pointer:]
			}
			var val uint8 = 0
			for i := 0; i < len(octetSegment); i++ {
				// Subtracting the ASCII value, getting digit
				digit := uint8(octetSegment[i] - '0')
				val = val*10 + digit
			}
			result = append(result, val)
			pointer = i + 1
		}
	}
	return result
}

// last three bytes to int
func LastThreeBytesToInt(items []uint8) (uint, error) {
	if len(items) != 3 {
		return 0, errors.New("must be length of 3")
	}
	var result uint
	for i, v := range items {
		result |= uint(v) << (16 - uint(i)*8)
	}
	return result, nil
}
