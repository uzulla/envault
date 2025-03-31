package utils

import (
	"os"
	"strings"
	"testing"
)

func TestGetPasswordFromStdin(t *testing.T) {
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
	}()

	testPassword := "testpassword\n"
	go func() {
		w.Write([]byte(testPassword))
		w.Close()
	}()

	password, err := GetPasswordFromStdin()
	if err != nil {
		t.Errorf("GetPasswordFromStdin() error = %v", err)
	}

	expectedPassword := strings.TrimSpace(testPassword)
	if password != expectedPassword {
		t.Errorf("GetPasswordFromStdin() = %v, want %v", password, expectedPassword)
	}
}

func TestExecuteScript(t *testing.T) {
	testScript := "#!/bin/bash\necho 'Test script executed'"

	err := ExecuteScript(testScript)
	if err != nil {
		t.Errorf("ExecuteScript() error = %v", err)
	}

}

func TestGetPasswordInteractive(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
}
