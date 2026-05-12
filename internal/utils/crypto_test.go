package utils

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	original := "my-secret-api-token"
	encrypted, err := Encrypt(original)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	if encrypted == original {
		t.Fatalf("Encrypted text is same as original")
	}
	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}
	if decrypted != original {
		t.Fatalf("Expected %s, got %s", original, decrypted)
	}
}

func TestEncryptDecryptEmpty(t *testing.T) {
	encrypted, err := Encrypt("")
	if err != nil {
		t.Fatalf("Failed to encrypt empty string: %v", err)
	}
	if encrypted != "" {
		t.Fatalf("Expected empty string, got %s", encrypted)
	}
	decrypted, err := Decrypt("")
	if err != nil {
		t.Fatalf("Failed to decrypt empty string: %v", err)
	}
	if decrypted != "" {
		t.Fatalf("Expected empty string, got %s", decrypted)
	}
}
