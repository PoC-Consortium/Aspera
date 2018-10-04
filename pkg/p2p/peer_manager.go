package p2p

import (
	"bytes"
	"net/url"
	"regexp"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	"github.com/lukechampine/randmap"
)

var hasPortRegexp = regexp.MustCompile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")

type PeerManager interface {
	RandomPeerURL() string
	InitPeers(client *Client, initialPeers []string)
}

type peerManager struct {
	peerURLs map[string]struct{}
}

func NewPeerManager() PeerManager {
	return &peerManager{
		peerURLs: make(map[string]struct{}),
	}
}

func (pm *peerManager) InitPeers(client *Client, initialPeers []string) {
	peerURLs := make(map[string]struct{})
	for _, peer := range initialPeers {
		u, err := peerToURL(peer)
		if err != nil {
			continue
		}
		peerURLs[u] = struct{}{}

		var s = new(pb.GetPeers)

		if res, err := client.buildRequest("getPeers").Post(u); err == nil {
			client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), s)
			for _, newPeer := range s.Peers {
				u, err := peerToURL(newPeer)
				if err != nil {
					continue
				}
				peerURLs[u] = struct{}{}
			}
		}
	}
	pm.peerURLs = peerURLs
}

func (pm *peerManager) RandomPeerURL() string {
	return randmap.Key(pm.peerURLs).(string)
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
