package compressed

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

var _ chunkyUploads.ChunkRepository = (*Gzip)(nil)

type Gzip struct {
	chunkyUploads.ChunkRepository
}

func (g *Gzip) ChunkCreate(
	ctx context.Context,
	f *chunkyUploads.File,
	c *chunkyUploads.Chunk,
) error {
	var (
		b bytes.Buffer
		w = gzip.NewWriter(&b)
	)
	defer w.Close()

	_, err := io.Copy(w, bytes.NewReader(c.Content))
	if err != nil {
		return err
	}
	return g.ChunkRepository.ChunkCreate(ctx, f, &chunkyUploads.Chunk{
		UUID:    c.UUID,
		Hash:    c.Hash, // todo: update hash? yes
		Content: b.Bytes(),
	})
}

func (g *Gzip) ChunkRetrieve(ctx context.Context, uuid chunkyUploads.UUID) (*chunkyUploads.Chunk, error) {
	c, err := g.ChunkRepository.ChunkRetrieve(ctx, uuid)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(c.Content))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	_, err = io.Copy(&b, r)
	if err != nil {
		return nil, err
	}

	return &chunkyUploads.Chunk{
		UUID:    c.UUID,
		Hash:    c.Hash,
		Content: b.Bytes(),
	}, nil
}
