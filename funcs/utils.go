package funcs

import (
	"math/rand"
	"reflect"
	"strings"
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

func TrimStringFields(model interface{}) {
	// Get the value and type of the model
	val := reflect.ValueOf(model).Elem()

	// Iterate through all fields in the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is a string and can be modified
		if field.Kind() == reflect.String && field.CanSet() {
			trimmedValue := strings.TrimSpace(field.String()) // Trim spaces
			field.SetString(trimmedValue)                     // Set the trimmed value
		}
	}
}
func DefaultUUID() string {
	return "00000000-0000-0000-0000-000000000000"
}
