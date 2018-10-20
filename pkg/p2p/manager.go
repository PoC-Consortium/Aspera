package p2p

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lukechampine/randmap/safe"
	"go.uber.org/zap"

	. "github.com/ac0v/aspera/pkg/log"
)

const (
	PeerTimeout                 = 1
	PeerDataIntegrity           = 2
	PeerDataIntegrityValidation = 3
)

var (
	hasPortRegexp   = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	hasSchemaRegexp = regexp.MustCompile("^[^:]+://")

	ErrUnsupportedInternetProtocol = errors.New("unsupported internet protocol")
)

type Manager interface {
	BlockPeer(peer string, reason int)
	RandomPeer() string
}

type manager struct {
	client Client

	initialPeers []string

	allPeers   map[string]struct{}
	allPeersMu sync.RWMutex

	blockedPeers   map[string]time.Time
	blockedPeersMu sync.Mutex

	scanForNewPeersInterval time.Duration

	blacklisted time.Time

	internetProtocols []string

	peerIterator   *iterator
	peerIteratorMu sync.Mutex
}

func NewManager(c Client, peers, internetProtocols []string, scanForNewPeersInterval time.Duration) Manager {
	m := &manager{
		initialPeers:            peers,
		client:                  c,
		scanForNewPeersInterval: scanForNewPeersInterval,
		allPeers:                make(map[string]struct{}),
		blockedPeers:            make(map[string]time.Time),
		internetProtocols:       internetProtocols,
		peerIterator:            &iterator{peers: []string{}},
	}

	m.initPeers()

	go m.jobs(c)

	return m
}

func (m *manager) jobs(c Client) {
	unblockTicker := time.NewTicker(5 * time.Minute)
	scanForPeersTicker := time.NewTicker(m.scanForNewPeersInterval)
	for {
		select {
		case <-unblockTicker.C:
			m.unblockPeers()
		case <-scanForPeersTicker.C:
			m.addPeersOf(m.RandomPeer())
		}
	}
}

func (m *manager) unblockPeers() {
	Log.Info("clearing peer blacklist")

	now := time.Now()
	m.blockedPeersMu.Lock()
	for p, due := range m.blockedPeers {
		if now.After(due) {
			delete(m.blockedPeers, p)
		}
	}
	m.blockedPeersMu.Unlock()
}

func (m *manager) initPeers() {
	for _, u := range m.initialPeers {
		peer, err := m.baseURLToAPIURL(u)
		if err != nil {
			continue
		}
		m.addPeersOf(peer)
	}
}

func (m *manager) addPeersOf(peer string) {
	getPeersMsg, err := m.client.GetPeersOf(peer)
	if err != nil {
		return
	}
	Log.Info("adding peers", zap.String("from", peer), zap.Int("count", len(getPeersMsg.Peers)))
	m.newPeers(append(getPeersMsg.Peers, peer))
}

func (m *manager) RandomPeer() string {
	peer := m.peerIterator.Next()
	// check if iterator is exhausted
	for peer == nil {
		if peer = m.resetPeerIterator(); peer == nil {
			peer = m.peerIterator.Next()
		} else {
			return *peer
		}
	}
	return *peer
}

func (m *manager) resetPeerIterator() *string {
	m.peerIteratorMu.Lock()
	// another go routine might have recreated the iterator, while we were waiting for the lock
	if peer := m.peerIterator.Next(); peer != nil {
		m.peerIteratorMu.Unlock()
		return peer
	}

	m.allPeersMu.Lock()
	lenAllPeers := len(m.allPeers)

	// if we don't have any peers at this point we need to reinitialize
	if lenAllPeers == 0 {
		m.allPeersMu.Unlock()
		m.initPeers()
		m.allPeersMu.Lock()
		lenAllPeers = len(m.allPeers)
		// TODO: if we still have 0 peers we really have a problem...
	}

	peers := make([]string, lenAllPeers)

	var ignore struct{}
	var peer string
	peerI := randmap.Iter(m.allPeers, &peer, &ignore)
	i := 0
	for peerI.Next() {
		peers[i] = peer
		i++
	}
	m.allPeersMu.Unlock()

	m.peerIterator = &iterator{peers: peers, idx: 0}
	m.peerIteratorMu.Unlock()

	return nil
}

func (m *manager) BlockPeer(peer string, reason int) {
	Log.Warn("blocking peer", zap.String("peer", peer), zap.Int("reason", reason))

	now := time.Now()
	var due time.Time
	switch reason {
	case PeerTimeout:
		due = now.Add(10 * time.Minute)
	case PeerDataIntegrity:
		due = now.Add(1 * time.Hour)
	case PeerDataIntegrityValidation:
		due = now.Add(1 * time.Hour)
	}

	m.blockedPeersMu.Lock()
	m.blockedPeers[peer] = due
	m.blockedPeersMu.Unlock()
}

func (m *manager) newPeers(baseURLs []string) {
	newPeers := make(map[string]struct{})
	for _, baseURL := range baseURLs {
		peer, err := m.baseURLToAPIURL(baseURL)
		if err != nil {
			Log.Warn("failed to create api url", zap.Error(err))
			continue
		}

		newPeers[peer] = struct{}{}
	}

	m.blockedPeersMu.Lock()
	for peer := range newPeers {
		if _, blocked := m.blockedPeers[peer]; blocked {
			delete(newPeers, peer)
		}
	}
	m.blockedPeersMu.Unlock()

	if len(newPeers) == 0 {
		return
	}

	m.allPeersMu.Lock()
	for peer := range newPeers {
		m.allPeers[peer] = struct{}{}
	}
	m.allPeersMu.Unlock()
}

func (m *manager) baseURLToAPIURL(u string) (string, error) {
	if !hasSchemaRegexp.MatchString(u) {
		u = "http://" + u
	}
	if !hasPortRegexp.MatchString(u) {
		u = u + ":8123"
	}
	apiURL, err := url.Parse(u + "/burst")

	err = m.checkProtocol(apiURL)
	return apiURL.String(), err
}

func (m *manager) checkProtocol(u *url.URL) error {
	host := u.Hostname()
	ips, err := net.LookupHost(host)
	if err != nil {
		return err
	}

	peerProtocolOf := map[string]bool{}
	for _, ip := range ips {
		if strings.Count(ip, ":") < 2 {
			peerProtocolOf["v4"] = true
		} else if strings.Count(ip, ":") >= 2 {
			peerProtocolOf["v6"] = true
		}
	}
	if len(peerProtocolOf) == 0 {
		return ErrUnsupportedInternetProtocol
	}

	for _, protocol := range m.internetProtocols {
		if peerProtocolOf[protocol] {
			return nil
		}
	}
	return ErrUnsupportedInternetProtocol
}

type iterator struct {
	peers []string
	idx   int32
}

func (i *iterator) Next() *string {
	idx := atomic.LoadInt32(&i.idx)
	if int(idx) > len(i.peers)-1 {
		return nil
	}
	atomic.AddInt32(&i.idx, 1)
	p := i.peers[idx]
	return &p
}
