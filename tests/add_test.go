package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			ADD, RAX, RCX,
		},
		output: []uint8{
			0x48, 0x01, 0xc8,
		},
	}}...)
}