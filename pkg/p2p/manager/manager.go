package p2p

import (
	"errors"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const maxPeers = 300

var (
	hasPortRegexp   = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	hasSchemaRegexp = regexp.MustCompile("^[^:]+://")

	ErrUnsupportedInternetProtocol = errors.New("unsupported internet protocol")
)

type Manager interface {
	NewIterator(height int32) Iterator
	RenewPeers()
}

type manager struct {
	peers   map[string]Peer
	peersMu sync.Mutex

	stopJobs   chan struct{}
	renewPeers chan struct{}

	internetProtocols []string
}

// NewManager returns a new object that manages peers
func NewManager(peerUrls, internetProtocols []string) Manager {
	m := &manager{
		internetProtocols: internetProtocols,
		peers:             make(map[string]Peer),
	}

	// set inital peers
	for _, baseUrl := range peerUrls {
		if apiUrl, err := m.baseURLToAPIURL(baseUrl); err == nil {
			p := NewPeer(apiUrl)
			p.SetHeight()
			m.peers[baseUrl] = p
		}
	}

	m.initPeers(maxPeers)
	go m.jobs()

	return m
}

// NewIterator returns a peer iterator that contains peers
// that already synced to a minimum Height of minHeight
func (m *manager) NewIterator(minHeight int32) Iterator {
	now := time.Now()
	m.peersMu.Lock()
	var peers []Peer
	for _, peer := range m.peers {
		if peer.IsUsable(minHeight, now) {
			peers = append(peers, peer)
		}
	}
	m.peersMu.Unlock()

	shufflePeers(peers)

	return NewIterator(peers)
}

// InitPeers initialises the peers new with the following rules:
// 1. remember throttled peers
// 2. keep an old peer with a chance of 50%, but only a max of maxPeers/2
// 3. try to fill up the remaining maxPeers/2 with random chunks from current peers
func (m *manager) initPeers(maxPeers int) {
	var unthrottledPeers int
	newPeers := make(map[string]Peer, maxPeers)
	now := time.Now()
	m.peersMu.Lock()

	// ask the current peers in a random order for new peers
	currentPeersRand := make([]Peer, len(m.peers))
	var i int
	for u, p := range m.peers {
		currentPeersRand[i] = p
		// we don't forget about peers that we throttled
		if p.IsThrottled(now) {
			newPeers[u] = p
		} else if unthrottledPeers < maxPeers/2 && rand.Int()&1 == 1 { // toss a coin
			newPeers[u] = p
			unthrottledPeers++
		}
		i++
	}
	m.peersMu.Unlock()

	shufflePeers(currentPeersRand)

	// new loop, because we don't want to make requests inside a lock
	sem := make(chan struct{}, 8)
	for _, peer := range currentPeersRand {
		peerUrls, err := peer.GetPeerUrls()
		if err != nil {
			continue
		}

		shuffleStrings(peerUrls)
		// TODO: only trust n peers from a single peer
		if len(peerUrls) > 40 {
			peerUrls = peerUrls[:40]
		}

		for _, baseUrl := range peerUrls {
			if _, exists := newPeers[baseUrl]; !exists {
				if apiUrl, err := m.baseURLToAPIURL(baseUrl); err == nil {
					p := NewPeer(apiUrl)
					// TODO: probably use a go routine pool
					sem <- struct{}{}
					go func() {
						p.SetHeight()
						<-sem
					}()
					newPeers[baseUrl] = p
					unthrottledPeers++
				}
			}
		}
		if len(newPeers) > maxPeers {
			break
		}
	}

	newPeersSlice := make([]Peer, len(newPeers))
	var j int
	for _, p := range newPeers {
		newPeersSlice[j] = p
		j++
	}

	m.peersMu.Lock()
	m.peers = newPeers
	m.peersMu.Unlock()
}

func shuffleStrings(xs []string) {
	for i := range xs {
		j := rand.Intn(i + 1)
		xs[i], xs[j] = xs[j], xs[i]
	}
}

func shufflePeers(peers []Peer) {
	for i := range peers {
		j := rand.Intn(i + 1)
		peers[i], peers[j] = peers[j], peers[i]
	}
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

// RenewPeers renews Peers if they weren't initialised recently
func (m *manager) RenewPeers() {
	m.renewPeers <- struct{}{}
}

func (m *manager) jobs() {
	renewPeersInterval := 10 * time.Minute
	renewPeersTicker := time.After(renewPeersInterval)

	lastRenew := time.Now()
	renewPeers := func() {
		// only accept renewing peers every 30 seconds
		if time.Now().Sub(lastRenew) < 30*time.Second {
			return
		}
		m.initPeers(maxPeers)
		renewPeersTicker = time.After(renewPeersInterval)
		lastRenew = time.Now()
	}
	for {
		select {
		case <-renewPeersTicker:
			renewPeers()
		case <-m.renewPeers:
			renewPeers()
		case <-m.stopJobs:
			return
		}
	}
}

type Iterator interface {
	Next() Peer
}

type iterator struct {
	peers []Peer
	idx   int
	sync.Mutex
}

func NewIterator(peers []Peer) Iterator {
	return &iterator{peers: peers}
}

func (i *iterator) Next() Peer {
	i.Lock()
	defer i.Unlock()
	if i.idx < len(i.peers) {
		p := i.peers[i.idx]
		i.idx++
		return p
	}
	return nil
}
