package p2p

import (
	"bytes"
	"errors"
	"gopkg.in/resty.v1"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/jsonpb"
)

const majority = 3

type Client struct {
	registry    *r.Registry
	peerManager PeerManager
	unmarshaler *jsonpb.Unmarshaler
}

func NewClient(registry *r.Registry) *Client {
	// TODO: timeout should be config option
	resty.SetTimeout(2 * time.Second)
	// resty.SetDebug(true)

	client := &Client{
		registry:    registry,
		unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true},
	}
	pm := NewPeerManager(client, registry, time.Minute)

	client.peerManager = pm

	return client
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func (client *Client) autoRequest(byMajority bool, params ...map[string]interface{}) ([]byte, error) {
	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(2, fpcs)
	if n == 0 {
		log.Fatal("can't figure out who called me")
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		log.Fatal("can't figure out who called me")
	}
	requestType := fun.Name()
	requestType = requestType[strings.LastIndex(requestType, ".")+1 : len(requestType)]
	requestType = strings.Replace(requestType, requestType[0:1], strings.ToLower(requestType[0:1]), 1)

	if !byMajority {
		req := client.buildRequest(requestType, params...)
		res, err := req.Post(client.peerManager.RandomPeer().apiURL)
		if err != nil {
			return nil, err
		}
		return res.Body(), nil
	}

	stop := make(chan struct{})

	// sem will ensure that there are only majority + n parallel requests
	// 2 is arbitrary
	sem := make(chan struct{}, majority+2)
	for i := 0; i < majority; i++ {
		sem <- struct{}{}
	}

	type PeerResponse struct {
		of         *Peer
		err        error
		statusCode int
		body       string
	}
	peerResponses := make(chan *PeerResponse)

	go func() {
		for {
			select {
			case <-stop:
				return
			case <-sem:
				go func() {
					req := client.buildRequest(requestType, params...)
					peer := client.peerManager.RandomPeer()

					peer.StartRequest()
					res, err := req.Post(peer.apiURL)
					peer.FinishRequest()

					peerResponses <- &PeerResponse{
						of:         peer,
						body:       string(res.Body()),
						err:        err,
						statusCode: res.StatusCode(),
					}

					sem <- struct{}{}
				}()
			}
		}
	}()

	type seen struct {
		count int
		peers map[*Peer]struct{}
	}

	seens := make(map[string]*seen)
	for peerResponse := range peerResponses {
		if peerResponse.err != nil || peerResponse.statusCode != http.StatusOK {
			client.peerManager.BlockPeer(peerResponse.of)
			continue
		}

		if _, knownResponse := seens[peerResponse.body]; !knownResponse {
			seens[peerResponse.body] = &seen{
				count: 1,
				peers: map[*Peer]struct{}{peerResponse.of: struct{}{}},
			}
		} else {
			seen := seens[peerResponse.body]
			if _, processedPeer := seen.peers[peerResponse.of]; processedPeer {
				continue
			}

			seen.count++
			seen.peers[peerResponse.of] = struct{}{}

			if seen.count >= majority {
				stop <- struct{}{}
				for otherBody, seen := range seens {
					if otherBody != peerResponse.body {
						for p := range seen.peers {
							client.peerManager.BlockPeer(p)
						}
					}
				}
				return []byte(peerResponse.body), nil
			}
		}

	}

	return nil, errors.New("unexpected error")
}

func (client *Client) buildRequest(requestType string, params ...map[string]interface{}) *resty.Request {
	param := map[string]interface{}{
		"protocol":    "B1",
		"requestType": requestType,
	}

	for _, m := range params {
		param = mergeMaps(param, m)
	}

	return resty.R().SetBody(param)
}

type Transaction struct {
	Type    byte `json:"type"`
	Subtype byte `json:"subtype"`
}

// ToDo: transactions: theType byte, subtype byte, timestamp int, deadline uint16, senderPublicKey hex, amountNQT uint64, feeNQT uint64, referencedTransactionFullHash string, signature string, version byte, attachment object, recipient string, ecBlockHeight int, ecBlockId long) + attachment

func (client *Client) GetCumulativeDifficulty() (*pb.GetCumulativeDifficultyResponse, error) {
	body, err := client.autoRequest(false)
	if err != nil {
		return nil, err
	}

	var s = new(pb.GetCumulativeDifficultyResponse)
	err = client.unmarshaler.Unmarshal(bytes.NewReader(body), s)
	if err != nil {
		log.Fatal(err)
	}

	return s, err

}
func (client *Client) GetInfo(announcedAddress string, application string, version string, platform string, shareAddress string) {
	client.autoRequest(
		false,
		map[string]interface{}{
			"announcedAddress": announcedAddress,
			"application":      application,
			"version":          version,
			"platform":         platform,
			"shareAddress":     shareAddress,
		},
	)
}
func (client *Client) GetMilestoneBlockIds(lastBlockId uint64, lastMilestoneBlockId uint64) {
	client.autoRequest(false, map[string]interface{}{
		"lastBlockId":          strconv.FormatUint(lastBlockId, 10),
		"lastMilestoneBlockId": strconv.FormatUint(lastMilestoneBlockId, 10),
	})
}

func (client *Client) GetNextBlockIds(blockId uint64) (*pb.GetNextBlockIdsResponse, error) {
	body, err := client.autoRequest(true, map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})
	if err != nil {
		return nil, err
	}

	var s = new(pb.GetNextBlockIdsResponse)
	err = client.unmarshaler.Unmarshal(bytes.NewReader(body), s)

	return s, err
}

func (client *Client) GetNextBlocks(blockId uint64) (*pb.GetNextBlocksResponse, error) {
	body, err := client.autoRequest(false, map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})
	if err != nil {
		return nil, err
	}

	var s = new(pb.GetNextBlocksResponse)
	err = client.unmarshaler.Unmarshal(bytes.NewReader(body), s)

	return s, err
}

func (client *Client) GetPeers()                   { client.autoRequest(false) }
func (client *Client) GetUnconfirmedTransactions() { client.autoRequest(false) }
func (client *Client) ProcessBlock()               { client.autoRequest(false) }
func (client *Client) ProcessTransactions(transactions ...*Transaction) {
	client.autoRequest(false, structs.Map(transactions[0]))
}
func (client *Client) GetAccountBalance(accountId uint64) {
	client.autoRequest(false, map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}
func (client *Client) GetAccountRecentTransactions(accountId uint64) {
	client.autoRequest(false, map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}
