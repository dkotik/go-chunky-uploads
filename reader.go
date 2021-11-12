package chunkyUploads

import (
	"bytes"
	"context"
	"io"
)

// Reader satisfies ReadSeekCloser interface.
func (u *Uploads) Reader(ctx context.Context, file *File) (*Reader, error) {

	return nil, nil
}

type Reader struct {
	uploads       *Uploads
	ctx           context.Context
	chunks        []*ChunkAttachment
	chunk         *Chunk
	size          int64
	cursor        int64
	cursorInChunk int64
}

// SeekStart means relative to the start of the file, SeekCurrent means relative to the current offset, and SeekEnd means relative to the end. Seek returns the new offset relative to the start of the file and an error, if any.
func (r *Reader) Seek(offset int64, whence int) (n int64, err error) {
	switch whence {
	case io.SeekCurrent:
		n = r.cursor + offset
	case io.SeekEnd:
		n = r.size - offset
	default: // io.SeekStart
		n = offset
	}
	if n < 0 || n > r.size {
		return 0, io.ErrUnexpectedEOF
	}

	for _, frame := range r.chunks {
		if n > frame.Start && n < frame.End {
			if r.chunk == nil || bytes.Compare(r.chunk.UUID, frame.Chunk) != 0 {
				r.chunk, err = r.uploads.chunks.ChunkRetrieve(r.ctx, frame.Chunk)
				if err != nil {
					return 0, err
				}
			}
			r.cursor = n
			r.cursorInChunk = n - frame.Start
			return n, nil
		}
	}
	return 0, io.EOF
}

func (r *Reader) Read(b []byte) (n int, err error) {
	if r.chunk == nil {
		if _, err = r.Seek(r.cursor, 0); err != nil {
			return 0, err
		}
	}
	n = copy(
		b,
		r.chunk.Content[r.cursorInChunk:],
	)
	r.cursor += int64(n)
	r.cursorInChunk += int64(n)
	if r.cursorInChunk > int64(len(r.chunk.Content)) {
		r.chunk = nil
	}
	return n, nil
}
