package p2p

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	PeerTimeout       = 1
	PeerDataIntegrity = 2
)

var (
	hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
)

type blacklisting struct {
	when   time.Time
	due    time.Time
	reason int
}

type Peer struct {
	baseUrl          string
	apiUrl           string
	startLastRequest time.Time

	blacklisting *blacklisting

	mu sync.Mutex
}

func NewPeer(baseUrl string, internetProtocols []string) (*Peer, error) {
	apiUrl, err := peerToUrl(baseUrl)
	if err != nil {
		return nil, err
	}

	ips, err := net.LookupHost(apiUrl.Hostname())
	if err != nil {
		return nil, err
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
		return nil, errors.New("unknown internet protocol of remote peer " + apiUrl.Hostname())
	}

	for _, protocol := range internetProtocols {
		if peerProtocolOf[protocol] {
			return &Peer{
				baseUrl: baseUrl,
				apiUrl:  apiUrl.String(),
			}, nil
		}
	}

	return nil, errors.New("internet protocol of remote peer " + apiUrl.Hostname() + " not supported by config")
}

func (p *Peer) Block(reason int) {
	p.mu.Lock()
	if p.blacklisting != nil {
		p.mu.Unlock()
		return
	}

	now := time.Now()
	var due time.Time
	switch reason {
	case PeerTimeout:
		due = now.Add(10 * time.Minute)
	case PeerDataIntegrity:
		due = now.Add(1 * time.Hour)
	}

	p.blacklisting = &blacklisting{
		reason: reason,
		when:   now,
		due:    due,
	}
	p.mu.Unlock()
}

func (p *Peer) Unblock() {
	p.mu.Lock()
	p.blacklisting = nil
	p.mu.Unlock()
}

func (p *Peer) TryUnblock() (unblocked bool) {
	now := time.Now()
	p.mu.Lock()
	if p.blacklisting.due.After(now) {
		p.blacklisting = nil
		p.mu.Unlock()
		return true
	}
	p.mu.Unlock()
	return false
}

func peerToUrl(peer string) (*url.URL, error) {
	hasSchemaRegexp, _ := regexp.Compile("^[^:]+://")
	if !hasSchemaRegexp.MatchString(peer) {
		peer = "http://" + peer
	}

	if !hasPortRegexp.MatchString(peer) {
		peer = peer + ":8123"
	}
	u, err := url.Parse(peer + "/burst")

	return u, err
}
