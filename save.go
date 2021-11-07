package chunkyUploads

import (
	"context"
	"io"
)

func (u *Uploads) Save(ctx context.Context, f *File, r io.Reader) (int64, error) {
	// upload chunks first?
	return 0, nil
}
