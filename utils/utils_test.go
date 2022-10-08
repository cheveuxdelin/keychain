package utils_test

import (
	"fmt"
	"testing"

	"github.com/cheveuxdelin/keychain/utils"
)

func TestIsNumber(t *testing.T) {
	type testStruct struct {
		have rune
		want bool
	}

	for _, tt := range []testStruct{
		{have: '0', want: true},
		{have: '1', want: true},
		{have: '2', want: true},
		{have: '3', want: true},
		{have: '4', want: true},
		{have: '5', want: true},
		{have: '6', want: true},
		{have: '7', want: true},
		{have: '8', want: true},
		{have: '9', want: true},
		{have: 'b', want: false},
		{have: 'c', want: false},
		{have: '{', want: false},
	} {
		testname := fmt.Sprintf("%d", tt.have)
		t.Run(testname, func(t *testing.T) {
			if utils.IsNumber(tt.have) != tt.want {
				t.Error("error")
			}
		})
	}
}
