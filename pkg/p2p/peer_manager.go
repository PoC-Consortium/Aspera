package p2p

import (
	"bytes"
	"sync"
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	"github.com/lukechampine/randmap/safe"
)

var (
	syncPeerCount = 10
)

type PeerManager interface {
	BlockPeer(p *Peer)
	RandomPeer() *Peer
}

type peerManager struct {
	allPeers   map[string]*Peer
	allPeersMu sync.RWMutex

	blockedPeers   map[string]*Peer
	blockedPeersMu sync.Mutex

	scanForNewPeersInterval time.Duration

	blacklisted time.Time
}

func NewPeerManager(client *Client, peersBaseURLs []string, scanForNewPeersInterval time.Duration) PeerManager {
	pm := &peerManager{
		scanForNewPeersInterval: scanForNewPeersInterval,
		allPeers:                make(map[string]*Peer),
		blockedPeers:            make(map[string]*Peer),
	}

	pm.initPeers(client, peersBaseURLs)

	go pm.scanForNewPeers(client)

	return pm
}

func (pm *peerManager) initPeers(client *Client, peerBaseURLs []string) {
	for _, url := range peerBaseURLs {
		peer, err := NewPeer(url)
		if err != nil {
			continue
		}
		pm.addPeersOf(client, peer)
	}
}

func (pm *peerManager) addPeersOf(client *Client, peer *Peer) {
	getPeersMsg := new(pb.GetPeers)
	res, err := client.buildRequest("getPeers").Post(peer.apiURL)
	if err != nil {
		return
	}

	client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), getPeersMsg)

	pm.allPeersMu.Lock()
	pm.allPeers[peer.baseURL] = peer
	for _, baseURL := range getPeersMsg.Peers {
		p, err := NewPeer(baseURL)
		if err != nil {
			continue
		}
		pm.allPeers[baseURL] = p
	}
	pm.allPeersMu.Unlock()
}

func (pm *peerManager) RandomPeer() *Peer {
	pm.allPeersMu.RLock()
	p := randmap.Val(pm.allPeers).(*Peer)
	pm.allPeersMu.RUnlock()
	return p
}

func (pm *peerManager) BlockPeer(pToBlock *Peer) {
	pm.blockedPeersMu.Lock()
	_, blocked := pm.blockedPeers[pToBlock.baseURL]
	if blocked {
		pm.blockedPeersMu.Unlock()
		return
	}
	pm.blockedPeers[pToBlock.baseURL] = pToBlock
	pm.blockedPeersMu.Unlock()

	pm.allPeersMu.Lock()
	delete(pm.allPeers, pToBlock.baseURL)
	pm.allPeersMu.Unlock()
}

func (pm *peerManager) scanForNewPeers(client *Client) {
	for range time.NewTicker(pm.scanForNewPeersInterval).C {
		pm.addPeersOf(client, pm.RandomPeer())
	}
}
