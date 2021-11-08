package httpservice

import (
	"html/template"
	"net/http"

	chunkyUploads "github.com/dkotik/go-chunky-uploads"
)

type HTTPHandler func(http.ResponseWriter, *http.Request) error

func Index(u chunkyUploads.Uploads) HTTPHandler {
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
                <form
                  enctype="multipart/form-data"
                  action="/upload"
                  method="POST"
                >
                  <input class="input file-input" type="file" name="file" multiple />
                  <button class="button" type="submit">Upload</button>
                </form>
            </body>
        </html>
    `))
	uploads, downloads := Uploads(u, 16<<20), Downloads(u)

	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodPost {
			if err := uploads(w, r); err != nil {
				return err
			}
			w.WriteHeader(http.StatusTemporaryRedirect)
			// TODO: add redirect
			return nil
		}
		if r.URL.Path != "/" {
			return downloads(w, r)
		}

		files, err := u.FileQuery(r.Context(), &chunkyUploads.FileQuery{
			PerPage: 10000,
		})
		if err != nil {
			return err
		}
		return t.Execute(w, files)
	}
}
