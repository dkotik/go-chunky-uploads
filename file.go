package chunkyUploads

import (
	"context"
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
		PerPage    uint8
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
