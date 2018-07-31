package db

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"regexp"
)

// NamePattern general name pattern
var NamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\._\-]{3,}$`)

// WildcardPattern wildcard pattern
var WildcardPattern = regexp.MustCompile(`^[a-zA-Z0-9\._\*\-]*$`)

type DB struct {
	*gorm.DB
}

// Open open or create a database file
func Open(filename string) (db *DB, err error) {
	// try make parent directories
	os.MkdirAll(filepath.Dir(filename), 0640)
	// open database
	var d *gorm.DB
	if d, err = gorm.Open("sqlite3", filename); err != nil {
		return
	}
	// construct wrapper struct
	db = &DB{DB: d}
	return
}

func (db *DB) Migrate() error {
	return db.AutoMigrate(new(User)).Error
}

func IsRecordNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
