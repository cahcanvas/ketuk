package pay

import (
	"crypto/hmac"
	"testing"
)

func TestVerifyCallbackHMAC(t *testing.T) {
	d := NewDuitku("DXXXX", "secret", "")
	sig := hmacSHA256("DXXXX"+"150000"+"abcde", "secret")
	if !d.VerifyCallback("DXXXX", "150000", "abcde", sig) {
		t.Fatal("hmac should verify")
	}
	if d.VerifyCallback("DXXXX", "150000", "abcde", "nope") {
		t.Fatal("bad signature")
	}
}

func TestHmacStable(t *testing.T) {
	a := hmacSHA256("DXXXX150000abcde", "secret")
	b := hmacSHA256("DXXXX150000abcde", "secret")
	if !hmac.Equal([]byte(a), []byte(b)) || a == "" {
		t.Fatal(a, b)
	}
}
