package nql

import (
	"testing"
	"time"
)

func TestTimeSelection_DuringPast(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		unit     TimeUnit
		expected string
	}{
		{"7 days", 7, Days, "during past 7d"},
		{"24 hours", 24, Hours, "during past 24h"},
		{"30 days", 30, Days, "during past 30d"},
		{"60 minutes", 60, Minutes, "during past 60min"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTimeSelection().
				DuringPast(tt.value, tt.unit).
				Build()

			if result != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestTimeSelection_FromTo(t *testing.T) {
	result := NewTimeSelection().
		From("2024-01-01").
		To("2024-01-31").
		Build()

	expected := "from 2024-01-01 to 2024-01-31"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestTimeSelection_On(t *testing.T) {
	result := NewTimeSelection().
		On("Feb 8, 2024").
		Build()

	expected := "on Feb 8, 2024"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestTimeSelection_FromRelativeToRelative(t *testing.T) {
	result := NewTimeSelection().
		FromRelative(21, Days).
		ToRelative(13, Days).
		Build()

	expected := "from 21d ago to 13d ago"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestTimeSelection_ByHighResolution(t *testing.T) {
	result := NewTimeSelection().
		DuringPast(2, Hours).
		ByHighResolution().
		Build()

	expected := "during past 2h by 30s"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestPredefinedTimeSelections(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Past7Days", Past7Days, "during past 7d"},
		{"Past24Hours", Past24Hours, "during past 24h"},
		{"Past30Days", Past30Days, "during past 30d"},
		{"Past1Hour", Past1Hour, "during past 1h"},
		{"Yesterday", Yesterday, "from 1d ago to 1d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestTimeUnit_String(t *testing.T) {
	tests := []struct {
		unit     TimeUnit
		expected string
	}{
		{Minutes, "min"},
		{Hours, "h"},
		{Days, "d"},
	}

	for _, tt := range tests {
		t.Run(string(tt.unit), func(t *testing.T) {
			if string(tt.unit) != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, string(tt.unit))
			}
		})
	}
}

func TestTimeGranularity_Values(t *testing.T) {
	tests := []struct {
		name        string
		granularity TimeGranularity
		expected    string
	}{
		{"15 minutes", Granularity15Min, "15 min"},
		{"1 hour", Granularity1Hour, "1 h"},
		{"1 day", Granularity1Day, "1 d"},
		{"7 days", Granularity7Days, "7 d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.granularity) != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, string(tt.granularity))
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	dt := time.Date(2024, 2, 8, 10, 15, 30, 0, time.UTC)

	result := FormatDateTime(dt)
	expected := "2024-02-08 10:15:30"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestFormatDate(t *testing.T) {
	dt := time.Date(2024, 2, 8, 10, 15, 30, 0, time.UTC)

	result := FormatDate(dt)
	expected := "2024-02-08"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestFormatDateShort(t *testing.T) {
	dt := time.Date(2024, 2, 8, 10, 15, 30, 0, time.UTC)

	result := FormatDateShort(dt)
	expected := "Feb 8, 2024"

	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}

func TestTimeSelection_Empty(t *testing.T) {
	result := NewTimeSelection().Build()

	if result != "" {
		t.Errorf("Expected empty string for empty time selection, got: %s", result)
	}
}
