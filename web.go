package webclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

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

	if len(r.URL.Path) < 21 {
		logger.Error("404 unexpected %v", r.URL.Path)
		return
	}

	btih, err := bittorrent.BTIDFromHex(r.URL.Path[1:41])
	if err != nil {
		logger.Error("got request for invalid BTID: %v", err)
		return
	}

	extension := string(r.URL.Path[41:])

	switch extension {
	case "":
		wc.handleTorrentPageRequest(w, r, btih)
	case ".torrent":
		wc.handleTorrentFileRequest(w, r, btih)
		return
	case ".json":
		wc.handleTorrentJsonInfoRequest(w, r, btih)
		return
	default:
		logger.Error("404 extension %v", extension)
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
	logger.Notice("Serving %v page request", infoHash)

	// Temporarily doing this client-side
	wc.handleIndex(w, r)
	return

	// TODO: Get rid of all JavaScript.
	// Use a Refresh header to reload while we don't have the metadata.
}

// /INFOHASH.torrent
//
func (wc *T) handleTorrentFileRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	logger.Notice("Serving %v.torrent file request", infoHash)

	info := wc.btClient.Swarm(infoHash, wc.dhtClient.GetPeers(infoHash).ReadNewPeers()).Info()

	data, err := bencoding.Encode(bencoding.Dict{
		"info":          info,
		"announce-list": bencoding.List{},
		"nodes":         bencoding.List{},
	})
	if err != nil {
		logger.Error("unable to encode torrent: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/x-bittorrent")
	// TODO: encode filename & filenmae8* properly
	w.Header().Set("Content-Disposition", "attachment;filename="+string(info["name"].(bencoding.String))+".bittorrent")
	w.Write(data)
}

// /INFOHASH.json
// Subset of torrent metadata in prety JSON format. Excludes info other than names and sizes.
func (wc *T) handleTorrentJsonInfoRequest(w http.ResponseWriter, r *http.Request, infoHash bittorrent.BTID) {
	logger.Notice("Serving %v.json request", infoHash)

	info := wc.btClient.Swarm(infoHash, wc.dhtClient.GetPeers(infoHash).ReadNewPeers()).Info()
	// Info() is mutable/shared -- don't edit it.

	s, err := json.Marshal(bencoding.Dict{
		"name": info["name"],
	})
	if err != nil {
		logger.Fatalf("JSON encoding failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.Write(s)

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
