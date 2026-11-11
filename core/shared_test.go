package core

import (
	"testing"
)

func TestCombineHost(t *testing.T) {
	tests := []struct {
		sub, domain, expected string
	}{
		{"www", "example.com", "www.example.com"},
		{"", "example.com", "example.com"},
		{"mail", "google.com", "mail.google.com"},
		{"sub.deep", "example.com", "sub.deep.example.com"},
	}

	for _, tt := range tests {
		result := combineHost(tt.sub, tt.domain)
		if result != tt.expected {
			t.Errorf("combineHost(%q, %q): expected %q, got %q", tt.sub, tt.domain, tt.expected, result)
		}
	}
}

func TestObfuscateDots(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"example.com", "example[[d0t]]com"},
		{"no dots", "no dots"},
		{"a.b.c", "a[[d0t]]b[[d0t]]c"},
		{"", ""},
	}

	for _, tt := range tests {
		result := obfuscateDots(tt.input)
		if result != tt.expected {
			t.Errorf("obfuscateDots(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestRemoveObfuscatedDots(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"example[[d0t]]com", "example.com"},
		{"no dots", "no dots"},
		{"a[[d0t]]b[[d0t]]c", "a.b.c"},
	}

	for _, tt := range tests {
		result := removeObfuscatedDots(tt.input)
		if result != tt.expected {
			t.Errorf("removeObfuscatedDots(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestObfuscateRoundtrip(t *testing.T) {
	input := "www.example.com"
	result := removeObfuscatedDots(obfuscateDots(input))
	if result != input {
		t.Errorf("roundtrip failed: expected %q, got %q", input, result)
	}
}

func TestStringExists(t *testing.T) {
	sa := []string{"apple", "banana", "cherry"}
	if !stringExists("banana", sa) {
		t.Error("expected 'banana' to exist")
	}
	if stringExists("grape", sa) {
		t.Error("expected 'grape' to not exist")
	}
	if stringExists("apple", []string{}) {
		t.Error("expected nothing to exist in empty slice")
	}
}

func TestIntExists(t *testing.T) {
	ia := []int{1, 2, 3, 5, 8}
	if !intExists(3, ia) {
		t.Error("expected 3 to exist")
	}
	if intExists(4, ia) {
		t.Error("expected 4 to not exist")
	}
	if intExists(1, []int{}) {
		t.Error("expected nothing to exist in empty slice")
	}
}

func TestRemoveString(t *testing.T) {
	sa := []string{"a", "b", "c", "d"}
	result := removeString("b", sa)
	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
	if stringExists("b", result) {
		t.Error("'b' should have been removed")
	}

	result2 := removeString("z", []string{"a", "b"})
	if len(result2) != 2 {
		t.Error("removing nonexistent element should not change slice")
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this...ing"},
		{"ab", 10, "ab"},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%q, %d): expected %q, got %q", tt.input, tt.maxLen, tt.expected, result)
		}
	}
}
