package chunkyUploads

type (
	UUID   []byte
	Hash   []byte
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
)

const (
	StatusUnkown = iota
	StatusUploading
	StatusCancelled
	StatusError
	StatusComplete
	StatusDeleted
)
