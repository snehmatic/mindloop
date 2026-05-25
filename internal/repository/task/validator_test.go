package task

import "testing"

func TestValidateTaskTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"empty title", "", ErrTaskTitleEmpty},
		{"valid title", "Valid Task", nil},
		{"valid title at max length", string(make([]byte, 100)), nil},
		{"title too long", string(make([]byte, 101)), ErrTaskTitleTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskTitle(tt.title)
			if err != tt.wantErr {
				t.Errorf("ValidateTaskTitle(%q) = %v, want %v", tt.title, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSubTaskTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"empty title", "", ErrSubTaskTitleEmpty},
		{"valid title", "Valid SubTask", nil},
		{"valid title at max length", string(make([]byte, 100)), nil},
		{"title too long", string(make([]byte, 101)), ErrSubTaskTitleTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSubTaskTitle(tt.title)
			if err != tt.wantErr {
				t.Errorf("ValidateSubTaskTitle(%q) = %v, want %v", tt.title, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTaskID(t *testing.T) {
	tests := []struct {
		name    string
		taskID  uint
		wantErr error
	}{
		{"zero taskID", 0, ErrSubTaskInvalidTaskID},
		{"valid taskID", 1, nil},
		{"valid taskID large", 999999, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskID(tt.taskID)
			if err != tt.wantErr {
				t.Errorf("ValidateTaskID(%d) = %v, want %v", tt.taskID, err, tt.wantErr)
			}
		})
	}
}
