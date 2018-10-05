package p2p

import (
	"net/url"
	"regexp"
	"sync"
	"time"
)

var hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")

type Peer struct {
	baseURL          string
	apiURL           string
	startLastRequest time.Time

	rtt              time.Duration
	answeredRequests int
	mu               sync.Mutex
}

func NewPeer(baseURL string) (*Peer, error) {
	apiURL, err := peerToURL(baseURL)
	if err != nil {
		return nil, err
	}
	return &Peer{
		baseURL: baseURL,
		apiURL:  apiURL,
	}, nil
}

func (p *Peer) StartRequest() {
	p.mu.Lock()
	p.startLastRequest = time.Now()
	p.mu.Unlock()
}

func (p *Peer) FinishRequest() {
	p.mu.Lock()
	p.answeredRequests++
	p.rtt = time.Now().Sub(p.startLastRequest)
	p.mu.Unlock()
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
