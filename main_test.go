package main

import (
	"testing"
)

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		command  string
		expected int
	}{
		{"", CMD_HELP},
		{"--help", CMD_HELP},
		{"-h", CMD_HELP},
		{"--version", CMD_VERSION},
		{"-v", CMD_VERSION},
		{"--update-key", CMD_UPDATE},
		{"-u", CMD_UPDATE},
		{"--delete-key", CMD_DELETE},
		{"-d", CMD_DELETE},
		{"random strings", CMD_COMPLETION},
		{"--test", CMD_COMPLETION},
		{"-test", CMD_COMPLETION},
		{"how to extract test.tar.gz", CMD_COMPLETION},
	}

	for _, test := range tests {
		result := handleCommand(test.command)
		if result != test.expected {
			t.Error("Expected", test.expected, "got", result)
		}
	}
}
