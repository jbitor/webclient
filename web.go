package webclient

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jbitor/bencoding"
	"github.com/jbitor/bittorrent"
	"github.com/jbitor/dht"
)

type T struct {
	dhtClient    dht.Client
	btClient     bittorrent.Client
	addr         string
	staticPath   string
	pageTemplate string
}

func New(dhtClient dht.Client, btClient bittorrent.Client) (wc T, err error) {
	wc.dhtClient = dhtClient
	wc.btClient = btClient

	wc.addr = "0.0.0.0:8080"

	// HACK
	path, err := filepath.Abs("./src/github.com/jbitor/webclient/static/")
	if err != nil {
		return
	}

	wc.staticPath = filepath.Clean(path)

	templatePath := wc.staticPath + "/index.html"
	data, err := ioutil.ReadFile(templatePath)
	if err != nil {
		logger.Fatalf("unable to read template: %v", err)
	}
	wc.pageTemplate = string(data)

	return
}

func (wc *T) handleRequest(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Got request")

	btih, _ := bittorrent.BTIDFromHex("0ea39049afdbaaea255ca1d0af662e2a0d503098")

	if r.URL.Path == "/" {
		wc.handleIndex(w, r)
		return
	}

	path := r.URL.Path[1:]
	pieces := strings.Split(path, ".")
	btih, _ = bittorrent.BTIDFromHex(pieces[0])

	if len(pieces) == 1 {
		wc.handleTorrentPageRequest(w, r, btih)
		return
	} else if len(pieces) == 2 {
		extension := pieces[1]

		switch extension {
		case "torrent":
			wc.handleTorrentFileRequest(w, r, btih)
			return
		default:
			logger.Error("404 extension %v", extension)
			return
		}
	} else {
		logger.Error("404 path components %v", path)
		return
	}

}

func (wc *T) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Write([]byte(wc.pageTemplate))
}

// TODO: use these

// PERMANENT REDIRECT /?btih=INFOHASH -> /INFOHASH
func (wc *T) handleTorrentPageQueryRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	http.Redirect(w, r, "/"+r.URL.Query().Get("btih"), 301)
}

// PERMANENT REDIRECT /.json?btih=INFOHASH -> /INFOHASH.json
func (wc *T) handleTorrentJsonInfoQueryRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	http.Redirect(w, r, "/.json"+r.URL.Query().Get("btih")+".json", 301)
}

// PERMANENT REDIRECT /.torrent?btih=INFOHASH -> /INFOHASH.torrent
func (wc *T) handleTorrentQueryRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	http.Redirect(w, r, "/.torrent"+r.URL.Query().Get("btih")+".torrent", 301)
}

// /INFOHASH
//
func (wc *T) handleTorrentPageRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	// Use a Refresh header to reload while we don't have the metadata.
	panic("NOT IMPLEMENTED")
}

// /INFOHASH.torrent
//
func (wc *T) handleTorrentFileRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	w.Header().Set("Content Type", "application/x-bittorrent")
	data, _ := bencoding.Encode(bencoding.Dict{
		"info": wc.btClient.Swarm(infoHash, wc.dhtClient.GetPeers(infoHash).ReadNewPeers()).Info(),
	})
	w.Write(data)
}

// /INFOHASH.json
// Subset of torrent metadata in prety JSON format. Excludes info other than names and sizes.
func (wc *T) handleTorrentJsonInfoRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	// s, err := json.Marshal(wc.serialize())
	// if err != nil {
	// 	logger.Fatalf("JSON encoding failed: %v", err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "text/html")
	// w.Write(s)

	panic("NOT IMPLEMENTED")
}

func (wc *T) ListenAndServe() (err error) {
	logger.Info("Serving web client of %v at %v.", wc.staticPath, wc.addr)

	http.Handle("/_s/", http.StripPrefix(
		"/_s/", http.FileServer(http.Dir(wc.staticPath))))

	http.HandleFunc("/", wc.handleRequest)

	return http.ListenAndServe(wc.addr, nil)
}

// func (wc *T) DELETE_handlePeerSearch(w http.ResponseWriter, r *http.Request) {
// 	hexInfohash := r.FormValue("infohash")
// 	infohash, err := bittorrent.BTIDFromHex(hexInfohash)
// 	if err != nil {
// 		panic("BTID shouldn't have been invalid!")
// 	}

// 	peerSearch := wc.dhtClient.GetPeers(infohash)

// 	wc.peerSearches = append(wc.peerSearches, peerSearch)
// }

// func (wc *T) serializePeerSearches() (serialized []interface{}) {
// 	serialized = make([]interface{}, 0)

// 	for _, search := range wc.peerSearches {
// 		serialized = append(serialized, map[string]interface{}{
// 			"infohash":          search.Infohash.String(),
// 			"searchDistance":    search.Infohash.BitDistance(wc.dhtClient.Id()),
// 			"peers":             search.PeersFound,
// 			"finished":          search.Finished(),
// 			"queriedNodes":      wc.serializeQueriedNodes(search.QueriedNodes, search.Infohash),
// 			"queriedNodesCount": len(search.QueriedNodes),
// 		})
// 	}
// 	return
// }

// func (wc *T) serializeQueriedNodes(queriedNodes map[string]*dht.RemoteNode, target bittorrent.BTID) (serialized map[string]interface{}) {
// 	serialized = make(map[string]interface{}, 0)

// 	for key, node := range queriedNodes {
// 		var sourceId string
// 		if node.Source != nil {
// 			sourceId = node.Source.Id.String()
// 		}

// 		serialized[key] = map[string]interface{}{
// 			"id":             node.Id.String(),
// 			"sourceId":       sourceId,
// 			"localDistance":  node.Id.BitDistance(wc.dhtClient.Id()),
// 			"targetDistance": node.Id.BitDistance(target),
// 		}
// 	}

// 	return
// }

// func (wc *T) serialize() (serialized map[string]interface{}) {
// 	return map[string]interface{}{
// 		"dht": map[string]interface{}{
// 			"peerSearches":   wc.serializePeerSearches(),
// 			"connectionInfo": wc.dhtClient.ConnectionInfo(),
// 		},
// 	}
// }
