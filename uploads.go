package chunkyUploads

import (
	"bytes"
	"context"
	"io"
	"os"
)

type Uploads struct {
	Repository

	ChunkSize int
}

func (u *Uploads) Writer(ctx context.Context, file *File) (io.WriteCloser, error) {
	if file == nil {
		return nil, os.ErrNotExist
	}

	file.Status = StatusUploading
	if err := u.Repository.Create(ctx, file); err != nil {
		return nil, err
	}

	w := &writer{
		ctx:    ctx,
		file:   file,
		buffer: &bytes.Buffer{},
	}
	w.buffer.Grow(u.ChunkSize)
	return w, nil
}
