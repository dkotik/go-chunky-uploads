package main

import (
	"fmt"
	"html/template"
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

func Index(u chunkyUploads.Uploads) chunkyUploads.HTTPHandler {
	t := template.Must(template.New("index").Parse(`
        <html>
            <head>
                <title>Downloads</title>
            </head>
            <body>
                <ul>
                    {{ range . }}
                    <li>
                        {{ .Path }}
                        <a href="{{ .UUID }}">↴</a>
                    </li>
                    {{ end }}
                </ul>
                <form enctype="multipart/form-data" action="/upload" method="POST">
                  <input class="input file-input" type="file" name="file" multiple />
                  <button class="button" type="submit">Upload</button>
                </form>
            </body>
        </html>
    `))
	uploads, downloads := u.Upload("upload", 16<<20), u.Download(u.FileByUUID())

	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodPost {
			if err := uploads(w, r); err != nil {
				return err
			}
			http.Redirect(w, r, r.URL.Path, http.StatusTemporaryRedirect)
			return nil
		}
		if r.URL.Path != "/" {
			return downloads(w, r)
		}

		// files, err := u.FileQuery(r.Context(), &chunkyUploads.FileQuery{
		// 	PerPage: 10000,
		// })
		// if err != nil {
		// 	return err
		// }
		// return t.Execute(w, files)
		fmt.Println(t)
		return nil
	}
}

func main() {

}
