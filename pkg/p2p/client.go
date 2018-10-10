package p2p

import (
	"bytes"
	"errors"
	"gopkg.in/resty.v1"
	"net/http"
	"strconv"
	"time"
	

	b "github.com/ac0v/aspera/pkg/block"
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	"github.com/ac0v/aspera/pkg/config"
	"github.com/golang/protobuf/jsonpb"
	"github.com/json-iterator/go"
)

const (
	majority    = 3
	parallelism = 5
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type GetNextBlocksResponse struct {
	NextBlocks []*b.Block `json:"nextBlocks"`
}

type Client interface {
	GetNextBlockIDs(blockID uint64) (*pb.GetNextBlockIdsResponse, []*Peer, error)
	GetNextBlocks(blockID uint64) (*GetNextBlocksResponse, []*Peer, error)
	GetPeersOf(apiUrl string) (*pb.GetPeers, error)
}

type client struct {
	manager     Manager
	unmarshaler *jsonpb.Unmarshaler
}

func NewClient(config *config.P2P, internetProtocols []string) Client {
	resty.SetTimeout(config.Timeout)
	resty.SetDebug(config.Debug)

	c := &client{
		unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true},
	}
	pm := NewManager(c, config.Peers, internetProtocols, time.Minute)

	c.manager = pm

	return c
}

func (c *client) request(req *resty.Request) (*resty.Response, *Peer, error) {
	peer := c.manager.RandomPeer()
	res, err := req.Post(peer.apiUrl)
	return res, peer, err
}

func (c *client) requestByMajority(requestType string, params map[string]interface{}) ([]byte, []*Peer, error) {
	foundMajority := make(chan struct{})

	sem := make(chan struct{}, parallelism)
	for i := 0; i < parallelism; i++ {
		sem <- struct{}{}
	}

	type peerResponse struct {
		of   *Peer
		body string
	}
	peerResponses := make(chan *peerResponse)

	go func() {
		for {
			select {
			case <-foundMajority:
				return
			case <-sem:
				go func() {
					peer := c.manager.RandomPeer()
					res, err := c.buildRequest(requestType, params).Post(peer.apiUrl)

					if err != nil || res == nil || res.StatusCode() != http.StatusOK {
						c.manager.BlockPeer(peer, PeerTimeout)
					} else {
						peerResponses <- &peerResponse{
							of:   peer,
							body: res.String(),
						}
					}

					sem <- struct{}{}
				}()
			}
		}
	}()

	seenBy := make(map[string]map[*Peer]struct{})
	for peerResponse := range peerResponses {
		peers := seenBy[peerResponse.body]
		if peers == nil {
			seenBy[peerResponse.body] = map[*Peer]struct{}{peerResponse.of: struct{}{}}
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
							c.manager.BlockPeer(p, PeerDataIntegrity)
						}
					}
				}
				var peersSlice []*Peer
				for p := range peers {
					peersSlice = append(peersSlice, p)
				}
				return []byte(peerResponse.body), peersSlice, nil
			}
		}

	}

	return nil, nil, errors.New("unexpected error")
}

func (c *client) buildRequest(requestType string, params map[string]interface{}) *resty.Request {
	paramsCopy := make(map[string]interface{}, len(params))
	for k, v := range params {
		paramsCopy[k] = v
	}

	paramsCopy["protocol"] = "B1"
	paramsCopy["requestType"] = requestType

	return resty.R().SetBody(paramsCopy)
}

func (c *client) GetNextBlockIDs(blockId uint64) (*pb.GetNextBlockIdsResponse, []*Peer, error) {
	body, peers, err := c.requestByMajority("getNextBlockIds", map[string]interface{}{
		"blockId": strconv.FormatUint(blockId, 10),
	})
	if err != nil {
		return nil, nil, err
	}

	var msg = new(pb.GetNextBlockIdsResponse)
	err = c.unmarshaler.Unmarshal(bytes.NewReader(body), msg)
	return msg, peers, err
}

func (c *client) GetNextBlocks(blockId uint64) (*GetNextBlocksResponse, []*Peer, error) {
	req := c.buildRequest("getNextBlocks", map[string]interface{}{
		"blockId": strconv.FormatUint(blockId, 10),
	})
	res, peers, err := c.request(req)
	if err != nil {
		return nil, nil, err
	}

	var msg = new(GetNextBlocksResponse)
	return msg, []*Peer{peers}, json.Unmarshal(res.Body(), msg)
}

func (c *client) GetPeersOf(apiUrl string) (*pb.GetPeers, error) {
	res, err := c.buildRequest("getPeers", map[string]interface{}{}).Post(apiUrl)
	if err != nil {
		return nil, err
	}

	msg := new(pb.GetPeers)
	return msg, c.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), msg)
}
