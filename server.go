package webclient

import (
	"encoding/json"
	"github.com/jbitor/dht"
	"net/http"
	"path/filepath"
)

func ServeForDhtClient(dhtClient dht.Client) (err error) {
	address := "127.0.0.1:47935"

	// XXX(JB): This is pretty horrible because it badly relies on the CWD.
	path, err := filepath.Abs("./src/github.com/jbitor/webclient/static/")
	if err != nil {
		return err
	}

	path = filepath.Clean(path)

	logger.Printf("Serving web client of %v at %v.\n", path, address)

	http.Handle("/_s/", http.StripPrefix(
		"/_s/", http.FileServer(http.Dir(path))))

	http.HandleFunc("/api/nodeCounts.json", func(w http.ResponseWriter, r *http.Request) {
		s, err := json.Marshal(dhtClient.ConnectionInfo())

		if err != nil {
			panic("How did JSON encoding possibly fail?")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Write(s)
	})

	http.Handle("/", http.StripPrefix(
		"/", http.FileServer(http.Dir(path))))

	return http.ListenAndServe(address, nil)
}
