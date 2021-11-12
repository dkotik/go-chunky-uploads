package chunkyUploads

import (
	"bytes"
	"context"
	"io"
	"time"
)

type (
	UUID       []byte
	Hash       []byte
	Range      struct{ Start, End int64 }
	Status     uint8
	QueryOrder uint8

	File struct {
		UUID                            UUID
		Hash                            Hash
		Path                            string
		Title                           string
		Description                     string
		ContentType                     string
		Status                          Status
		Size                            int64
		CreatedAt, UpdatedAt, DeletedAt time.Time
	}

	FileRepository interface {
		FileCreate(context.Context, *File) error
		FileRetrieve(context.Context, UUID) (*File, error)
		FileUpdate(context.Context, UUID, func(*File) error) error
		FileDelete(context.Context, UUID) error

		FileQuery(context.Context, *FileQuery) (*FileQueryResult, error)
	}

	FileQuery struct {
		Page       uint32
		PerPage    uint32
		OrderBy    QueryOrder
		Descending bool
		Status     Status
		Path       string
		Size       *Range
		CreatedAt  *Range
		UpdatedAt  *Range
		DeletedAt  *Range
	}

	FileQueryResult struct {
		Files             []*File
		Start, End, Total uint64
	}
)

const (
	StatusUnkown Status = iota
	StatusUploading
	StatusCancelled
	StatusError
	StatusComplete
	StatusDeleted

	QueryOrderByUnknown QueryOrder = iota
	QueryOrderByStatus
	QueryOrderBySize
	QueryOrderByCreated
	QueryOrderByUpdated
	QueryOrderByDeleted
)

func (u *Uploads) Save(ctx context.Context, f *File, r io.ReadSeeker) (err error) {
	if f.ContentType == "" {
		f.ContentType, err = u.contentTypeDetector(r)
		if err != nil {
			return err
		}
	}
	f.UUID = u.uuidProvider()
	f.Status = StatusUploading
	if err = u.files.FileCreate(ctx, f); err != nil {
		return err
	}
	defer u.files.FileUpdate(ctx, f.UUID, func(nf *File) error {
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
