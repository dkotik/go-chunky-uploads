package chunkyUploads

import (
	"bytes"
	"context"
	"io"
)

func (u *Uploads) append(ctx context.Context, f *File, b []byte) (int64, error) {
	size := int64(len(b))
	usage, err := u.ChunkRepository.ChunkStorageUsage(ctx)
	if usage+uint64(size) > u.chunkStorageLimit {
		return 0, ErrStorageFull
	}

	hash := u.hashProvider()
	_, err = io.Copy(hash, bytes.NewReader(b))
	if err != nil {
		return 0, err
	}

	return size, u.ChunkRepository.ChunkCreate(ctx, f,
		&Chunk{
			UUID:    u.uuidProvider(),
			Content: b,
			Hash:    hash.Sum(nil),
		},
	)
}

func (u *Uploads) Save(ctx context.Context, f *File, r io.ReadSeeker) (err error) {
	if f.ContentType == "" {
		f.ContentType, err = u.contentTypeDetector(r)
		if err != nil {
			return err
		}
	}
	f.UUID = u.uuidProvider()
	f.Status = StatusUploading
	if err = u.FileRepository.FileCreate(ctx, f); err != nil {
		return err
	}
	defer u.FileRepository.FileUpdate(ctx, f.UUID, func(nf *File) error {
		nf.Size = f.Size
		if err != nil {
			nf.Status = StatusError
		} else {
			nf.Status = StatusComplete
		}
		return nil
	})

	b := &bytes.Buffer{}
	b.Grow(u.chunkSize)

	for err == nil {
		_, err = io.CopyN(b, r, int64(u.chunkSize))
		if err != nil {
			return err
		}
		n, err := u.append(ctx, f, b.Bytes())
		if err != nil {
			return err
		}
		f.Size += n
	}

	if err == io.EOF {
		err = nil
	}
	return err
}
