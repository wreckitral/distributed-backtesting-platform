package domain

import "testing"

func TestBacktestStatusString(t *testing.T) {
	tests := []struct {
		status   BacktestStatus
		expected string
	}{
		{StatusPending, "PENDING"},
		{StatusQueued, "QUEUED"},
		{StatusRunning, "RUNNING"},
		{StatusCompleted, "COMPLETED"},
		{StatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.expected {
			t.Errorf("Status.String() = %v, want %v", got, tt.expected)
		}
	}
}
