package sqlite

import (
	"context"
	"time"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func (s *SQLiteDriver) FileCreate(ctx context.Context, f *chunkyUploads.File) error {
	t := time.Now().Unix()
	_, err := s.sqlFileCreate.ExecContext(
		ctx,
		&f.UUID,
		&f.Hash,
		&f.Path,
		&f.Title,
		&f.Description,
		&f.ContentType,
		&f.Status,
		&t, &t,
	)
	return err
}

func (s *SQLiteDriver) FileRetrieve(ctx context.Context, uuid chunkyUploads.UUID) (*chunkyUploads.File, error) {
	row := s.sqlFileRetrieve.QueryRowContext(ctx, uuid)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	result := &chunkyUploads.File{UUID: uuid}
	var CreatedAt, UpdatedAt, DeletedAt int64
	if err = row.Scan(
		&result.Hash,
		&result.Path,
		&result.Title,
		&result.Description,
		&result.ContentType,
		&result.Status,
		&result.Size,
		&CreatedAt,
		&UpdatedAt,
		&DeletedAt,
	); err != nil {
		return nil, err
	}
	result.CreatedAt = time.Unix(CreatedAt, 0)
	result.UpdatedAt = time.Unix(UpdatedAt, 0)
	result.DeletedAt = time.Unix(DeletedAt, 0)
	return result, nil
}

func (s *SQLiteDriver) FileUpdate(ctx context.Context, uuid chunkyUploads.UUID, update func(f *chunkyUploads.File) error) error {
	f, err := s.FileRetrieve(ctx, uuid)
	if err != nil {
		return err
	}
	hash := f.Hash
	size := f.Size
	if err = update(f); err != nil {
		return err
	}
	_, err = s.sqlFileUpdate.ExecContext(
		ctx,
		hash,
		f.Path,
		f.Title,
		f.Description,
		f.ContentType,
		f.Status,
		size,
		time.Now().Unix(),
	)
	return err
}

func (s *SQLiteDriver) FileDelete(ctx context.Context, uuid chunkyUploads.UUID) error {
	_, err := s.sqlFileDelete.ExecContext(ctx, uuid)
	return err
}
