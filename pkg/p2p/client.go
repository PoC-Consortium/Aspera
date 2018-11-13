package p2p

import (
	"bytes"
	"errors"
	"gopkg.in/resty.v1"
	"sync"

	api "github.com/ac0v/aspera/pkg/api/p2p"
	compat "github.com/ac0v/aspera/pkg/api/p2p/compat"
	b "github.com/ac0v/aspera/pkg/block"
	"github.com/ac0v/aspera/pkg/config"
	"github.com/ac0v/aspera/pkg/p2p/manager"

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
	GetNextBlockIDs(blockID uint64, height int32) (*api.GetNextBlockIdsResponse, []manager.Peer, error)
	GetNextBlocks(blockID uint64, height int32) (*api.GetNextBlocksResponse, []manager.Peer, error)
}

type client struct {
	manager     manager.Manager
	unmarshaler *jsonpb.Unmarshaler
}

func NewClient(config *config.P2P, manager manager.Manager) Client {
	resty.SetTimeout(config.Timeout)
	resty.SetDebug(config.Debug)

	return &client{
		manager:     manager,
		unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true},
	}
}

func (c *client) requestByMajority(height int32, req func(p manager.Peer) ([]byte, error)) ([]byte, []manager.Peer, error) {
	foundMajority := make(chan struct{})

	type peerResponse struct {
		of   manager.Peer
		body string
	}

	peerResponses := make(chan *peerResponse, majority)
	sem := make(chan struct{}, majority)
	var wg sync.WaitGroup
	go func() {
		for {
			select {
			case <-foundMajority:
				return
			case sem <- struct{}{}:
				wg.Add(1)
				go func() {
					peer := c.manager.RandomPeer(height)

					body, err := req(peer)
					if err == nil {
						peerResponses <- &peerResponse{
							of:   peer,
							body: string(body),
						}
					} else {
						peer.Throttle()
					}

					<-sem
					wg.Done()
				}()
			}
		}
	}()

	seenBy := make(map[string]map[manager.Peer]struct{})
	for peerResponse := range peerResponses {
		peers := seenBy[peerResponse.body]
		if peers == nil {
			seenBy[peerResponse.body] = map[manager.Peer]struct{}{peerResponse.of: struct{}{}}
		} else {
			if _, processedPeer := peers[peerResponse.of]; processedPeer {
				continue
			}

			peers[peerResponse.of] = struct{}{}

			if len(peers) >= majority {
				foundMajority <- struct{}{}
				for otherBody, peers := range seenBy {
					if otherBody != peerResponse.body {
						for p := range peers {
							p.Throttle()
						}
					}
				}
				var peersSlice []manager.Peer
				for p := range peers {
					p.DeThrottle()
					peersSlice = append(peersSlice, p)
				}
				wg.Wait()
				return []byte(peerResponse.body), peersSlice, nil
			}
		}

	}

	return nil, nil, errors.New("unexpected error")
}

func (c *client) GetNextBlockIDs(blockId uint64, height int32) (*api.GetNextBlockIdsResponse, []manager.Peer, error) {
	req := func(p manager.Peer) ([]byte, error) { return p.GetNextBlockIDsBody(blockId) }
	body, peers, err := c.requestByMajority(height, req)
	if err != nil {
		return nil, nil, err
	}

	var msg = new(api.GetNextBlockIdsResponse)
	err = c.unmarshaler.Unmarshal(bytes.NewReader(body), msg)
	return msg, peers, err
}

func (c *client) GetNextBlocks(blockId uint64, height int32) (*api.GetNextBlocksResponse, []manager.Peer, error) {
	p := c.manager.RandomPeer(height)
	body, err := p.GetNextBlocksBody(blockId)
	if err != nil {
		return nil, []manager.Peer{p}, err
	}

	var json []byte
	if json, err = compat.Upgrade(body); err != nil {
		return nil, nil, err
	}

	var msg = new(api.GetNextBlocksResponse)
	return msg, []manager.Peer{p}, c.unmarshaler.Unmarshal(bytes.NewReader(json), msg)
}
