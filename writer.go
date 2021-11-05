package chunkyUploads

import (
	"bytes"
	"context"
)

type writer struct {
	ctx     context.Context
	uploads *Uploads
	buffer  *bytes.Buffer
	file    *File
}

func (w *writer) Write(b []byte) (n int, err error) {
	remaining := w.buffer.Cap() - w.buffer.Len()
	defer func() {
		if err != nil {
			w.buffer.Truncate(w.buffer.Cap() - remaining)
			w.uploads.Update(w.ctx, w.file.UUID, func(f *File) error {
				f.Status = StatusError
				return nil
			})
		}
	}()

	if len(b) > remaining {
		b = b[:n]
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

func (w *writer) Flush() error {
	if w.buffer.Len() == 0 {
		return nil
	}

	err := w.uploads.PushChunk(w.ctx, w.file.UUID, w.buffer.Bytes())
	if err != nil {
        w.uploads.Update(w.ctx, w.file.UUID, func(f *File) error {
            f.Status = StatusError
            return nil
        })
		return err
	}
	w.buffer.Reset()
	return nil
}

func (w *writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	return w.uploads.Update(w.ctx, w.file.UUID, func(f *File) error {
		f.Status = StatusComplete
		return nil
	})
}
