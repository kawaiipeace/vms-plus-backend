package funcs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func DefaultUUID() string {
	return "00000000-0000-0000-0000-000000000000"
}

func CalculateAge(date time.Time) string {
	today := time.Now()
	years := today.Year() - date.Year()
	months := today.Month() - date.Month()

	// Adjust if birthday hasn't occurred yet this year
	if today.Day() < date.Day() {
		months--
	}

	if months < 0 {
		years--
		months += 12
	}

	return fmt.Sprintf("%d ปี %d เดือน", years, months)
}

func GetEmpImage(empID string) string {
	return fmt.Sprintf("https://pictureapi.pea.co.th/MyphotoAPI/api/v1/Main/GetPicImg?EmpCode=%s&Type=2&SType=2", empID)
}

func GetDuration(createdAt time.Time) string {
	duration := time.Since(createdAt)

	if duration.Hours() < 24 {
		if duration.Hours() < 1 {
			minutes := int(duration.Minutes())
			return fmt.Sprintf("%d นาทีที่แล้ว", minutes)
		}
		hours := int(duration.Hours())
		return fmt.Sprintf("%d ชั่วโมงที่แล้ว", hours)
	}

	days := int(duration.Hours() / 24)
	return fmt.Sprintf("%d วันที่แล้ว", days)
}

// ParseCSV parses a CSV file and returns a slice of maps where each map represents a row with column names as keys.
func ParseCSV(reader io.Reader) ([]map[string]string, error) {
	csvReader := csv.NewReader(reader)
	headers, err := csvReader.Read()
	if err != nil {
		return nil, errors.New("failed to read CSV headers")
	}

	var records []map[string]string
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("failed to read CSV row")
		}

		record := make(map[string]string)
		for i, header := range headers {
			record[header] = row[i]
		}
		records = append(records, record)
	}

	return records, nil
}
