package httpservice

import (
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

type HTTPHandler func(http.ResponseWriter, *http.Request) error

func Index(u chunkyUploads.Uploads) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {

		return nil
	}
}
