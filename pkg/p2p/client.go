package p2p

import (
	"encoding/json"
	//"fmt"
	pb "github.com/ac0v/brsx/internal/api/protobuf-spec"
	"github.com/fatih/structs"
	"gopkg.in/resty.v1"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const url string = "http://wallet.burst.cryptoguru.org:8123/burst"

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func request(params ...map[string]interface{}) *resty.Response {
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

	res, _ := resty.R().SetBody(param).Post(url)

	return res
}

type Transaction struct {
	Type    byte `json:"type"`
	Subtype byte `json:"subtype"`
}

// ToDo: transactions: theType byte, subtype byte, timestamp int, deadline uint16, senderPublicKey hex, amountNQT uint64, feeNQT uint64, referencedTransactionFullHash string, signature string, version byte, attachment object, recipient string, ecBlockHeight int, ecBlockId long) + attachment

func AddPeers(peers ...string) {
	request(map[string]interface{}{"peers": peers})
}
func GetCumulativeDifficulty() (*pb.GetCumulativeDifficultyResponse, error) {
	res := request()
	var s = new(pb.GetCumulativeDifficultyResponse)
	err := json.Unmarshal(res.Body(), &s)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v\n", s)
	return s, err

}
func GetInfo(announcedAddress string, application string, version string, platform string, shareAddress string) {
	request(
		map[string]interface{}{
			"announcedAddress": announcedAddress,
			"application":      application,
			"version":          version,
			"platform":         platform,
			"shareAddress":     shareAddress,
		},
	)
}
func GetMilestoneBlockIds(lastBlockId uint64, lastMilestoneBlockId uint64) {
	request(map[string]interface{}{
		"lastBlockId":          strconv.FormatUint(lastBlockId, 10),
		"lastMilestoneBlockId": strconv.FormatUint(lastMilestoneBlockId, 10),
	})
}
func GetNextBlockIds(blockId uint64) {
	request(map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})
}
func GetNextBlocks(blockId uint64) (*pb.GetNextBlocksResponse, error) {
	res := request(map[string]interface{}{"blockId": strconv.FormatUint(blockId, 10)})

	var s = new(pb.GetNextBlocksResponse)
	err := json.Unmarshal(res.Body(), &s)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v\n", s)
	return s, err
}
func GetPeers()                   { request() }
func GetUnconfirmedTransactions() { request() }
func ProcessBlock()               { request() } // ToDo
func ProcessTransactions(transactions ...*Transaction) {
	request(structs.Map(transactions[0]))
}
func GetAccountBalance(accountId uint64) {
	request(map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}
func GetAccountRecentTransactions(accountId uint64) {
	request(map[string]interface{}{"account": strconv.FormatUint(accountId, 10)})
}
