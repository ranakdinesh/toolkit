package toolkit

import (
	"fmt"
	"net/http"
	"path"
)

// DownloadStaticFile download a static file and tries to force the browser to download it
func (t *Tools) DownloadStaticFile(w http.ResponseWriter, r *http.Request, p, file, displayName string) {
	fp := path.Join(p, file)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))
	http.ServeFile(w, r, fp)
}
