// ABOUTME: Characterizes the public command contract for the Go-only Folio application.
// ABOUTME: Exercises dispatch, help, version, diagnostics, and stable process status.
package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunTopLevelContract(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStatus int
		wantOut    []string
		wantErr    []string
	}{
		{
			name:       "help",
			args:       []string{"--help"},
			wantStatus: 0,
			wantOut:    []string{"format converter", "convert", "letter", "manuscript"},
		},
		{
			name:       "short help",
			args:       []string{"-h"},
			wantStatus: 0,
			wantOut:    []string{"format converter", "convert", "letter", "manuscript"},
		},
		{
			name:       "version",
			args:       []string{"--version"},
			wantStatus: 0,
			wantOut:    []string{"folio 0.4.10"},
		},
		{
			name:       "missing command",
			wantStatus: 1,
			wantErr:    []string{"Error: no subcommand given", "Usage: folio <command>"},
		},
		{
			name:       "unknown command",
			args:       []string{"unknown"},
			wantStatus: 1,
			wantErr:    []string{"Error: unknown subcommand 'unknown'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			status := Run(tt.args, strings.NewReader(""), &stdout, &stderr)
			if status != tt.wantStatus {
				t.Fatalf("status = %d, want %d\nstdout:\n%s\nstderr:\n%s", status, tt.wantStatus, stdout.String(), stderr.String())
			}
			assertTextContains(t, stdout.String(), tt.wantOut)
			assertTextContains(t, stderr.String(), tt.wantErr)
		})
	}
}

func assertTextContains(t *testing.T, text string, fragments []string) {
	t.Helper()
	for _, fragment := range fragments {
		if !strings.Contains(text, fragment) {
			t.Errorf("output missing %q:\n%s", fragment, text)
		}
	}
}
