package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultVaultedFileName = ".env.vaulted"
)

var (
	ErrFileNotFound = errors.New("ファイルが見つかりません")
	ErrEmptyFile    = errors.New("ファイルが空です")
)

func ReadEnvFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	if len(data) == 0 {
		return nil, ErrEmptyFile
	}

	return data, nil
}

func WriteVaultedFile(data []byte, outputPath string) error {
	if outputPath == "" {
		outputPath = DefaultVaultedFileName
	}

	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("ファイル '%s' は既に存在します。上書きしますか？ [y/N]: ", outputPath)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("入力の読み取りに失敗しました: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			return fmt.Errorf("操作がキャンセルされました")
		}
	}

	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

func ReadVaultedFile(filePath string) ([]byte, error) {
	if filePath == "" {
		filePath = DefaultVaultedFileName
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	if len(data) == 0 {
		return nil, ErrEmptyFile
	}

	return data, nil
}
