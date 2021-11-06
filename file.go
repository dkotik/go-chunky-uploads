package chunkyUploads

import "context"

type (
	UUID   []byte
	Hash   []byte
	Range  struct{ Start, End int64 }
	Status uint8

	File struct {
		UUID                            UUID
		Hash                            Hash
		Path                            string
		Title                           string
		Description                     string
		ContentType                     string
		Status                          Status
		Size                            int64
		CreatedAt, UpdatedAt, DeletedAt int64
	}

	FileRepository interface {
		FileCreate(context.Context, *File) error
		FileRetrieve(context.Context, UUID) (*File, error)
		FileUpdate(context.Context, UUID, func(*File) error) error
		FileDelete(context.Context, UUID) error

		FileQuery(context.Context, *FileQuery) ([]*File, error)
	}

	FileQuery struct {
		Path    string
		Status  Status
		Size    *Range
		Created *Range
		Updated *Range
		Deleted *Range
	}
)

const (
	StatusUnkown = iota
	StatusUploading
	StatusCancelled
	StatusError
	StatusComplete
	StatusDeleted
)
