package sqlite

import (
	"path"
	"testing"
)

func TestSQliteDriver(t *testing.T) {
	db, err := NewSQLiteDriver("test", path.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			t.Fatal(err)
		}
	}()
}
