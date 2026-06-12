package handlers

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashing(t *testing.T) {
	password := "my_secure_music_password_123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate password hash: %v", err)
	}

	if len(hashed) == 0 {
		t.Fatal("Hashed password is empty")
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Fatalf("Failed to verify password with correct credentials: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrong_password"))
	if err == nil {
		t.Fatal("Expected error with incorrect password but got nil")
	}
}
