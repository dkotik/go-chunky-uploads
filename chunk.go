package chunkyUploads

import (
	"bytes"
	"context"
	"io"
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
		ChunkStorageUsage(context.Context) (uint64, error)
	}

	// ChunkQuery(context.Context, *ChunkQuery) ([]*Chunk, error)
	// ChunkQuery struct {
	// 	Path    string
	// 	Status  Status
	// 	Size    Range
	// }
)

func (u *Uploads) append(ctx context.Context, f *File, b []byte) (int64, error) {
	size := int64(len(b))
	usage, err := u.chunks.ChunkStorageUsage(ctx)
	if usage+uint64(size) > u.chunkStorageLimit {
		return 0, ErrStorageFull
	}

	hash := u.hashProvider()
	_, err = io.Copy(hash, bytes.NewReader(b))
	if err != nil {
		return 0, err
	}

	return size, u.chunks.ChunkCreate(ctx, f,
		&Chunk{
			UUID:    u.uuidProvider(),
			Content: b,
			Hash:    hash.Sum(nil),
		},
	)
}
