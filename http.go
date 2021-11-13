package chunkyUploads

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
)

type (
	HTTPHandler     func(http.ResponseWriter, *http.Request) error
	HTTPFileLocator func(*http.Request) (*File, error)
)

func (u *Uploads) Upload(field string, sizeLimit int64) HTTPHandler {
	saveOneFile := func(ctx context.Context, f *File, h *multipart.FileHeader) error {
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
			f := &File{
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

func (u *Uploads) FileByUUID() HTTPFileLocator {
	match := regexp.MustCompile(`([^.\/]{4,64})(\.\w+)?$`)
	return func(r *http.Request) (*File, error) {
		uuid, err := base64.StdEncoding.DecodeString(match.FindString(r.URL.Path))
		if len(uuid) == 0 || err != nil {
			return nil, os.ErrNotExist
		}
		file, err := u.files.FileRetrieve(r.Context(), uuid)
		if err != nil {
			return nil, os.ErrNotExist // TODO: check if actually not found or another error
		}
		return file, nil
	}
}

func (u *Uploads) Download(using HTTPFileLocator) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		file, err := using(r)
		if err != nil {
			return err
		}
		reader, err := u.Reader(r.Context(), file)
		if err != nil {
			return err
		}

		header := w.Header()
		// header.Set("Accept-Ranges", "bytes")
		header.Set("ETag",
			base64.StdEncoding.EncodeToString(file.Hash))
		header.Set("Content-Type", file.ContentType)
		header.Set("Content-Disposition",
			fmt.Sprintf("Content-Disposition: attachment; filename=%q", file.Path))
		header.Set("Content-Length", fmt.Sprintf("%d", file.Size))

		http.ServeContent(
			w,
			r,
			path.Base(file.Path),
			file.UpdatedAt,
			reader,
		)
		return nil

		// rangeHeader := r.Header.Get("Range")
		// if rangeHeader == "" {
		// 	return u.Copy(r.Context(), w, file.UUID)
		// }
		// header.Set("Content-Range", fmt.Sprintf("0-100/%d", file.Size))
		// w.WriteHeader(http.StatusPartialContent)
		// return u.Copy(r.Context(), w, file.UUID)
	}
}
