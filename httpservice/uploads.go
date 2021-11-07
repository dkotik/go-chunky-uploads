package httpservice

import (
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func Upload(u chunkyUploads.Uploads) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {

		return nil
	}
}

func Download(u chunkyUploads.Uploads) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {

		return nil
	}
}
