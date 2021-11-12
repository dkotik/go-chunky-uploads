package chunkyUploads

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type (
	HTTPHandler     func(http.ResponseWriter, *http.Request) error
	HTTPFileLocator func(*http.Request) (*File, error)
)

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

		// TODO: use

		// func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker)

		header := w.Header()
		header.Set("Accept-Ranges", "bytes")
		header.Set("Content-Type", file.ContentType)
		header.Set("Content-Disposition",
			fmt.Sprintf("Content-Disposition: attachment; filename=%q", file.Path))
		header.Set("Content-Length", fmt.Sprintf("%d", file.Size))

		rangeHeader := r.Header.Get("Range")
		if rangeHeader == "" {
			return u.Copy(r.Context(), w, file.UUID)
		}

		header.Set("Content-Range", fmt.Sprintf("0-100/%d", file.Size))
		w.WriteHeader(http.StatusPartialContent)

		return u.Copy(r.Context(), w, file.UUID)
	}
}
