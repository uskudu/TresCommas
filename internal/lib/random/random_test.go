package random

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewRandomString(t *testing.T) {
	tests := []int{-2, 0, 1, 5, 19, 26, 27}

	assert.PanicMatches(t, func() { NewRandomString(tests[0]) }, "runtime error: slice bounds out of range [:-2]")
	assert.IsEqual(NewRandomString(tests[1]), "")
	assert.MatchRegex(t, NewRandomString(tests[2]), "^[A-Za-z0-9]$")
	assert.MatchRegex(t, NewRandomString(tests[3]), "^[A-Za-z0-9]{5}$")
	assert.MatchRegex(t, NewRandomString(tests[4]), "^[A-Za-z0-9]{19}$")
	assert.MatchRegex(t, NewRandomString(tests[5]), "^[A-Za-z0-9]{26}$")
	assert.PanicMatches(t, func() { NewRandomString(tests[6]) }, "runtime error: slice bounds out of range [:27] with length 26")
}
