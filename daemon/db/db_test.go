package db

import (
	"testing"
	"crypto/rand"
	"os"
	"path/filepath"
	"encoding/hex"
)

func temporaryFile() string {
	buf := make([]byte, 8, 8)
	rand.Read(buf)
	return filepath.Join(os.TempDir(), "bnktestdb"+hex.EncodeToString(buf)+".sqlite3")
}

func TestDB_Migrate(t *testing.T) {
	db, err := Open(temporaryFile())
	db.LogMode(true)
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Migrate(); err != nil {
		t.Fatal(err)
	}
	u := &User{
		Account: "testuser",
		IsAdmin: 1,
	}
	u.SetPassword("qwerty")
	if err = db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	id := u.ID
	u.ID = 0
	if err = db.Create(u).Error; err == nil {
		t.Fatal("unique index failed")
	}
	m := User{}
	if err = db.Find(&m, id).Error; err != nil {
		t.Fatal("failed to retrieve")
	}
	if m.IsAdmin == 0 {
		t.Fatal("failed to retrieve bool")
	}
}
