package httpservice

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func Uploads(u chunkyUploads.Uploads, sizeLimit int64) HTTPHandler {
	save := func(ctx context.Context, f *chunkyUploads.File, h *multipart.FileHeader) error {
		handle, err := h.Open()
		if err != nil {
			return err
		}
		defer handle.Close()

		f.ContentType, err = DetectContentType(handle)
		if err != nil {
			return err
		}
		if err = u.FileCreate(ctx, f); err != nil {
			return err
		}
		_, err = u.Save(ctx, f, handle)
		if err != nil {
			return err
		}
		// f.Size = n
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)
		defer r.Body.Close()
		err := r.ParseMultipartForm(sizeLimit)
		if err != nil {
			return fmt.Errorf("could not parse form: %w", err)
		}

		ctx := r.Context()
		files := r.MultipartForm.File["upload"]
		for _, file := range files {
			f := &chunkyUploads.File{
				Path:  file.Filename,
				Title: file.Filename,
				Size:  file.Size,
			}
			if err = save(ctx, f, file); err != nil {
				return err
			}
		}

		return nil
	}
}
