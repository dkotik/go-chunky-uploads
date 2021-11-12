package httpservice

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func Uploads(u chunkyUploads.Uploads, field string, sizeLimit int64) HTTPHandler {
	saveOneFile := func(ctx context.Context, f *chunkyUploads.File, h *multipart.FileHeader) error {
		handle, err := h.Open()
		if err != nil {
			return err
		}
		defer handle.Close()
		return u.Save(ctx, f, handle)
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)
		defer r.Body.Close()
		err := r.ParseMultipartForm(sizeLimit)
		if err != nil {
			return fmt.Errorf("could not parse form: %w", err)
		}

		ctx := r.Context()
		files := r.MultipartForm.File[field]
		for _, file := range files {
			f := &chunkyUploads.File{
				Path:  file.Filename,
				Title: file.Filename,
				Size:  file.Size,
			}
			if err = saveOneFile(ctx, f, file); err != nil {
				return err
			}
		}
		return nil
	}
}
