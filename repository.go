package chunkyUploads

import (
	"context"
)

type (
	Repository interface {
		Create(context.Context, *File) error
		Retrieve(context.Context, UUID) (*File, error)
		Update(context.Context, UUID, func(*File) error) error
		Delete(context.Context, UUID) error

		Query(context.Context, *Query) ([]*File, error)
		Chunks(context.Context, UUID) ([]*ChunkAttachment, error)
		PushChunk(context.Context, UUID, []byte) error
		PullChunk(context.Context, UUID) ([]byte, error)
	}

	Range struct{ Start, End int64 }

	Query struct {
		Path    string
		Status  Status
		Size    *Range
		Created *Range
		Updated *Range
		Deleted *Range
	}
)
