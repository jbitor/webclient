package webclient

import (
	"encoding/json"
	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
	"net/http"
	"path/filepath"
)

type peerRequest struct {
	HexInfohash string
	infohash    bittorrent.BTID
	Peers       []*bittorrent.RemotePeer
}

type T struct {
	peerRequests []*peerRequest
	dhtClient    dht.Client
	addr         string
	staticPath   string
}

func NewForDhtClient(dhtClient dht.Client) (wc T, err error) {
	wc.peerRequests = make([]*peerRequest, 0)
	wc.dhtClient = dhtClient

	wc.addr = "127.0.0.1:47935"

	// XXX(JB): This is pretty horrible because it badly relies on the CWD.
	path, err := filepath.Abs("./src/github.com/jbitor/webclient/static/")
	if err != nil {
		return
	}

	wc.staticPath = filepath.Clean(path)
	return
}

func (wc *T) ListenAndServe() (err error) {
	logger.Printf("Serving web client of %v at %v.\n",
		wc.staticPath, wc.addr)

	http.Handle("/_s/", http.StripPrefix(
		"/_s/", http.FileServer(http.Dir(wc.staticPath))))

	http.HandleFunc("/api/clientState.json", wc.handleClientState)

	http.HandleFunc("/api/peerRequest", wc.handlePeerRequest)

	http.Handle("/", http.StripPrefix(
		"/", http.FileServer(http.Dir(wc.staticPath))))

	return http.ListenAndServe(wc.addr, nil)
}

func (wc *T) handleClientState(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Serving clientState.json to %v.\n", r.RemoteAddr)

	s, err := json.Marshal(map[string]interface{}{
		"peerRequests": wc.peerRequests,
		"nodeCounts":   wc.dhtClient.ConnectionInfo(),
	})
	if err != nil {
		logger.Fatalf("How did JSON encoding possibly fail? %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(s)
}

// XXX(JB): No CSRF protection or anything.
func (wc *T) handlePeerRequest(w http.ResponseWriter, r *http.Request) {
	hexInfohash := r.FormValue("infohash")
	infohash, err := bittorrent.BTIDFromHex(hexInfohash)
	if err != nil {
		panic("BTID shouldn't have been invalid!")
	}

	peerRequest := peerRequest{
		hexInfohash,
		infohash,
		nil,
	}

	wc.peerRequests = append(wc.peerRequests, &peerRequest)

	go func() {
		peers, err := wc.dhtClient.GetPeers(infohash)
		if err != nil {
			// XXX(JB): obviouly a bad way of handling errors
			peerRequest.HexInfohash = "error"
			return
		}

		logger.Printf("Got peers for %v: %v.",
			peerRequest.infohash, peers)

		peerRequest.Peers = peers
	}()
}
