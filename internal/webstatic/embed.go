package webstatic

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"time"
)

var (
	//go:embed static
	staticFS    embed.FS
	StaticFS, _ = fs.Sub(staticFS, "static")
)

func Register(mux *http.ServeMux) {
	t := time.Now()
	fs.WalkDir(StaticFS, ".", func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		mux.HandleFunc("/"+p, func(w http.ResponseWriter, r *http.Request) {
			f, err := StaticFS.Open(p)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			rs, ok := f.(io.ReadSeeker)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.ServeContent(w, r, d.Name(), t, rs)
		})
		return nil
	})
}
