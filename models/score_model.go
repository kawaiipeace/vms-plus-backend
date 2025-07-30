package models

import "fmt"

type Score float64

func (s Score) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.2f\"", s)), nil
}
