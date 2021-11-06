package sqlite

import (
	"context"
	"database/sql"
	"time"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func (s *SQLiteDriver) ChunkCreate(ctx context.Context, f *chunkyUploads.File, c *chunkyUploads.Chunk) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
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

	attach := tx.Stmt(s.sqlChunkCreateAttachment)
	defer attach.Close()
	growth := f.Size + int64(len(c.Content))
	_, err = attach.ExecContext(ctx, f.UUID, c.UUID, f.Size, growth)
	if err != nil {
		return err
	}

	update := tx.Stmt(s.sqlFileUpdate)
	defer update.Close()
	_, err = update.ExecContext(
		ctx,
		f.Hash,
		f.Path,
		f.Title,
		f.Description,
		f.ContentType,
		f.Status,
		growth,
		time.Now().Unix(),
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SQLiteDriver) ChunkRetrieve(ctx context.Context, uuid chunkyUploads.UUID) (*chunkyUploads.Chunk, error) {
	row := s.sqlChunkRetrieve.QueryRowContext(ctx, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	result := &chunkyUploads.Chunk{UUID: uuid}
	if err = row.Scan(&result.Hash, &result.Content); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SQLiteDriver) ChunkAttachmentList(ctx context.Context, uuid chunkyUploads.UUID) ([]*chunkyUploads.ChunkAttachment, error) {
	rows, err := s.sqlChunkAttachmentList.QueryContext(ctx, uuid)
	if err != nil {
		return nil, err
	}

	result := make([]*chunkyUploads.ChunkAttachment, 0)
	for rows.Next() {
		attachment := &chunkyUploads.ChunkAttachment{File: uuid}
		if err = rows.Scan(&attachment.Chunk, &attachment.Start, &attachment.End); err != nil {
			return nil, err
		}
		result = append(result, attachment)
	}
	return result, nil
}

func (s *SQLiteDriver) ChunkDelete(ctx context.Context, uuid chunkyUploads.UUID) error {
	_, err := s.sqlChunkDelete.ExecContext(ctx, uuid)
	return err
}
