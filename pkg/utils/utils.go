package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func GetPasswordFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("パスワードの読み込みに失敗しました: %w", err)
	}
	
	password = strings.TrimSpace(password)
	
	return password, nil
}

func GetPasswordInteractive(prompt string) (string, error) {
	if prompt == "" {
		prompt = "パスワードを入力してください: "
	}
	
	fmt.Print(prompt)
	
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // 改行を追加
	
	if err != nil {
		return "", fmt.Errorf("パスワードの読み込みに失敗しました: %w", err)
	}
	
	password := strings.TrimSpace(string(passwordBytes))
	
	return password, nil
}

func ExecuteScript(script string) error {
	tmpFile, err := os.CreateTemp("", "envault-*.sh")
	if err != nil {
		return fmt.Errorf("一時ファイルの作成に失敗しました: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := io.WriteString(tmpFile, script); err != nil {
		return fmt.Errorf("スクリプトの書き込みに失敗しました: %w", err)
	}
	
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("一時ファイルのクローズに失敗しました: %w", err)
	}
	
	if err := os.Chmod(tmpFile.Name(), 0700); err != nil {
		return fmt.Errorf("実行権限の設定に失敗しました: %w", err)
	}
	
	fmt.Printf("source %s\n", tmpFile.Name())
	
	return nil
}
