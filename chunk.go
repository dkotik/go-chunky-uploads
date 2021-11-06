package chunkyUploads

import (
	"context"
)

type (
	Chunk struct {
		UUID    UUID
		Hash    Hash
		Content []byte
	}

	ChunkAttachment struct {
		File  UUID
		Chunk UUID
		Start int64
		End   int64
	}

	ChunkRepository interface {
		ChunkCreate(context.Context, *File, *Chunk) error
		ChunkRetrieve(context.Context, UUID) (*Chunk, error)
		ChunkDelete(context.Context, UUID) error
		ChunkAttachmentList(context.Context, UUID) ([]*ChunkAttachment, error)
	}

	// ChunkQuery(context.Context, *ChunkQuery) ([]*Chunk, error)
	// ChunkQuery struct {
	// 	Path    string
	// 	Status  Status
	// 	Size    Range
	// }
)
