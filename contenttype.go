package chunkyUploads

import (
	"errors"
	"io"
	"net/http"
)

var (
	ErrContentTypeNotAllowed = errors.New("detected content type is not allowed")
)

// ContentTypeDetector reads the head of a ReadSeeker and attempts to determine the underlying content type.
type ContentTypeDetector func(io.ReadSeeker) (string, error)

// NewContentTypeDetector returns a function that returns a content type, if it is in the allowed list. If the allowed list is `nil`, all content types are accepted. Inherits "application/octet-stream" as default.
func NewContentTypeDetector(allowed ...string) ContentTypeDetector {
	b := make([]byte, 512)

	if allowed == nil { // anything goes
		return func(r io.ReadSeeker) (string, error) {
			_, err := r.Read(b)
			if err != nil {
				return "", err
			}
			_, err = r.Seek(0, io.SeekStart)
			if err != nil {
				return "", err
			}
			return http.DetectContentType(b), nil
		}
	}

	return func(r io.ReadSeeker) (string, error) {
		_, err := r.Read(b)
		if err != nil {
			return "", err
		}
		_, err = r.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}
		t := http.DetectContentType(b)
		for _, one := range allowed {
			if one == t {
				return t, nil
			}
		}
		return "", ErrContentTypeNotAllowed
	}
}
