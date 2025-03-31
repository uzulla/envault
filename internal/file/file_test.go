package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadEnvFile(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.env")
	testContent := []byte("TEST_VAR1=value1\nTEST_VAR2=value2")

	err := os.WriteFile(testFilePath, testContent, 0600)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	data, err := ReadEnvFile(testFilePath)
	if err != nil {
		t.Errorf("ReadEnvFile() error = %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("ReadEnvFile() = %v, want %v", string(data), string(testContent))
	}

	_, err = ReadEnvFile(filepath.Join(tempDir, "nonexistent.env"))
	if err != ErrFileNotFound {
		t.Errorf("存在しないファイルの読み込みで期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrFileNotFound, err)
	}

	emptyFilePath := filepath.Join(tempDir, "empty.env")
	err = os.WriteFile(emptyFilePath, []byte{}, 0600)
	if err != nil {
		t.Fatalf("空ファイルの作成に失敗しました: %v", err)
	}

	_, err = ReadEnvFile(emptyFilePath)
	if err != ErrEmptyFile {
		t.Errorf("空ファイルの読み込みで期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrEmptyFile, err)
	}
}

func TestWriteVaultedFile(t *testing.T) {
	tempDir := t.TempDir()
	testContent := []byte("暗号化されたテストデータ")
	testFilePath := filepath.Join(tempDir, DefaultVaultedFileName)

	err := WriteVaultedFile(testContent, testFilePath)
	if err != nil {
		t.Errorf("WriteVaultedFile() error = %v", err)
	}

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("ファイルが作成されませんでした: %v", err)
	}

	data, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("ファイルの内容が期待と異なります。期待: %v, 実際: %v", string(testContent), string(data))
	}

}

func TestReadVaultedFile(t *testing.T) {
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, DefaultVaultedFileName)
	testContent := []byte("暗号化されたテストデータ")

	err := os.WriteFile(testFilePath, testContent, 0600)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	data, err := ReadVaultedFile(testFilePath)
	if err != nil {
		t.Errorf("ReadVaultedFile() error = %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("ReadVaultedFile() = %v, want %v", string(data), string(testContent))
	}

	_, err = ReadVaultedFile(filepath.Join(tempDir, "nonexistent.env.vaulted"))
	if err != ErrFileNotFound {
		t.Errorf("存在しないファイルの読み込みで期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrFileNotFound, err)
	}

	emptyFilePath := filepath.Join(tempDir, "empty.env.vaulted")
	err = os.WriteFile(emptyFilePath, []byte{}, 0600)
	if err != nil {
		t.Fatalf("空ファイルの作成に失敗しました: %v", err)
	}

	_, err = ReadVaultedFile(emptyFilePath)
	if err != ErrEmptyFile {
		t.Errorf("空ファイルの読み込みで期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrEmptyFile, err)
	}

}
