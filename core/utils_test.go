package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenRandomToken(t *testing.T) {
	token := GenRandomToken()
	if len(token) != 64 {
		t.Errorf("expected token length 64, got %d", len(token))
	}

	token2 := GenRandomToken()
	if token == token2 {
		t.Error("two generated tokens should not be equal")
	}
}

func TestGenRandomString(t *testing.T) {
	s := GenRandomString(16)
	if len(s) != 16 {
		t.Errorf("expected string length 16, got %d", len(s))
	}

	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			t.Errorf("unexpected character '%c' in random string", c)
		}
	}
}

func TestGenRandomStringZeroLength(t *testing.T) {
	s := GenRandomString(0)
	if len(s) != 0 {
		t.Errorf("expected empty string, got '%s'", s)
	}
}

func TestGenRandomAlphanumString(t *testing.T) {
	s := GenRandomAlphanumString(20)
	if len(s) != 20 {
		t.Errorf("expected string length 20, got %d", len(s))
	}

	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Errorf("unexpected character '%c' in random alphanum string", c)
		}
	}
}

func TestParseDurationString(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"1d", 24 * time.Hour, false},
		{"2h", 2 * time.Hour, false},
		{"30m", 30 * time.Minute, false},
		{"45s", 45 * time.Second, false},
		{"1d2h3m4s", 24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second, false},
		{"1h30m", time.Hour + 30*time.Minute, false},
		{"0s", 0, false},
		{"1x", 0, true},
		{"abc", 0, true},
		{"1s1d", 0, true},
	}

	for _, tt := range tests {
		dur, err := ParseDurationString(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseDurationString(%q): expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseDurationString(%q): unexpected error: %v", tt.input, err)
			} else if dur != tt.expected {
				t.Errorf("ParseDurationString(%q): expected %v, got %v", tt.input, tt.expected, dur)
			}
		}
	}
}

func TestGetDurationString(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		expire   time.Time
		expected string
	}{
		{now.Add(24*time.Hour + 2*time.Hour + 3*time.Minute + 4*time.Second), "1d2h3m4s"},
		{now.Add(2*time.Hour + 30*time.Minute), "2h30m0s"},
		{now.Add(45 * time.Second), "45s"},
		{now.Add(-1 * time.Hour), ""},
		{now, ""},
	}

	for _, tt := range tests {
		result := GetDurationString(now, tt.expire)
		if result != tt.expected {
			t.Errorf("GetDurationString(now, %v): expected %q, got %q", tt.expire, tt.expected, result)
		}
	}
}

func TestCreateDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "testdir")
	err := CreateDir(dir, 0700)
	if err != nil {
		t.Fatalf("CreateDir: unexpected error: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory")
	}

	err = CreateDir(dir, 0700)
	if err != nil {
		t.Errorf("CreateDir on existing dir should not fail: %v", err)
	}
}

func TestSaveToFileAndReadFromFile(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "test.txt")
	content := []byte("hello world")

	err := SaveToFile(content, fpath, 0644)
	if err != nil {
		t.Fatalf("SaveToFile: unexpected error: %v", err)
	}

	data, err := ReadFromFile(fpath)
	if err != nil {
		t.Fatalf("ReadFromFile: unexpected error: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(data))
	}
}

func TestReadFromFileNotFound(t *testing.T) {
	_, err := ReadFromFile("/nonexistent/file/path")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

