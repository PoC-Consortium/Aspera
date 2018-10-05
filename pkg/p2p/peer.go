package p2p

import (
	"errors"
	r "github.com/ac0v/aspera/pkg/registry"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")

type Peer struct {
	baseUrl          string
	apiUrl           string
	startLastRequest time.Time

	rtt              time.Duration
	answeredRequests int
	mu               sync.Mutex
}

func NewPeer(registry *r.Registry, baseUrl string) (*Peer, error) {
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

	for _, protocol := range registry.Config.Network.InternetProtocols {
		if peerProtocolOf[protocol] {
			return &Peer{
				baseUrl: baseUrl,
				apiUrl:  apiUrl.String(),
			}, nil
		}
	}

	return nil, errors.New("internet protocol of remote peer " + apiUrl.Hostname() + " not supported by config")
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
