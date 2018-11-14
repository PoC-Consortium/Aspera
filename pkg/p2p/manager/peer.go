package manager

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	api "github.com/ac0v/aspera/pkg/api/p2p"

	"github.com/json-iterator/go"
	"gopkg.in/resty.v1"
)

func init() {
	// TODO: use dynamic timeout
	resty.SetTimeout(5 * time.Second)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Peer interface {
	SetHeight() error
	GetHeight() int32
	GetNextBlockIDsBody(blockId uint64) ([]byte, error)
	GetNextBlocksBody(blockId uint64) ([]byte, error)
	GetPeerUrls() ([]string, error)
	Throttle()
	DeThrottle()
	IsThrottled(now time.Time) bool
	IsUsable(height int32, now time.Time) bool
}

type peer struct {
	apiUrl string

	height int32

	throttleWeight int
	throttledUntil time.Time

	sync.Mutex
}

func NewPeer(apiUrl string) Peer {
	return &peer{
		apiUrl: apiUrl,
	}
}

func (p *peer) SetHeight() error {
	body, err := p.postRequest("getCumulativeDifficulty", map[string]interface{}{})
	if err != nil {
		return err
	}

	var msg api.GetCumulativeDifficultyResponse
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}
	atomic.StoreInt32(&p.height, msg.BlockchainHeight)
	return nil
}

func (p *peer) GetHeight() int32 {
	return atomic.LoadInt32(&p.height)
}

func (p *peer) GetPeerUrls() ([]string, error) {
	body, err := p.postRequest("getPeers", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var msg api.GetPeers
	if err := json.Unmarshal(body, &msg); err == nil {
		return msg.Peers, nil
	} else {
		return nil, err
	}
}

func (p *peer) GetNextBlockIDsBody(blockId uint64) ([]byte, error) {
	return p.postRequest("getNextBlockIds", map[string]interface{}{
		"blockId": strconv.FormatUint(blockId, 10),
	})
}

func (p *peer) GetNextBlocksBody(blockId uint64) ([]byte, error) {
	return p.postRequest("getNextBlocks", map[string]interface{}{
		"blockId": strconv.FormatUint(blockId, 10),
	})
}

func (p *peer) DeThrottle() {
	p.Lock()
	p.throttleWeight = 0
	p.Unlock()
}

func (p *peer) Throttle() {
	p.Lock()
	p.throttleWeight++
	var throttleDuration time.Duration
	// TODO: use more reasonable values
	switch p.throttleWeight {
	default:
		throttleDuration = 20 * time.Second
		// case 1:
		// 	throttleDuration = 30 * time.Second
		// case 2:
		// 	throttleDuration = 2 * time.Minute
		// case 3:
		// 	throttleDuration = 5 * time.Minute
		// case 4:
		// 	throttleDuration = 10 * time.Minute
		// default:
		// 	throttleDuration = 30 * time.Minute
	}
	p.throttledUntil = time.Now().Add(throttleDuration)
	p.Unlock()
}

func (p *peer) IsThrottled(now time.Time) bool {
	p.Lock()
	defer p.Unlock()
	if now.Before(p.throttledUntil) {
		return false
	}
	return true
}

func (p *peer) IsUsable(minHeight int32, now time.Time) bool {
	p.Lock()
	defer p.Unlock()
	if p.height < minHeight && now.Before(p.throttledUntil) {
		return false
	}
	return true
}

func (p *peer) postRequest(requestType string, params map[string]interface{}) ([]byte, error) {
	paramsCopy := make(map[string]interface{}, len(params))
	for k, v := range params {
		paramsCopy[k] = v
	}

	paramsCopy["protocol"] = "B1"
	paramsCopy["requestType"] = requestType

	req := resty.R().SetBody(paramsCopy)
	if res, err := req.Post(p.apiUrl); err == nil {
		return res.Body(), err
	} else {
		return nil, err
	}
}
