package cli

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestNewCLI(t *testing.T) {
	cli := NewCLI()
	if cli == nil {
		t.Errorf("NewCLI() = nil, want non-nil")
	}
	
	// rootCmdがnilでないことを確認
	if cli.rootCmd == nil {
		t.Errorf("cli.rootCmd = nil, want non-nil")
	}
}

func captureOutput(f func() error) (string, error) {
	// 標準出力をキャプチャするためのバッファを作成
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	
	// 元の標準出力を保存
	stdout := os.Stdout
	// 標準出力をパイプに切り替え
	os.Stdout = w
	
	// 関数を実行
	err = f()
	
	// パイプを閉じて標準出力を元に戻す
	w.Close()
	os.Stdout = stdout
	
	// パイプからすべての出力を読み取る
	var buf bytes.Buffer
	io.Copy(&buf, r)
	
	return buf.String(), err
}

func TestRunHelp(t *testing.T) {
	cli := NewCLI()
	
	// helpコマンドのテスト
	output, err := captureOutput(func() error {
		return cli.Run([]string{"help"})
	})
	if err != nil {
		t.Errorf("Run(help) error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run(help) output is empty")
	}
	
	// -hフラグのテスト
	output, err = captureOutput(func() error {
		return cli.Run([]string{"-h"})
	})
	if err != nil {
		t.Errorf("Run(-h) error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run(-h) output is empty")
	}
	
	// --helpフラグのテスト
	output, err = captureOutput(func() error {
		return cli.Run([]string{"--help"})
	})
	if err != nil {
		t.Errorf("Run(--help) error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run(--help) output is empty")
	}
}

func TestRunVersion(t *testing.T) {
	cli := NewCLI()
	
	// versionコマンドのテスト
	output, err := captureOutput(func() error {
		return cli.Run([]string{"version"})
	})
	if err != nil {
		t.Errorf("Run(version) error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run(version) output is empty")
	}
	
	// --versionフラグのテスト
	output, err = captureOutput(func() error {
		return cli.Run([]string{"--version"})
	})
	if err != nil {
		t.Errorf("Run(--version) error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run(--version) output is empty")
	}
}

func TestRunWithNoArgs(t *testing.T) {
	cli := NewCLI()
	
	// 引数なしの実行テスト
	output, err := captureOutput(func() error {
		return cli.Run([]string{})
	})
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if len(output) == 0 {
		t.Errorf("Run() output is empty")
	}
}

func TestExportCommand(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
	// exportコマンドの基本的な動作確認
	cli := NewCLI()
	err := cli.Run([]string{"export", "-h"})
	if err != nil {
		t.Errorf("Run(export -h) error = %v", err)
	}
}

func TestUnsetCommand(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
	// unsetコマンドの基本的な動作確認
	cli := NewCLI()
	err := cli.Run([]string{"unset", "-h"})
	if err != nil {
		t.Errorf("Run(unset -h) error = %v", err)
	}
}

func TestDumpCommand(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
	// dumpコマンドの基本的な動作確認
	cli := NewCLI()
	err := cli.Run([]string{"dump", "-h"})
	if err != nil {
		t.Errorf("Run(dump -h) error = %v", err)
	}
}

func TestRunEncrypt(t *testing.T) {
	t.Skip("このテストは対話的な入力が必要なため、スキップします")
}