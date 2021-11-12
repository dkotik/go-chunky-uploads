package chunkyUploads

import (
	"bytes"
	"context"
	"errors"
	"hash"
	"io"
)

var (
	ErrIncompleteCopy = errors.New("failed to copy all the bytes")
	ErrStorageFull    = errors.New("storage is overwhelmed")
)

type Uploads struct {
	files               FileRepository
	chunks              ChunkRepository
	hashProvider        func() hash.Hash
	uuidProvider        func() UUID
	contentTypeDetector ContentTypeDetector
	chunkSize           int
	chunkStorageLimit   uint64
}

func (u *Uploads) Copy(ctx context.Context, w io.Writer, uuid UUID) error {
	list, err := u.chunks.ChunkAttachmentList(ctx, uuid)
	if err != nil {
		return err
	}
	for _, one := range list {
		chunk, err := u.chunks.ChunkRetrieve(ctx, one.Chunk)
		if err != nil {
			return err
		}
		n, err := io.Copy(w, bytes.NewReader(chunk.Content))
		if err != nil {
			return err
		}
		if int(n) != len(chunk.Content) {
			return ErrIncompleteCopy
		}
	}
	return nil
}

func (u *Uploads) CopyRange(ctx context.Context, w io.Writer, uuid UUID, ra Range) error {
	list, err := u.chunks.ChunkAttachmentList(ctx, uuid)
	if err != nil {
		return err
	}
	left := ra.End - ra.Start
	for _, one := range list {
		if one.End < ra.Start {
			continue // chunk does not apply
		}
		chunk, err := u.chunks.ChunkRetrieve(ctx, one.Chunk)
		if err != nil {
			return err
		}
		n, err := io.Copy(w, io.LimitReader(bytes.NewReader(chunk.Content), left))
		if err != nil {
			return err
		}
		if int(n) != len(chunk.Content) {
			return ErrIncompleteCopy
		}
		left -= n
	}
	if left > 0 {
		return ErrIncompleteCopy
	}
	return nil
}
