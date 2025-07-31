package models

import "fmt"

type Score float64

func (s Score) MarshalJSON() ([]byte, error) {
	//if s is 0, return "0"
	if s == 0 {
		return []byte("\"0\""), nil
	}
	return []byte(fmt.Sprintf("\"%.2f\"", s)), nil
}
