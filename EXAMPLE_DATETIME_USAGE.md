# Practical Example: Using DateTime Conversion in VMS Plus BE

This document shows how to integrate the new datetime conversion functionality into the existing VMS Plus BE codebase.

## Important Note: Timezone Handling

**The system treats UTC times as if they are already in +07:00 timezone.**

- Input: `"2025-03-26T08:00:00Z"` 
- Output: `"2025-03-26T08:00:00+07:00"`
- **NOT**: `"2025-03-26T15:00:00+07:00"` (no hour conversion)

This means when you send `2025-03-26T08:00:00Z`, it's interpreted as 8:00 AM in Bangkok time (+07:00), not converted from UTC.

## Example 1: Updating a Model to Use TimeWithZone

Let's update the `VehicleReportTripDetail` model in `models/vehicle_management_model.go`:

### Before (Current Implementation)
```go
type VehicleReportTripDetail struct {
    // ... other fields ...
    TripStartDatetime TimeWithZone `gorm:"column:trip_start_datetime;type:timestamp with time zone" json:"trip_start_datetime" example:"2025-03-26T08:00:00+07:00"`
    TripEndDatetime   TimeWithZone `gorm:"column:trip_end_datetime;type:timestamp" json:"trip_end_datetime" example:"2025-03-26T10:00:00+07:00"`
    // ... other fields ...
}
```

### After (With TimeWithZone)
```go
type VehicleReportTripDetail struct {
    // ... other fields ...
    TripStartDatetime funcs.TimeWithZone `gorm:"column:trip_start_datetime;type:timestamp with time zone" json:"trip_start_datetime" example:"2025-03-26T08:00:00+07:00"`
    TripEndDatetime   funcs.TimeWithZone `gorm:"column:trip_end_datetime;type:timestamp" json:"trip_end_datetime" example:"2025-03-26T10:00:00+07:00"`
    // ... other fields ...
}
```

## Example 2: Using in a Handler

Here's how to use it in a handler like `handlers/vehicle_management_handler.go`:

```go
// CreateTripRequest represents the JSON request for creating a trip
type CreateTripRequest struct {
    VehicleUID        string              `json:"vehicle_uid"`
    TripStartDatetime funcs.TimeWithZone  `json:"trip_start_datetime"`
    TripEndDatetime   funcs.TimeWithZone  `json:"trip_end_datetime"`
    Destination       string              `json:"destination"`
}

func (h *VehicleManagementHandler) CreateTrip(c *gin.Context) {
    var request CreateTripRequest
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // The times are automatically converted to +07:00 timezone
    // Input:  {"trip_start_datetime": "2025-03-26T08:00:00Z"}
    // Result: request.TripStartDatetime.Time will be 2025-03-26T08:00:00+07:00
    
    // Access the underlying time.Time value
    startTime := request.TripStartDatetime.Time
    endTime := request.TripEndDatetime.Time
    
    // Use the times in your business logic
    fmt.Printf("Trip start: %s\n", startTime.Format("2006-01-02T15:04:05-07:00"))
    fmt.Printf("Trip end: %s\n", endTime.Format("2006-01-02T15:04:05-07:00"))
    
    // ... rest of your logic
}
```

## Example 3: API Request/Response Flow

### Client sends JSON request:
```json
{
    "vehicle_uid": "123e4567-e89b-12d3-a456-426614174000",
    "trip_start_datetime": "2025-03-26T08:00:00Z",
    "trip_end_datetime": "2025-03-26T10:00:00Z",
    "destination": "Bangkok"
}
```

### Server receives and processes:
```go
// The TimeWithZone automatically converts:
// "2025-03-26T08:00:00Z" -> 2025-03-26T08:00:00+07:00 (same time, different timezone)
// "2025-03-26T10:00:00Z" -> 2025-03-26T10:00:00+07:00 (same time, different timezone)
```

### Server responds with JSON:
```json
{
    "trip_id": "456e7890-e89b-12d3-a456-426614174001",
    "trip_start_datetime": "2025-03-26T08:00:00+07:00",
    "trip_end_datetime": "2025-03-26T10:00:00+07:00",
    "status": "created"
}
```

## Example 4: Using Utility Functions Directly

If you need to convert times manually in your business logic:

```go
func (h *VehicleManagementHandler) ProcessExternalData(c *gin.Context) {
    // External API returns UTC times
    externalData := map[string]string{
        "start_time": "2025-03-26T08:00:00Z",
        "end_time": "2025-03-26T10:00:00Z",
    }
    
    // Convert to Bangkok timezone (treats as if already in +07:00)
    startTime, err := funcs.ConvertUTCToBangkokTime(externalData["start_time"])
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
        return
    }
    
    endTime, err := funcs.ConvertUTCToBangkokTime(externalData["end_time"])
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
        return
    }
    
    fmt.Printf("Converted start time: %s\n", startTime) // 2025-03-26T08:00:00+07:00
    fmt.Printf("Converted end time: %s\n", endTime)     // 2025-03-26T10:00:00+07:00
    
    // ... process the converted times
}
```

## Example 5: Database Operations

When working with GORM and the TimeWithZone type:

```go
// Creating a record
trip := models.VehicleReportTripDetail{
    MasVehicleUID: "123e4567-e89b-12d3-a456-426614174000",
    TripStartDatetime: funcs.TimeWithZone{Time: time.Now()},
    TripEndDatetime: funcs.TimeWithZone{Time: time.Now().Add(2 * time.Hour)},
}

// Save to database
if err := config.DB.Create(&trip).Error; err != nil {
    // Handle error
}

// Query from database
var trips []models.VehicleReportTripDetail
if err := config.DB.Find(&trips).Error; err != nil {
    // Handle error
}

// The times will be automatically formatted as +07:00 when marshaled to JSON
```

## Migration Checklist

To migrate existing code to use the new datetime conversion:

1. **Update imports** (if not already present):
   ```go
   import "vms_plus_be/funcs"
   ```

2. **Replace time.Time with funcs.TimeWithZone** in structs that need automatic conversion:
   ```go
   // Before
   CreatedAt TimeWithZone `json:"created_at"`
   
   // After
   CreatedAt funcs.TimeWithZone `json:"created_at"`
   ```

3. **Update any direct time.Time assignments**:
   ```go
   // Before
   record.CreatedAt = time.Now()
   
   // After
   record.CreatedAt = funcs.TimeWithZone{Time: time.Now()}
   ```

4. **Access underlying time.Time when needed**:
   ```go
   // For database operations or time calculations
   actualTime := record.CreatedAt.Time
   ```

5. **Test the conversion** with sample JSON requests containing UTC times.

## Benefits

- **Automatic Conversion**: No need to manually convert UTC times in handlers
- **Consistent Format**: All datetime fields output in +07:00 timezone format
- **Backward Compatible**: Existing code continues to work
- **Flexible**: Can handle multiple input formats
- **Type Safe**: Compile-time checking for datetime operations
- **GORM Compatible**: Works seamlessly with database operations

## Timezone Behavior Summary

| Input Format | Output Format | Behavior |
|--------------|---------------|----------|
| `"2025-03-26T08:00:00Z"` | `"2025-03-26T08:00:00+07:00"` | Treats as 8:00 AM Bangkok time |
| `"2025-03-26T08:00:00+07:00"` | `"2025-03-26T08:00:00+07:00"` | Preserved as-is |
| `"2025-03-26T08:00:00.000Z"` | `"2025-03-26T08:00:00+07:00"` | Treats as 8:00 AM Bangkok time | 