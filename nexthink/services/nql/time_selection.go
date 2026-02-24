package nql

import (
	"fmt"
	"time"
)

// Time selection helpers for NQL queries
// Provides a fluent API for building NQL time selection clauses

// =============================================================================
// Time Units
// =============================================================================

// TimeUnit represents a time unit for NQL queries
type TimeUnit string

const (
	Minutes TimeUnit = "min"
	Hours   TimeUnit = "h"
	Days    TimeUnit = "d"
)

// String returns the string representation of the time unit
func (tu TimeUnit) String() string {
	return string(tu)
}

// =============================================================================
// Time Granularity (for summarize by)
// =============================================================================

// TimeGranularity represents time bucket granularity for aggregations
type TimeGranularity string

const (
	// Minute granularities (must be multiples of 15)
	Granularity15Min = "15 min"
	Granularity30Min = "30 min"
	Granularity45Min = "45 min"
	
	// Hour granularities (whole numbers)
	Granularity1Hour  = "1 h"
	Granularity2Hours = "2 h"
	Granularity3Hours = "3 h"
	Granularity6Hours = "6 h"
	Granularity12Hours = "12 h"
	
	// Day granularities (whole numbers)
	Granularity1Day  = "1 d"
	Granularity7Days = "7 d"
	Granularity30Days = "30 d"
)

// String returns the string representation of the time granularity
func (tg TimeGranularity) String() string {
	return string(tg)
}

// =============================================================================
// TimeSelection Builder
// =============================================================================

// TimeSelection represents a time selection clause in NQL
type TimeSelection struct {
	clause string
}

// NewTimeSelection creates a new time selection builder
func NewTimeSelection() *TimeSelection {
	return &TimeSelection{}
}

// DuringPast creates a "during past" time selection
// Example: DuringPast(7, Days) -> "during past 7d"
func (ts *TimeSelection) DuringPast(value int, unit TimeUnit) *TimeSelection {
	ts.clause = fmt.Sprintf("during past %d%s", value, unit)
	return ts
}

// From creates a "from ... to ..." time selection with absolute dates
// Example: From("2024-01-01").To("2024-01-31") -> "from 2024-01-01 to 2024-01-31"
func (ts *TimeSelection) From(date string) *TimeSelectionFrom {
	return &TimeSelectionFrom{
		ts:       ts,
		fromDate: date,
	}
}

// FromRelative creates a "from ... ago to ... ago" time selection
// Example: FromRelative(21, DaysAgo).ToRelative(13, DaysAgo) -> "from 21d ago to 13d ago"
func (ts *TimeSelection) FromRelative(value int, unit TimeUnit) *TimeSelectionFromRelative {
	return &TimeSelectionFromRelative{
		ts:        ts,
		fromValue: value,
		fromUnit:  unit,
	}
}

// On creates an "on" time selection for a specific date
// Example: On("Feb 8, 2024") -> "on Feb 8, 2024"
func (ts *TimeSelection) On(date string) *TimeSelection {
	ts.clause = fmt.Sprintf("on %s", date)
	return ts
}

// OnTime creates an "on" time selection with time.Time
// Example: OnTime(time.Date(2024, 2, 8, 0, 0, 0, 0, time.UTC)) -> "on 2024-02-08"
func (ts *TimeSelection) OnTime(t time.Time) *TimeSelection {
	ts.clause = fmt.Sprintf("on %s", t.Format("2006-01-02"))
	return ts
}

// ByHighResolution adds high-resolution qualifier for VDI data (30s resolution)
// Only works with VDI event data for the past 2 days
// Example: DuringPast(1, Days).ByHighResolution() -> "during past 1d by 30s"
func (ts *TimeSelection) ByHighResolution() *TimeSelection {
	ts.clause = fmt.Sprintf("%s by 30s", ts.clause)
	return ts
}

// String returns the time selection clause
func (ts *TimeSelection) String() string {
	return ts.clause
}

// Build returns the final time selection string
func (ts *TimeSelection) Build() string {
	return ts.clause
}

// =============================================================================
// TimeSelectionFrom (for absolute date ranges)
// =============================================================================

