package customtype

import "time"

// DBOpsSummary represents the summary of database operations
type DBOpsSummary struct {
	Status       string    `json:"status"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     string    `json:"duration"`
	NewRecs      int       `json:"new_records"`
	ModifiedRecs int       `json:"modified_records"`
	TotalOps     int       `json:"total_operations"`
	TotalRecs    int       `json:"total_records"`
}
