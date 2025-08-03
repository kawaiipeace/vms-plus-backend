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
	"vms_plus_be/models"
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

func CalculateAgeInt(date time.Time) int {
	now := time.Now()
	// Subtract the registration year from the current year
	age := now.Year() - date.Year()

	// Adjust if the current date is before the registration date in the year
	if now.YearDay() < date.YearDay() {
		age--
	}
	return age
}

func GetEmpImage(empID string) string {
	return fmt.Sprintf("https://pictureapi.pea.co.th/MyphotoAPI/api/v1/Main/GetPicImg?EmpCode=%s&Type=2&SType=2", empID)
}

func GetDuration(createdAt time.Time) string {
	duration := time.Since(createdAt)
	switch {
	case duration.Minutes() < 1 && duration.Hours() < 1:
		return "ตอนนี้"
	case duration.Hours() < 1:
		return fmt.Sprintf("%d นาทีที่แล้ว", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d ชั่วโมงที่แล้ว", int(duration.Hours()))
	default:
		return fmt.Sprintf("%d วันที่แล้ว", int(duration.Hours()/24))
	}
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

func IsHoliday(date time.Time, holidays []models.VmsMasHolidays) bool {
	// Check if the date is a weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return true
	}

	// Check if the date is in the holidays list
	for _, holiday := range holidays {
		if holiday.HolidaysDate.Time.Year() == date.Year() &&
			holiday.HolidaysDate.Time.Month() == date.Month() &&
			holiday.HolidaysDate.Time.Day() == date.Day() {
			return true
		}
	}

	return false
}

func GetDateBuddhistYear(date time.Time) string {
	day := date.Day()
	month := int(date.Month())
	year := date.Year() + 543
	return fmt.Sprintf("%02d/%02d/%04d", day, month, year)
}

func GetDateWithZone(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	thaiYear := date.Year() + 543
	return fmt.Sprintf("%02d-%02d-%04d %02d:%02d:%02d", date.Day(), date.Month(), thaiYear, date.Hour(), date.Minute(), date.Second())
}

func GetReportNumber(number float64) string {
	if number == 0 {
		return ""
	}
	return fmt.Sprintf("%.0f", number)
}
