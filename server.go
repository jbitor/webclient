package webclient

import (
	"github.com/jbitor/dht"
	"net/http"
	"path/filepath"
)

func ServeForDhtClient(dhtClient dht.Client) (err error) {
	address := "127.0.0.1:8080"

	// XXX(JB): This is pretty horrible because it badly relies on the CWD.
	path, err := filepath.Abs("./src/github.com/jbitor/webclient/static/")
	if err != nil {
		return err
	}

	path = filepath.Clean(path)

	logger.Printf("Serving web client of %v at %v.\n", path, address)

	http.Handle("/_s/", http.StripPrefix(
		"/_s/", http.FileServer(http.Dir(path))))

	http.Handle("/", http.StripPrefix(
		"/", http.FileServer(http.Dir(path))))

	return http.ListenAndServe(address, nil)
}
