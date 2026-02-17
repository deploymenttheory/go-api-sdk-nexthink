package nql

// isTerminalStatus checks if the export status is in a terminal state
func isTerminalStatus(status string) bool {
	switch status {
	case ExportStatusCompleted, ExportStatusError:
		return true
	default:
		return false
	}
}
