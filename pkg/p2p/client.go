package p2p

import (
	"encoding/json"
	//"fmt"
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/fatih/structs"
	"gopkg.in/resty.v1"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type peerIterator struct {
	current int
	peers   []string
}

func (it *peerIterator) NextPeer() string {
	it.current++
	if it.current >= len(it.peers) {
		it.current = 0
	}
	// shuffle peers to ask them in a different order on each iteration
	if it.current == 0 {
		for i := range it.peers {
			j := rand.Intn(i + 1)
			it.peers[i], it.peers[j] = it.peers[j], it.peers[i]
		}
	}
	return it.peers[it.current]
}

func NewPeerIterator(peers []string) *peerIterator {
	return &peerIterator{peers: peers, current: -1}
}

type Client struct {
	registry     *r.Registry
	peerIterator *peerIterator
}

func NewClient(registry *r.Registry) *Client {
	return &Client{peerIterator: NewPeerIterator(registry.Config.Peers), registry: registry}
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

func (client *Client) request(params ...map[string]interface{}) *resty.Response {
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

	//resty.SetDebug(true)

	param := map[string]interface{}{
		"protocol":    "B1",
		"requestType": requestType,
	}
	for _, m := range params {
		param = mergeMaps(param, m)
	}

	res, _ := resty.R().SetBody(param).Post(client.peerIterator.NextPeer())

	return res
}

type Transaction struct {
	Type    byte `json:"type"`
	Subtype byte `json:"subtype"`
}

// ToDo: transactions: theType byte, subtype byte, timestamp int, deadline uint16, senderPublicKey hex, amountNQT uint64, feeNQT uint64, referencedTransactionFullHash string, signature string, version byte, attachment object, recipient string, ecBlockHeight int, ecBlockId long) + attachment

func (client *Client) AddPeers(peers ...string) {
	client.request(map[string]interface{}{"peers": peers})
}
func (client *Client) GetCumulativeDifficulty() (*pb.GetCumulativeDifficultyResponse, error) {
	res := client.request()
	var s = new(pb.GetCumulativeDifficultyResponse)
	err := json.Unmarshal(res.Body(), &s)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v\n", s)
	return s, err

}
func (client *Client) GetInfo(announcedAddress string, application string, version string, platform string, shareAddress string) {
	client.request(
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
	client.request(map[string]interface{}{
		"lastBlockId":          strconv.FormatUint(lastBlockId, 10),
		"lastMilestoneBlockId": strconv.FormatUint(lastMilestoneBlockId, 10),
	})
}
func (client *Client) GetNextBlockIds(blockId uint64) {
	client.request(map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})
}

func (client *Client) GetNextBlocks(blockId uint64) (*pb.GetNextBlocksResponse, error) {
	res := client.request(map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})

	var s = new(pb.GetNextBlocksResponse)
	err := json.Unmarshal(res.Body(), &s)

	return s, err
}

func (client *Client) GetPeers()                   { client.request() }
func (client *Client) GetUnconfirmedTransactions() { client.request() }
func (client *Client) ProcessBlock()               { client.request() } // ToDo
func (client *Client) ProcessTransactions(transactions ...*Transaction) {
	client.request(structs.Map(transactions[0]))
}
func (client *Client) GetAccountBalance(accountId uint64) {
	client.request(map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}
func (client *Client) GetAccountRecentTransactions(accountId uint64) {
	client.request(map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}

func (client *Client) GetNextBlocksByMajority(blockId uint64) (*pb.GetNextBlocksResponse, error) {
	var responses []*pb.GetNextBlocksResponse
	for {
		r, err := client.GetNextBlocks(blockId)
		if err != nil {
			continue
		}
		responses = append(responses, r)

		// majority vote after having more than 4 responses
		if len(responses) > 4 {
			var major = responses[0]
			var count = 1
			for _, current := range responses[1:] {
				if reflect.DeepEqual(major, current) {
					count++
				} else {

					count--
					if count == 0 {
						major = current
					}
				}
			}
			if count > 4 {
				return major, err
			}
		}
	}
}
