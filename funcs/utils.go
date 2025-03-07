package funcs

import (
	"math/rand"
	"time"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomRefCode(n int) string {
	src := rand.NewSource(time.Now().UnixNano()) // Create a new random source
	r := rand.New(src)                           // Use a new rand instance

	refCode := make([]byte, n)
	for i := range refCode {
		refCode[i] = letters[r.Intn(len(letters))]
	}
	return string(refCode)
}
