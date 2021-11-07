package sqlite

import (
	"context"
	"testing"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func testChunkOperations(t *testing.T, db *SQLiteDriver) {
	ctx := context.Background()
	uuid := chunkyUploads.UUID(`test`)

	err := db.ChunkCreate(ctx, &chunkyUploads.File{
		UUID: uuid,
		Hash: []byte(`test`),
	}, &chunkyUploads.Chunk{
		UUID:    uuid,
		Hash:    []byte(`test`),
		Content: []byte(`test`),
	})
	if err != nil {
		t.Fatal(err)
	}

	n, err := db.ChunkStorageUsage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("there is no disk usage for chunks")
	}
	if n != 4 {
		t.Fatal("unexpected chunk sum size", n)
	}
}
