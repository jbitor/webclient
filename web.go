package webclient

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
)

type T struct {
	peerSearches []*dht.GetPeersSearch
	dhtClient    dht.Client
	addr         string
	staticPath   string
}

func NewForDhtClient(dhtClient dht.Client) (wc T, err error) {
	wc.peerSearches = make([]*dht.GetPeersSearch, 0)
	wc.dhtClient = dhtClient

	wc.addr = "0.0.0.0:8080"

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

	http.HandleFunc("/api/peerSearch", wc.handlePeerSearch)

	http.Handle("/", http.StripPrefix(
		"/", http.FileServer(http.Dir(wc.staticPath))))

	return http.ListenAndServe(wc.addr, nil)
}

func (wc *T) serializePeerSearches() (serialized []interface{}) {
	serialized = make([]interface{}, 0)

	for _, search := range wc.peerSearches {
		peersFound := make([]interface{}, 0)
		for _, peer := range search.PeersFound {
			// TODO: serialize peers usefully
			peersFound = append(peersFound, peer.String())
		}

		serialized = append(serialized, map[string]interface{}{
			"infohash":          search.Infohash.String(),
			"searchDistance":    search.Infohash.BitDistance(wc.dhtClient.Id()),
			"peers":             peersFound,
			"finished":          search.Finished(),
			"queriedNodes":      wc.serializeQueriedNodes(search.QueriedNodes, search.Infohash),
			"queriedNodesCount": len(search.QueriedNodes),
		})
	}
	return
}

func (wc *T) serializeQueriedNodes(queriedNodes map[string]*dht.RemoteNode, target bittorrent.BTID) (serialized map[string]interface{}) {
	serialized = make(map[string]interface{}, 0)

	for key, node := range queriedNodes {
		var sourceId string
		if node.Source != nil {
			sourceId = node.Source.Id.String()
		}

		serialized[key] = map[string]interface{}{
			"id":             node.Id.String(),
			"sourceId":       sourceId,
			"localDistance":  node.Id.BitDistance(wc.dhtClient.Id()),
			"targetDistance": node.Id.BitDistance(target),
		}
	}

	return
}

func (wc *T) serialize() (serialized map[string]interface{}) {
	return map[string]interface{}{
		"dht": map[string]interface{}{
			"peerSearches":   wc.serializePeerSearches(),
			"connectionInfo": wc.dhtClient.ConnectionInfo(),
		},
	}
}

func (wc *T) handleClientState(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Serving clientState.json to %v.\n", r.RemoteAddr)

	s, err := json.Marshal(wc.serialize())
	if err != nil {
		logger.Fatalf("JSON encoding failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(s)
}

// XXX(JB): No CSRF protection or anything.
func (wc *T) handlePeerSearch(w http.ResponseWriter, r *http.Request) {
	hexInfohash := r.FormValue("infohash")
	infohash, err := bittorrent.BTIDFromHex(hexInfohash)
	if err != nil {
		panic("BTID shouldn't have been invalid!")
	}

	peerSearch := wc.dhtClient.GetPeers(infohash)

	wc.peerSearches = append(wc.peerSearches, peerSearch)
}
