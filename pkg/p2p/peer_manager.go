package p2p

import (
	"bytes"
	"net/url"
	"regexp"
	"sync"
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	"github.com/lukechampine/randmap/safe"
)

var hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")

type PeerManager interface {
	RandomPeerURL() string
	InitPeers(client *Client, initialPeers []string)
}

type peerManager struct {
	peerURLs                map[string]struct{}
	peerURLsMu              sync.RWMutex
	scanForNewPeersInterval time.Duration
}

func NewPeerManager(scanForNewPeersInterval time.Duration) PeerManager {
	return &peerManager{
		scanForNewPeersInterval: scanForNewPeersInterval,
		peerURLs:                make(map[string]struct{}),
	}
}

func (pm *peerManager) scanForNewPeers(client *Client) {
	for range time.NewTicker(pm.scanForNewPeersInterval).C {
		pm.addPeersOf(client, pm.RandomPeerURL())
	}
}

func (pm *peerManager) addPeersOf(client *Client, peerURL string) {
	getPeersMsg := new(pb.GetPeers)
	if res, err := client.buildRequest("getPeers").Post(peerURL); err == nil {
		client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), getPeersMsg)
	}

	pm.peerURLsMu.Lock()
	pm.peerURLs[peerURL] = struct{}{}
	for _, peer := range getPeersMsg.Peers {
		u, err := peerToURL(peer)
		if err != nil {
			continue
		}

		pm.peerURLs[u] = struct{}{}
	}
	pm.peerURLsMu.Unlock()
}

func (pm *peerManager) InitPeers(client *Client, initialPeers []string) {
	for _, peer := range initialPeers {
		u, err := peerToURL(peer)
		if err != nil {
			continue
		}

		pm.addPeersOf(client, u)
	}
	go pm.scanForNewPeers(client)
}

func (pm *peerManager) RandomPeerURL() string {
	pm.peerURLsMu.RLock()
	u := randmap.Key(pm.peerURLs).(string)
	pm.peerURLsMu.RUnlock()
	return u
}

func peerToURL(peer string) (string, error) {
	hasSchemaRegexp, _ := regexp.Compile("^[^:]+://")
	if !hasSchemaRegexp.MatchString(peer) {
		peer = "http://" + peer
	}

	if !hasPortRegexp.MatchString(peer) {
		peer = peer + ":8123"
	}
	u, err := url.Parse(peer + "/burst")

	return u.String(), err
}
