package httpservice

import (
	"io"
	"net/http"
)

func DetectContentType(r io.ReadSeeker) (string, error) {
	b := make([]byte, 512)
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
