package cli

import (
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli := NewCLI()
	if cli == nil {
		t.Errorf("NewCLI() = nil, want non-nil")
	}
}

func TestRunHelp(t *testing.T) {
	cli := NewCLI()
	
	err := cli.Run([]string{"help"})
	if err != nil {
		t.Errorf("Run(help) error = %v", err)
	}
	
	err = cli.Run([]string{"-h"})
	if err != nil {
		t.Errorf("Run(-h) error = %v", err)
	}
	
	err = cli.Run([]string{"--help"})
	if err != nil {
		t.Errorf("Run(--help) error = %v", err)
	}
}

func TestRunVersion(t *testing.T) {
	cli := NewCLI()
	
	err := cli.Run([]string{"version"})
	if err != nil {
		t.Errorf("Run(version) error = %v", err)
	}
	
	err = cli.Run([]string{"-v"})
	if err != nil {
		t.Errorf("Run(-v) error = %v", err)
	}
	
	err = cli.Run([]string{"--version"})
	if err != nil {
		t.Errorf("Run(--version) error = %v", err)
	}
}

func TestRunWithNoArgs(t *testing.T) {
	cli := NewCLI()
	
	err := cli.Run([]string{})
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
}


func TestRunEncrypt(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
}

func TestRunExport(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
}

func TestRunUnset(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
}
