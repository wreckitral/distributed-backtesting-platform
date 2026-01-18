package domain

import "testing"

func TestBacktestStatusString(t *testing.T) {
	tests := []struct {
		status   BacktestStatus
		expected string
	}{
		{BacktestStatusPending, "PENDING"},
		{BacktestStatusQueued, "QUEUED"},
		{BacktestStatusRunning, "RUNNING"},
		{BacktestStatusCompleted, "COMPLETED"},
		{BacktestStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.expected {
			t.Errorf("Status.String() = %v, want %v", got, tt.expected)
		}
	}
}
