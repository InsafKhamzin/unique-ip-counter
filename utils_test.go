package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_IpToIntArray(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  []uint8
	}{
		{"Success new line", "192.122.44.1\n", []uint8{192, 122, 44, 1}},
		{"Success new line", "192.122.44.1", []uint8{192, 122, 44, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IpToIntArray([]byte(tt.input))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUtils_LastThreeBytesToInt(t *testing.T) {
	var tests = []struct {
		name  string
		input []uint8
		want  uint
		err   error
	}{
		{"Invalid length", []uint8{128, 64, 32, 16}, 0, errors.New("must be length of 3")},
		{"Success", []uint8{64, 32, 16}, 4202512, nil}, //01000000 00100000 00010000
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LastThreeBytesToInt(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
