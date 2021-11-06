package sqlite

import (
	"context"
	"database/sql"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func (s *SQLiteDriver) ChunkCreate(ctx context.Context, f *chunkyUploads.File, c *chunkyUploads.Chunk) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	create := tx.Stmt(s.sqlChunkCreate)
	defer create.Close()

	_, err = create.ExecContext(ctx, c.UUID, c.Hash, c.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteDriver) ChunkDelete(ctx context.Context, uuid chunkyUploads.UUID) error {
	_, err := s.sqlChunkDelete.ExecContext(ctx, uuid)
	return err
}
