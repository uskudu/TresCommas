package random

import (
	"crypto/rand"
)

func NewRandomString(length int) string {
	return rand.Text()[:length]
}
