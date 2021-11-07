package sqlite

import (
	"context"
	"testing"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func testQueryOperations(t *testing.T, db *SQLiteDriver) {
	ctx := context.Background()

	result, err := db.FileQuery(ctx, &chunkyUploads.FileQuery{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != 1 {
		t.Fatal("there should be one row in the database")
	}
	// t.Logf("%+v", result.Files[0])
	// t.Fatal()
}
