package httpserver

import (
	"testing"
	"time"
)

func TestRateLimiterBlocksAfterLimit(t *testing.T) {
	l := newRateLimiter(2, time.Minute)
	if !l.allow("1.2.3.4") || !l.allow("1.2.3.4") {
		t.Fatal("first two attempts should pass")
	}
	if l.allow("1.2.3.4") {
		t.Fatal("third attempt should be blocked")
	}
	if !l.allow("5.6.7.8") {
		t.Fatal("a different caller must have its own window")
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	l := newRateLimiter(1, time.Millisecond)
	if !l.allow("1.2.3.4") {
		t.Fatal("first attempt should pass")
	}
	time.Sleep(2 * time.Millisecond)
	if !l.allow("1.2.3.4") {
		t.Fatal("window should have reset")
	}
}
