package chunkyUploads

import (
	"bytes"
	"context"
	"hash"
	"io"
	"os"
)

func (u *Uploads) Delete(ctx context.Context, file *File) error {
	if file == nil {
		return os.ErrNotExist
	}
	return nil
}

func (u *Uploads) Writer(ctx context.Context, file *File) (io.WriteCloser, error) {
	if file == nil {
		return nil, os.ErrNotExist
	}

	file.Status = StatusUploading
	if err := u.files.FileCreate(ctx, file); err != nil {
		return nil, err
	}

	w := &writer{
		ctx:    ctx,
		file:   file,
		buffer: &bytes.Buffer{},
		hash:   u.hashProvider(),
	}
	w.buffer.Grow(u.chunkSize)
	return w, nil
}

type writer struct {
	ctx     context.Context
	uploads *Uploads
	buffer  *bytes.Buffer
	hash    hash.Hash
	file    *File
}

func (w *writer) Write(b []byte) (n int, err error) {
	remaining := w.buffer.Cap() - w.buffer.Len()
	defer func() {
		if err != nil {
			w.buffer.Truncate(w.buffer.Cap() - remaining)
			w.uploads.files.FileUpdate(w.ctx, w.file.UUID, func(f *File) error {
				f.Status = StatusError
				return nil
			})
		}
	}()

	if len(b) > remaining {
		b = b[:remaining]
	}
	n, err = w.buffer.Write(b)
	if err != nil {
		return 0, err
	}
	if remaining-n == 0 {
		if err = w.Flush(); err != nil {
			return 0, err
		}
	}
	return
}

func (w *writer) Flush() (err error) {
	if w.buffer.Len() == 0 {
		return nil
	}
	hash := w.uploads.hashProvider()
	_, err = io.Copy(hash, w.buffer)
	if err != nil {
		return err
	}
	chunk := &Chunk{
		UUID:    w.uploads.uuidProvider(),
		Content: w.buffer.Bytes(),
		Hash:    hash.Sum(nil),
	}
	if err = w.uploads.chunks.ChunkCreate(w.ctx, w.file, chunk); err != nil {
		w.uploads.files.FileUpdate(w.ctx, w.file.UUID, func(f *File) error {
			f.Status = StatusError
			return nil
		})
		return err
	}
	_, err = io.Copy(w.hash, w.buffer)
	if err != nil {
		return err
	}
	w.buffer.Reset()
	return nil
}

func (w *writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.uploads.files.FileUpdate(w.ctx, w.file.UUID, func(f *File) error {
		f.Hash = w.hash.Sum(nil)
		f.Status = StatusComplete
		return nil
	})
}
