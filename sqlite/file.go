package sqlite

import (
	"context"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func (s *SQLiteDriver) FileDelete(ctx context.Context, uuid chunkyUploads.UUID) error {
	_, err := s.sqlFileDelete.ExecContext(ctx, uuid)
	return err
}
