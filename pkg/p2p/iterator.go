package p2p

import (
	"math/rand"
	"net/url"
	"regexp"
)

type peerIterator struct {
	current int
	peers   []string
}

func (it *peerIterator) Next() string {
	// shuffle peers to ask them in a different order on each iteration
	if it.current == 0 {
		for i := range it.peers {
			j := rand.Intn(i + 1)
			it.peers[i], it.peers[j] = it.peers[j], it.peers[i]
		}
	}

	it.current++
	if it.current >= len(it.peers) {
		it.current = 0
	}

	return it.peers[it.current]
}

func (it *peerIterator) Add(newPeers []string) {
	var validNewPeers []string
	for _, peer := range newPeers {
		peer, err := peerToUrl(peer)
		if err == nil {
			validNewPeers = append(validNewPeers, peer)
		}
	}

	seen := make(map[string]bool, len(it.peers)+len(validNewPeers))

	for _, peer := range it.peers {
		seen[peer] = true
	}

	for _, peer := range validNewPeers {
		if !seen[peer] {
			it.peers = append(it.peers, peer)
		}
		seen[peer] = true
	}

}

func NewPeerIterator(peers []string) *peerIterator {
	var validPeers []string
	for _, peer := range peers {
		peer, err := peerToUrl(peer)
		if err == nil {
			validPeers = append(validPeers, peer)
		}
	}

	return &peerIterator{peers: validPeers, current: -1}
}

func peerToUrl(peer string) (string, error) {
	hasSchemaRegexp, _ := regexp.Compile("^[^:]+://")
	if !hasSchemaRegexp.MatchString(peer) {
		peer = "http://" + peer
	}

	hasPortRegexp, _ := regexp.Compile(":([1-9]|[1-8][0-9]|9[0-9]|[1-8][0-9]{2}|9[0-8][0-9]|99[0-9]|[1-8][0-9]{3}|9[0-8][0-9]{2}|99[0-8][0-9]|999[0-9]|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$")
	if !hasPortRegexp.MatchString(peer) {
		peer = peer + ":8123"
	}
	url, err := url.Parse(peer + "/burst")

	return url.String(), err
}
