package p2p

import (
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/jsonpb"
	//	"go.uber.org/zap"
	"bytes"
	"gopkg.in/resty.v1"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	registry     *r.Registry
	peerIterator *peerIterator
	unmarshaler  *jsonpb.Unmarshaler
}

func NewClient(registry *r.Registry) *Client {
	// client := &Client{peerIterator: NewPeerIterator(registry.Config.Peers), registry: registry, unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: false}}
	client := &Client{peerIterator: NewPeerIterator(registry.Config.Peers), registry: registry, unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true}}

	for range registry.Config.Peers {
		var s = new(pb.GetPeers)

		res := client.request(client.peerIterator.Next(), "getPeers")
		client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), s)

		client.peerIterator.Add(s.Peers)
	}
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

func (client *Client) autoRequest(byMajority bool, params ...map[string]interface{}) *resty.Response {
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

	if byMajority {
		seenCount := map[string]int8{}
		seenData := map[string]*resty.Response{}

		for {
			res := client.request(client.peerIterator.Next(), requestType, params...)
			seenCount[string(res.Body()[:])]++
			seenData[string(res.Body()[:])] = res

			var major string
			majorCount := int8(0)

			for value, count := range seenCount {
				if count > majorCount && count > 2 {
					major = value
				}
			}
			if len(major) > 0 {
				return seenData[major]
			}
		}
	}

	return client.request(client.peerIterator.Next(), requestType, params...)
}

func (client *Client) request(peer string, requestType string, params ...map[string]interface{}) *resty.Response {
	param := map[string]interface{}{
		"protocol":    "B1",
		"requestType": requestType,
	}
	for _, m := range params {
		param = mergeMaps(param, m)
	}

	//resty.SetDebug(true)
	//resty.SetDebugBodyLimit(1)
	resty.SetTimeout(1 * time.Second)

	//client.registry.Logger.Info("requesting", zap.String("peer", peer))
	res, _ := resty.R().SetBody(param).Post(peer)

	return res
}

type Transaction struct {
	Type    byte `json:"type"`
	Subtype byte `json:"subtype"`
}

// ToDo: transactions: theType byte, subtype byte, timestamp int, deadline uint16, senderPublicKey hex, amountNQT uint64, feeNQT uint64, referencedTransactionFullHash string, signature string, version byte, attachment object, recipient string, ecBlockHeight int, ecBlockId long) + attachment

func (client *Client) AddPeers(peers ...string) {
	client.autoRequest(false, map[string]interface{}{"peers": peers})
}
func (client *Client) GetCumulativeDifficulty() (*pb.GetCumulativeDifficultyResponse, error) {
	res := client.autoRequest(false)
	var s = new(pb.GetCumulativeDifficultyResponse)
	err := client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), s)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v\n", s)
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
	res := client.autoRequest(true, map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})

	var s = new(pb.GetNextBlockIdsResponse)
	err := client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), s)

	return s, err
}

func (client *Client) GetNextBlocks(blockId uint64) (*pb.GetNextBlocksResponse, error) {
	res := client.autoRequest(false, map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})

	var s = new(pb.GetNextBlocksResponse)
	err := client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), s)

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
