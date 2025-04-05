package executor

import (
	"runtime"
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	// Skip certain tests on Windows
	isWindows := runtime.GOOS == "windows"

	tests := []struct {
		name        string
		command     string
		expectError bool
		contains    string
		skipWindows bool
	}{
		{
			name:        "Echo command",
			command:     "echo Hello",
			expectError: false,
			contains:    "Hello",
		},
		{
			name:        "Empty command",
			command:     "",
			expectError: true,
		},
		{
			name:        "Invalid command",
			command:     "nonexistentcommand",
			expectError: true,
		},
		{
			name:        "List directory",
			command:     "ls",
			expectError: false,
			skipWindows: true,
		},
		{
			name:        "List empty pattern",
			command:     "ls nonexistentpattern*",
			expectError: false,
			contains:    "No matching files found",
			skipWindows: true,
		},
	}

	for _, tt := range tests {
		if tt.skipWindows && isWindows {
			t.Logf("Skipping %q on Windows", tt.name)
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			output, err := ExecuteCommand(tt.command)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.contains != "" && !strings.Contains(output, tt.contains) {
					t.Errorf("expected output to contain %q, got: %q", tt.contains, output)
				}
			}
		})
	}
}
