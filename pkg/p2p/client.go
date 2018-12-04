package p2p

import (
	"bytes"
	"errors"
	"gopkg.in/resty.v1"

	api "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	compat "github.com/PoC-Consortium/Aspera/pkg/api/p2p/compat"
	b "github.com/PoC-Consortium/Aspera/pkg/block"
	"github.com/PoC-Consortium/Aspera/pkg/config"
	. "github.com/PoC-Consortium/Aspera/pkg/p2p/manager"

	"github.com/golang/protobuf/jsonpb"
	"github.com/json-iterator/go"
)

const (
	majority = 3
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type GetNextBlocksResponse struct {
	NextBlocks []*b.Block `json:"nextBlocks"`
}

type Client interface {
	GetNextBlockIDs(blockID uint64, height int32) (*api.GetNextBlockIdsResponse, []Peer, error)
	GetNextBlocks(blockID uint64, height int32) (*api.GetNextBlocksResponse, []Peer, error)
}

type client struct {
	manager     Manager
	unmarshaler *jsonpb.Unmarshaler
}

func NewClient(config *config.P2P, manager Manager) Client {
	resty.SetTimeout(config.Timeout)
	resty.SetDebug(config.Debug)

	return &client{
		manager:     manager,
		unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true},
	}
}

type vote struct {
	voter *voter
	peer  Peer
	body  string
}

type voter struct {
	token chan struct{}
	quit  chan struct{}
}

func (v *voter) vote() {
	v.token <- struct{}{}
}

func (v *voter) stop() {
	v.quit <- struct{}{}
}

func newVoter(votes chan *vote, height int32, req func(p Peer) ([]byte, error), pm Manager) *voter {
	v := &voter{
		token: make(chan struct{}),
		quit:  make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-v.quit:
				return
			case <-v.token:
			retryVote:
				peer := pm.RandomPeer(height)
				if body, err := req(peer); err == nil {
					votes <- &vote{
						voter: v,
						peer:  peer,
						body:  string(body),
					}
				} else {
					peer.Throttle()
					if len(v.quit) == 0 {
						goto retryVote
					}
					return
				}
			}
		}
	}()
	return v
}

func (c *client) requestByMajority(height int32, req func(p Peer) ([]byte, error)) ([]byte, []Peer, error) {
	votes := make(chan *vote, majority)
	voters := make([]*voter, majority)
	for i := range voters {
		v := newVoter(votes, height, req, c.manager)
		v.vote()
		voters[i] = v
	}

	seenBy := make(map[string]map[Peer]struct{})
	for vote := range votes {
		peers := seenBy[vote.body]
		if peers == nil {
			seenBy[vote.body] = map[Peer]struct{}{vote.peer: struct{}{}}
		} else if _, processedPeer := peers[vote.peer]; !processedPeer {
			peers[vote.peer] = struct{}{}

			if len(peers) >= majority {
				for _, v := range voters {
					v.stop()
				}
				// close(votes)
				for otherBody, peers := range seenBy {
					if otherBody != vote.body {
						for p := range peers {
							p.Throttle()
						}
					}
				}
				var peersSlice []Peer
				for p := range peers {
					p.DeThrottle()
					peersSlice = append(peersSlice, p)
				}
				return []byte(vote.body), peersSlice, nil
			}
		}
		vote.voter.vote()
	}

	return nil, nil, errors.New("unexpected error")
}

func (c *client) GetNextBlockIDs(blockId uint64, height int32) (*api.GetNextBlockIdsResponse, []Peer, error) {
	req := func(p Peer) ([]byte, error) { return p.GetNextBlockIDsBody(blockId) }
	body, peers, err := c.requestByMajority(height, req)
	if err != nil {
		return nil, nil, err
	}

	var msg = new(api.GetNextBlockIdsResponse)
	err = c.unmarshaler.Unmarshal(bytes.NewReader(body), msg)
	return msg, peers, err
}

func (c *client) GetNextBlocks(blockId uint64, height int32) (*api.GetNextBlocksResponse, []Peer, error) {
	p := c.manager.RandomPeer(height)
	body, err := p.GetNextBlocksBody(blockId)
	if err != nil {
		return nil, []Peer{p}, err
	}

	var json []byte
	if json, err = compat.Upgrade(body); err != nil {
		return nil, nil, err
	}

	var msg = new(api.GetNextBlocksResponse)
	return msg, []Peer{p}, c.unmarshaler.Unmarshal(bytes.NewReader(json), msg)
}
