package httpservice

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func Downloads(u chunkyUploads.Uploads) HTTPHandler {
	match := regexp.MustCompile(`^/(.{4,64})\.\w+$`)

	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		uuid := match.Find([]byte(r.URL.Path))
		if len(uuid) == 0 {
			return os.ErrNotExist
		}
		file, err := u.FileRetrieve(ctx, uuid)
		if err != nil {
			return err
		}

		header := w.Header()
		header.Set("Content-Type", file.ContentType)
		header.Set("Content-Disposition",
			fmt.Sprintf("Content-Disposition: attachment; filename=%q", file.Path))
		header.Set("Content-Length", fmt.Sprintf("%d", file.Size))
		return u.Copy(ctx, w, file.UUID)
	}
}
