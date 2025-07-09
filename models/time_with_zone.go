package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TimeWithZone is a custom time type that automatically converts UTC to +07:00 timezone
type TimeWithZone struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for TimeWithZone
func (t *TimeWithZone) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	str := strings.Trim(string(data), `"`)

	// If empty, set to zero time
	if str == "" || str == "null" {
		t.Time = time.Time{}
		return nil
	}

	// Try to parse the time string
	var parsedTime time.Time
	var err error

	// Try different time formats
	formats := []string{
		"2006-01-02T15:04:05Z",          // 2025-03-26T08:00:00Z
		"2006-01-02T15:04:05.000Z",      // 2025-03-26T08:00:00.000Z
		"2006-01-02T15:04:05-07:00",     // 2025-03-26T08:00:00+07:00
		"2006-01-02T15:04:05.000-07:00", // 2025-03-26T08:00:00.000+07:00
		"2006-01-02 15:04:05",           // 2025-03-26 08:00:00
		"2006-01-02",                    // 2025-03-26
	}

	for _, format := range formats {
		parsedTime, err = time.Parse(format, str)
		if err == nil {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}

	// If the time is in UTC (ends with Z), treat it as if it's already in +07:00 timezone
	// This means 2025-03-26T08:00:00Z should become 2025-03-26T08:00:00+07:00
	if strings.HasSuffix(str, "Z") {
		// Create location for +07:00 timezone
		loc := time.FixedZone("Asia/Bangkok", 7*60*60)
		// Use the same time value but in +07:00 timezone
		parsedTime = time.Date(
			parsedTime.Year(),
			parsedTime.Month(),
			parsedTime.Day(),
			parsedTime.Hour(),
			parsedTime.Minute(),
			parsedTime.Second(),
			parsedTime.Nanosecond(),
			loc,
		)
	}

	t.Time = parsedTime
	return nil
}

// MarshalJSON implements custom JSON marshaling for TimeWithZone
func (t TimeWithZone) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	// Convert the time to Asia/Bangkok local time for JSON output
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// fallback to fixed zone if loading fails
		loc = time.FixedZone("Asia/Bangkok", 7*60*60)
	}
	convertedTime := t.Time.In(loc)

	// Format as RFC3339 with +07:00 timezone
	formatted := convertedTime.Format("2006-01-02T15:04:05-07:00")
	return json.Marshal(formatted)
}

// Value implements the driver.Valuer interface for GORM
func (t TimeWithZone) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan implements the sql.Scanner interface for GORM
func (t *TimeWithZone) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case []byte:
		// Try to parse as string
		str := string(v)
		if str == "" {
			t.Time = time.Time{}
			return nil
		}

		// Try different time formats
		formats := []string{
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02T15:04:05-07:00",
			"2006-01-02T15:04:05.000-07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if parsedTime, err := time.Parse(format, str); err == nil {
				t.Time = parsedTime
				return nil
			}
		}

		return fmt.Errorf("cannot scan %v into TimeWithZone", value)
	default:
		return fmt.Errorf("cannot scan %v into TimeWithZone", value)
	}
}

// ConvertUTCToTimezone converts a UTC datetime string to the specified timezone format
// Note: This treats the UTC time as if it's already in the target timezone
// Example: "2025-03-26T08:00:00Z" -> "2025-03-26T08:00:00+07:00"
func ConvertUTCToTimezone(utcString string, timezoneOffset int) (string, error) {
	// Parse the UTC time
	utcTime, err := time.Parse("2006-01-02T15:04:05Z", utcString)
	if err != nil {
		return "", fmt.Errorf("failed to parse UTC time: %v", err)
	}

	// Create location with the specified offset
	loc := time.FixedZone("Custom", timezoneOffset*60*60)

	// Treat the UTC time as if it's already in the target timezone
	localTime := time.Date(
		utcTime.Year(),
		utcTime.Month(),
		utcTime.Day(),
		utcTime.Hour(),
		utcTime.Minute(),
		utcTime.Second(),
		utcTime.Nanosecond(),
		loc,
	)

	// Format as RFC3339 with timezone offset
	formatted := localTime.Format("2006-01-02T15:04:05-07:00")

	return formatted, nil
}

// ConvertUTCToBangkokTime converts a UTC datetime string to Bangkok timezone (+07:00)
// Note: This treats the UTC time as if it's already in Bangkok timezone
// Example: "2025-03-26T08:00:00Z" -> "2025-03-26T08:00:00+07:00"
func ConvertUTCToBangkokTime(utcString string) (string, error) {
	return ConvertUTCToTimezone(utcString, 7)
}
