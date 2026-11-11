package parser

import (
	"testing"
)

func TestParseSimpleCommand(t *testing.T) {
	args, err := Parse("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0] != "hello" || args[1] != "world" {
		t.Errorf("expected [hello, world], got %v", args)
	}
}

func TestParseEmptyString(t *testing.T) {
	args, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestParseDoubleQuotedString(t *testing.T) {
	args, err := Parse(`set key "hello world"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[2] != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", args[2])
	}
}

func TestParseSingleQuotedString(t *testing.T) {
	args, err := Parse(`set key 'hello world'`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[2] != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", args[2])
	}
}

func TestParseEscapedCharacters(t *testing.T) {
	args, err := Parse(`hello\ world`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	if args[0] != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", args[0])
	}
}

func TestParseUnterminatedQuote(t *testing.T) {
	_, err := Parse(`"hello world`)
	if err == nil {
		t.Error("expected error for unterminated double quote")
	}
}

func TestParseUnterminatedSingleQuote(t *testing.T) {
	_, err := Parse(`'hello world`)
	if err == nil {
		t.Error("expected error for unterminated single quote")
	}
}

func TestParseTrailingBackslash(t *testing.T) {
	_, err := Parse(`hello\`)
	if err == nil {
		t.Error("expected error for trailing backslash")
	}
}

func TestParseMultipleSpaces(t *testing.T) {
	args, err := Parse("hello    world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0] != "hello" || args[1] != "world" {
		t.Errorf("expected [hello, world], got %v", args)
	}
}

func TestParseEmptyQuotedString(t *testing.T) {
	args, err := Parse(`set key ""`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
	if args[2] != "" {
		t.Errorf("expected empty string, got '%s'", args[2])
	}
}

func TestParseMixedQuotes(t *testing.T) {
	args, err := Parse(`"it's" 'say "hi"'`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0] != "it's" {
		t.Errorf("expected \"it's\", got '%s'", args[0])
	}
	if args[1] != `say "hi"` {
		t.Errorf("expected 'say \"hi\"', got '%s'", args[1])
	}
}

func TestParseTabsAndNewlines(t *testing.T) {
	args, err := Parse("hello\tworld\nfoo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(args))
	}
}

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Fatal("NewParser returned nil")
	}
	args, err := p.Parse("test command")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
}

func TestIsSpace(t *testing.T) {
	spaces := []rune{' ', '\t', '\r', '\n'}
	for _, r := range spaces {
		if !isSpace(r) {
			t.Errorf("expected isSpace(%q) to be true", r)
		}
	}
	nonSpaces := []rune{'a', '1', '-', '.'}
	for _, r := range nonSpaces {
		if isSpace(r) {
			t.Errorf("expected isSpace(%q) to be false", r)
		}
	}
}
