package store

import (
	"fmt"
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/dixonwille/skywalker"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type RawStore struct {
	BasePath string
	registry *r.Registry
	Current  *RawCurrent
}

type RawCurrent struct {
	Height int32
	Block  *pb.Block
}

type LookupWorker struct {
	*sync.Mutex
	found []string
}

func (lookupWorker *LookupWorker) Work(path string) {
	lookupWorker.Lock()
	defer lookupWorker.Unlock()
	lookupWorker.found = append(lookupWorker.found, path)
}

func NewRawStore(registry *r.Registry) *RawStore {
	var rawStore RawStore
	rawStore.registry = registry

	rawStore.BasePath = filepath.Join(registry.Config.Storage.Path, "raw")
	if _, err := os.Stat(rawStore.BasePath); os.IsNotExist(err) {
		os.MkdirAll(rawStore.BasePath, os.ModePerm)
	}

	lookupWorker := new(LookupWorker)
	lookupWorker.Mutex = new(sync.Mutex)

	sw := skywalker.New(rawStore.BasePath, lookupWorker)
	err := sw.Walk()
	if err != nil {
		rawStore.registry.Logger.Fatal("Fatal", zap.Error(err))
	}
	sort.Sort(sort.StringSlice(lookupWorker.found))

	// update current; create genesis or load most recent block
	height := int64(-1)
	for _, f := range lookupWorker.found {
		filePath := strings.Replace(f, sw.Root, "", 1)
		currentHeight, err := strconv.ParseInt(strings.Replace(filePath, string(os.PathSeparator), "", 10)[0:10], 10, 32)
		if err == nil && currentHeight == height+1 {
			height++
		} else {
			// looks like we found some non- raw storage stuff / out of order blocks which
			// could be the result of a interupted async blockchain sync
			rawStore.registry.Logger.Info("removing orphaned file from raw storage", zap.String("path", filePath))
			os.Remove(rawStore.BasePath + string(os.PathSeparator) + filePath)
		}
	}
	rawStore.Current = &RawCurrent{Height: int32(height)}

	if height == -1 {
		block := &pb.Block{Block: registry.Config.Network.P2P.Milestones[0].Id}
		rawStore.Store(block, registry.Config.Network.P2P.Milestones[0].Height)
	} else {
		rawStore.Current.Block = rawStore.load(rawStore.Current.Height)
	}
	rawStore.registry.Logger.Info("loaded Raw Storage", zap.Int("height", int(rawStore.Current.Height)))

	return &rawStore
}

func (rawStore *RawStore) Push(block *pb.Block) {
	rawStore.Store(block, rawStore.Current.Height+1)
}

func (rawStore *RawStore) convertHeightToPathInfo(height int32) string {
	parts := []rune(fmt.Sprintf("%010d", int(height)))
	var path string
	for _, part := range parts {
		path = filepath.Join(path, string(part))
	}
	return filepath.Join(rawStore.BasePath, path+".bin")
}

func (rawStore *RawStore) Store(block *pb.Block, height int32) {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	block.Height = height
	/* ToDo:
	block.PayloadLength
	block.TotalAmountNQT
	block.Generator
	block.BaseTarget
	block.generatorRS
	block.BlockReward
	block.NextBlock
	block.ScoopNum
	block.NumberOfTransactions
	block.Transactions
	block.TotalFeeNQT
	block.Block
	*/

	data, _ := proto.Marshal(block)
	ioutil.WriteFile(path, data, os.ModePerm)

	rawStore.Current.Height = height
	rawStore.Current.Block = block
}

func (rawStore *RawStore) load(height int32) *pb.Block {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	in, err := ioutil.ReadFile(path)
	if err != nil {
		rawStore.registry.Logger.Fatal("Error reading file:", zap.Error(err))
	}
	block := &pb.Block{}
	if err := proto.Unmarshal(in, block); err != nil {
		rawStore.registry.Logger.Fatal("Error parse block file:", zap.Error(err))
	}

	return block
}
