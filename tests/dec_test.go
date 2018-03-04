package tests

import . "github.com/modern-go/amd64"

func init() {
	testCases = append(testCases, []testCase{{
		input: input{
			DEC, EAX,
		},
		comment: "dec set opcode extension to 1",
		output: []uint8{
			0xff, 0xc8,
		},
	},{
		input: input{
			DEC, AL,
		},
		output: []uint8{
			0xfe, 0xc8,
		},
	}}...)
}
