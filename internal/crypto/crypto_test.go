package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	testData := []byte("TEST_VAR1=value1\nTEST_VAR2=value2")
	password := "testpassword"

	encrypted, err := Encrypt(testData, password)
	if err != nil {
		t.Fatalf("暗号化に失敗しました: %v", err)
	}

	if !bytes.HasPrefix(encrypted, []byte(MagicBytes)) {
		t.Errorf("暗号化されたデータにマジックバイトが含まれていません")
	}

	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("復号化に失敗しました: %v", err)
	}

	if !bytes.Equal(testData, decrypted) {
		t.Errorf("復号化されたデータが元のデータと一致しません\n元のデータ: %s\n復号化されたデータ: %s", testData, decrypted)
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	testData := []byte("TEST_VAR1=value1\nTEST_VAR2=value2")
	password := "correctpassword"
	wrongPassword := "wrongpassword"

	encrypted, err := Encrypt(testData, password)
	if err != nil {
		t.Fatalf("暗号化に失敗しました: %v", err)
	}

	_, err = Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Errorf("間違ったパスワードで復号化が成功しました")
	}

	if err != ErrDecryptionFailed {
		t.Errorf("期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrDecryptionFailed, err)
	}
}

func TestInvalidEncryptedData(t *testing.T) {
	invalidData := []byte("INVALID")
	password := "testpassword"

	_, err := Decrypt(invalidData, password)
	if err == nil {
		t.Errorf("無効なデータの復号化が成功しました")
	}

	if err != ErrInvalidFile {
		t.Errorf("期待されるエラーが返されませんでした。期待: %v, 実際: %v", ErrInvalidFile, err)
	}
}

func TestEncryptionStrength(t *testing.T) {
	testData := []byte("TEST_VAR1=value1\nTEST_VAR2=value2")
	password := "testpassword"

	encrypted, err := Encrypt(testData, password)
	if err != nil {
		t.Fatalf("暗号化に失敗しました: %v", err)
	}

	encryptedStr := string(encrypted)
	if strings.Contains(encryptedStr, "TEST_VAR1") || strings.Contains(encryptedStr, "TEST_VAR2") {
		t.Errorf("暗号化されたデータに平文が含まれています")
	}

	if strings.Contains(encryptedStr, "value1") || strings.Contains(encryptedStr, "value2") {
		t.Errorf("暗号化されたデータに平文が含まれています")
	}

	encrypted2, err := Encrypt(testData, password)
	if err != nil {
		t.Fatalf("2回目の暗号化に失敗しました: %v", err)
	}

	if bytes.Equal(encrypted, encrypted2) {
		t.Errorf("2回の暗号化結果が同じです。ソルトまたはnonceがランダムでない可能性があります")
	}
}
