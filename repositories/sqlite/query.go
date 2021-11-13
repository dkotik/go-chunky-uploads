package sqlite

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func queryOrder(q *chunkyUploads.FileQuery) string {
	order := "created_at"
	switch q.OrderBy {
	case chunkyUploads.QueryOrderBySize:
		order = "size"
	}
	if q.Descending {
		return order + " DESC"
	}
	return order
}

func queryFilter(q *chunkyUploads.FileQuery, s sq.SelectBuilder) sq.SelectBuilder {
	if q.Path != "" {
		s.Where("path LIKE ?", "%"+q.Path+"%")
	}
	if q.Status > 0 {
		s.Where("status = ?", q.Status)
	}
	if q.Size != nil {
		s = s.
			Where("size > ?", q.Size.Start).
			Where("size < ?", q.Size.End)
	}
	if q.CreatedAt != nil {
		s = s.
			Where("created_at > ?", q.CreatedAt.Start).
			Where("created_at < ?", q.CreatedAt.End)
	}
	if q.UpdatedAt != nil {
		s = s.
			Where("updated_at > ?", q.UpdatedAt.Start).
			Where("updated_at < ?", q.UpdatedAt.End)
	}
	if q.DeletedAt != nil {
		s = s.
			Where("deleted_at > ?", q.DeletedAt.Start).
			Where("deleted_at < ?", q.DeletedAt.End)
	}
	return s.OrderBy(queryOrder(q))
}

func (s *SQLiteDriver) FileQuery(ctx context.Context, q *chunkyUploads.FileQuery) (*chunkyUploads.FileQueryResult, error) {
	// q.CreatedAt = &chunkyUploads.Range{0, 0}

	start := uint64(q.PerPage) * uint64(q.Page-1)
	countq := queryFilter(q, sq.Select("COUNT(uuid)")).
		From(s.tableFiles).
		RunWith(s.db)

	var count sql.NullInt64
	err := countq.QueryRowContext(ctx).Scan(&count)
	if err != nil {
		return nil, err
	}
	final := &chunkyUploads.FileQueryResult{
		Files: make([]*chunkyUploads.File, 0),
		Start: start,
		End:   start,
	}
	if count.Int64 == 0 {
		return final, nil
	}

	query := queryFilter(q, sq.Select("uuid, hash, path, title, description, content_type, status, size, created_at, updated_at, deleted_at")).
		Limit(uint64(q.PerPage)).
		Offset(uint64(q.PerPage) * uint64(q.Page-1)).
		From(s.tableFiles).
		RunWith(s.db)
	// log.Println(query.ToSql())
	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result := &chunkyUploads.File{}
		var CreatedAt, UpdatedAt, DeletedAt int64
		if err = rows.Scan(
			&result.UUID,
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
		final.Files = append(final.Files, result)
		final.End++
	}
	// if err = rows

	return final, nil
}
