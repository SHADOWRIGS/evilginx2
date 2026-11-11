package database

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestDB(t *testing.T) *Database {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	return db
}

func TestNewDatabase(t *testing.T) {
	db := newTestDB(t)
	if db == nil {
		t.Fatal("NewDatabase returned nil")
	}
}

func TestNewDatabaseInvalidPath(t *testing.T) {
	_, err := NewDatabase("/nonexistent/dir/test.db")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestCreateAndListSessions(t *testing.T) {
	db := newTestDB(t)

	err := db.CreateSession("sid-1", "example", "https://example.com", "Mozilla/5.0", "192.168.1.1")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	sessions, err := db.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].SessionId != "sid-1" {
		t.Errorf("expected session_id 'sid-1', got '%s'", sessions[0].SessionId)
	}
	if sessions[0].Phishlet != "example" {
		t.Errorf("expected phishlet 'example', got '%s'", sessions[0].Phishlet)
	}
}

func TestCreateDuplicateSession(t *testing.T) {
	db := newTestDB(t)

	err := db.CreateSession("sid-dup", "example", "https://example.com", "UA", "1.2.3.4")
	if err != nil {
		t.Fatalf("first CreateSession: %v", err)
	}

	err = db.CreateSession("sid-dup", "example", "https://example.com", "UA", "1.2.3.4")
	if err == nil {
		t.Error("expected error for duplicate session")
	}
}

func TestSetSessionUsername(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-u", "test", "https://test.com", "UA", "1.1.1.1")

	err := db.SetSessionUsername("sid-u", "admin")
	if err != nil {
		t.Fatalf("SetSessionUsername: %v", err)
	}

	sessions, _ := db.ListSessions()
	if sessions[0].Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", sessions[0].Username)
	}
}

func TestSetSessionPassword(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-p", "test", "https://test.com", "UA", "1.1.1.1")

	err := db.SetSessionPassword("sid-p", "secret123")
	if err != nil {
		t.Fatalf("SetSessionPassword: %v", err)
	}

	sessions, _ := db.ListSessions()
	if sessions[0].Password != "secret123" {
		t.Errorf("expected password 'secret123', got '%s'", sessions[0].Password)
	}
}

func TestSetSessionCustom(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-c", "test", "https://test.com", "UA", "1.1.1.1")

	err := db.SetSessionCustom("sid-c", "token", "abc123")
	if err != nil {
		t.Fatalf("SetSessionCustom: %v", err)
	}

	sessions, _ := db.ListSessions()
	if sessions[0].Custom["token"] != "abc123" {
		t.Errorf("expected custom token 'abc123', got '%s'", sessions[0].Custom["token"])
	}
}

func TestSetSessionBodyTokens(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-bt", "test", "https://test.com", "UA", "1.1.1.1")

	tokens := map[string]string{"csrf": "token123"}
	err := db.SetSessionBodyTokens("sid-bt", tokens)
	if err != nil {
		t.Fatalf("SetSessionBodyTokens: %v", err)
	}

	sessions, _ := db.ListSessions()
	if sessions[0].BodyTokens["csrf"] != "token123" {
		t.Errorf("expected body token 'token123', got '%s'", sessions[0].BodyTokens["csrf"])
	}
}

func TestSetSessionHttpTokens(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-ht", "test", "https://test.com", "UA", "1.1.1.1")

	tokens := map[string]string{"auth": "bearer-xyz"}
	err := db.SetSessionHttpTokens("sid-ht", tokens)
	if err != nil {
		t.Fatalf("SetSessionHttpTokens: %v", err)
	}

	sessions, _ := db.ListSessions()
	if sessions[0].HttpTokens["auth"] != "bearer-xyz" {
		t.Errorf("expected http token 'bearer-xyz', got '%s'", sessions[0].HttpTokens["auth"])
	}
}

func TestSetSessionCookieTokens(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-ct", "test", "https://test.com", "UA", "1.1.1.1")

	tokens := map[string]map[string]*CookieToken{
		"example.com": {
			"session": {
				Name:     "session",
				Value:    "abc",
				Path:     "/",
				HttpOnly: true,
			},
		},
	}
	err := db.SetSessionCookieTokens("sid-ct", tokens)
	if err != nil {
		t.Fatalf("SetSessionCookieTokens: %v", err)
	}

	sessions, _ := db.ListSessions()
	ct := sessions[0].CookieTokens["example.com"]["session"]
	if ct == nil || ct.Value != "abc" {
		t.Errorf("expected cookie token value 'abc'")
	}
}

func TestDeleteSession(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-del", "test", "https://test.com", "UA", "1.1.1.1")

	err := db.DeleteSession("sid-del")
	if err != nil {
		t.Fatalf("DeleteSession: %v", err)
	}

	sessions, _ := db.ListSessions()
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions after delete, got %d", len(sessions))
	}
}

func TestDeleteSessionById(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-delid", "test", "https://test.com", "UA", "1.1.1.1")

	sessions, _ := db.ListSessions()
	id := sessions[0].Id

	err := db.DeleteSessionById(id)
	if err != nil {
		t.Fatalf("DeleteSessionById: %v", err)
	}

	sessions, _ = db.ListSessions()
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions after delete, got %d", len(sessions))
	}
}

func TestDeleteNonexistentSession(t *testing.T) {
	db := newTestDB(t)
	err := db.DeleteSession("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestMultipleSessions(t *testing.T) {
	db := newTestDB(t)

	for i := 0; i < 5; i++ {
		sid := "sid-multi-" + string(rune('a'+i))
		err := db.CreateSession(sid, "test", "https://test.com", "UA", "1.1.1.1")
		if err != nil {
			t.Fatalf("CreateSession %s: %v", sid, err)
		}
	}

	sessions, err := db.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions: %v", err)
	}
	if len(sessions) != 5 {
		t.Errorf("expected 5 sessions, got %d", len(sessions))
	}
}

func TestGenIndex(t *testing.T) {
	db := newTestDB(t)
	idx := db.genIndex("sessions", 42)
	if idx != "sessions:42" {
		t.Errorf("expected 'sessions:42', got '%s'", idx)
	}
}

func TestFlush(t *testing.T) {
	db := newTestDB(t)
	db.CreateSession("sid-flush", "test", "https://test.com", "UA", "1.1.1.1")
	db.Flush()

	sessions, _ := db.ListSessions()
	if len(sessions) != 1 {
		t.Error("flush should not remove data")
	}
}

func TestSetSessionUsernameNonexistent(t *testing.T) {
	db := newTestDB(t)
	err := db.SetSessionUsername("nonexistent", "admin")
	if err == nil {
		t.Error("expected error for nonexistent session")
	}
}

func TestDatabasePersistence(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "persist.db")

	db1, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	db1.CreateSession("persist-1", "test", "https://test.com", "UA", "1.1.1.1")
	db1.Flush()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("database file should exist")
	}

	db2, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase reopen: %v", err)
	}
	sessions, _ := db2.ListSessions()
	if len(sessions) != 1 {
		t.Errorf("expected 1 session after reopen, got %d", len(sessions))
	}
}
