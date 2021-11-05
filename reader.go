package chunkyUploads

import "context"

type reader struct {
	ctx     context.Context
	uploads *Uploads
	file    *File
	chunks  []*ChunkAttachment
	current *ChunkAttachment
	cursor  int
}

// func (r *Reader) Read([]byte) (n int, err error) {
//
// }
