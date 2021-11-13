package sqlite

import (
	"context"
	"testing"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func testFileOperations(t *testing.T, db *SQLiteDriver) {
	ctx := context.Background()
	uuid := chunkyUploads.UUID(`test`)

	err := db.FileCreate(ctx, &chunkyUploads.File{
		UUID: uuid,
		Hash: []byte(`test`),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.FileUpdate(ctx, uuid, func(f *chunkyUploads.File) error {
		f.Title = "updated title"
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