// TimeSelectionFrom represents a time selection starting with "from"
type TimeSelectionFrom struct {
	ts       *TimeSelection
	fromDate string
}

// To completes the "from ... to ..." clause
func (tsf *TimeSelectionFrom) To(date string) *TimeSelection {
	tsf.ts.clause = fmt.Sprintf("from %s to %s", tsf.fromDate, date)
	return tsf.ts
}

// ToTime completes the "from ... to ..." clause with time.Time
func (tsf *TimeSelectionFrom) ToTime(t time.Time) *TimeSelection {
	toDate := t.Format("2006-01-02 15:04:05")
	tsf.ts.clause = fmt.Sprintf("from %s to %s", tsf.fromDate, toDate)
	return tsf.ts
}

// =============================================================================
// TimeSelectionFromRelative (for relative date ranges)
// =============================================================================

// TimeSelectionFromRelative represents a relative time selection starting with "from ... ago"
type TimeSelectionFromRelative struct {
	ts        *TimeSelection
	fromValue int
	fromUnit  TimeUnit
}

// ToRelative completes the "from ... ago to ... ago" clause
func (tsfr *TimeSelectionFromRelative) ToRelative(value int, unit TimeUnit) *TimeSelection {
	tsfr.ts.clause = fmt.Sprintf("from %d%s ago to %d%s ago", tsfr.fromValue, tsfr.fromUnit, value, unit)
	return tsfr.ts
}

// =============================================================================
// Helper Functions
// =============================================================================

// DuringPastMinutes creates a "during past Xmin" time selection
func DuringPastMinutes(minutes int) string {
	return fmt.Sprintf("during past %dmin", minutes)
}

// DuringPastHours creates a "during past Xh" time selection
func DuringPastHours(hours int) string {
	return fmt.Sprintf("during past %dh", hours)
}

// DuringPastDays creates a "during past Xd" time selection
func DuringPastDays(days int) string {
	return fmt.Sprintf("during past %dd", days)
}

// OnDate creates an "on" time selection for a specific date
func OnDate(date string) string {
	return fmt.Sprintf("on %s", date)
}

// FromTo creates a "from ... to ..." time selection
func FromTo(from, to string) string {
	return fmt.Sprintf("from %s to %s", from, to)
}

// FromToRelative creates a "from X ago to Y ago" time selection
func FromToRelative(fromValue int, fromUnit TimeUnit, toValue int, toUnit TimeUnit) string {
	return fmt.Sprintf("from %d%s ago to %d%s ago", fromValue, fromUnit, toValue, toUnit)
}

// =============================================================================
// Predefined Common Time Selections
// =============================================================================

// Common time selections for convenience
var (
	// Past periods
	Past15Minutes = "during past 15min"
	Past30Minutes = "during past 30min"
	Past1Hour     = "during past 1h"
	Past24Hours   = "during past 24h"
	Past7Days     = "during past 7d"
	Past30Days    = "during past 30d"
	
	// Yesterday
	Yesterday = "from 1d ago to 1d ago"
)

// =============================================================================
// Time Format Helpers
// =============================================================================

// FormatDateTime formats a time.Time for NQL date-time expressions
// Example: 2024-02-08 15:30:00
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate formats a time.Time for NQL date expressions
// Example: 2024-02-08
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateShort formats a time.Time for NQL short date expressions
// Example: Feb 8, 2024
func FormatDateShort(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

// ParseRelativeTime converts a Go duration to NQL relative time format
// Example: 7 * 24 * time.Hour -> "7d ago"
func ParseRelativeTime(d time.Duration) string {
	if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	} else if d.Hours() >= 1 {
		hours := int(d.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		minutes := int(d.Minutes())
		return fmt.Sprintf("%dmin ago", minutes)
	}
}

// =============================================================================
// Validation
// =============================================================================

// ValidateTimeSelection validates a time selection string
func ValidateTimeSelection(selection string) error {
	if selection == "" {
		return fmt.Errorf("time selection cannot be empty")
	}
	
	// Basic validation - could be expanded
	// For now, just check that it's not empty
	return nil
}
