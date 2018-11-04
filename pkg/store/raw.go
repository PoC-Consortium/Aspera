package store

import (
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"

	"github.com/dgraph-io/badger"
	"github.com/dixonwille/skywalker"
	"github.com/golang/protobuf/proto"
)

type RawStore struct {
	BasePath string
	Current  *RawCurrent
	queue    *badger.DB
}

type RawCurrent struct {
	Height int32
	Block  *api.Block
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

func NewRawStore(path string, genesisMilestone config.Milestone) *RawStore {
	var rawStore RawStore

	rawStore.BasePath = filepath.Join(path, "raw")
	if _, err := os.Stat(rawStore.BasePath); os.IsNotExist(err) {
		os.MkdirAll(rawStore.BasePath, os.ModePerm)
	}

	lookupWorker := new(LookupWorker)
	lookupWorker.Mutex = new(sync.Mutex)

	sw := skywalker.New(rawStore.BasePath, lookupWorker)
	err := sw.Walk()
	if err != nil {
		Log.Fatal("Fatal", zap.Error(err))
	}
	sort.Sort(sort.StringSlice(lookupWorker.found))

	// initialize queue
	basePath := filepath.Join(path, "raw.queue")
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}

	opts := badger.DefaultOptions
	opts.Dir = basePath
	opts.ValueDir = basePath
	queue, err := badger.Open(opts)
	if err != nil {
		zap.Error(err)
	}
	rawStore.queue = queue

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
			Log.Info("removing orphaned file from raw storage", zap.String("path", filePath))
			os.Remove(rawStore.BasePath + string(os.PathSeparator) + filePath)
		}
	}
	rawStore.Current = &RawCurrent{Height: int32(height)}

	if height == -1 {
		payloadHash, _ := hex.DecodeString(`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`)
		genesisBlock := &api.Block{
			Id:                  genesisMilestone.Id,
			Version:             -1,
			GeneratorPublicKey:  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			BlockSignature:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			GenerationSignature: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			PayloadHash:         payloadHash,
		}
		rawStore.Store(genesisBlock, genesisMilestone.Height)
	} else {
		rawStore.Current.Block = rawStore.load(rawStore.Current.Height)
	}
	Log.Info("loaded Raw Storage", zap.Int("height", int(rawStore.Current.Height)))

	return &rawStore
}

func (rawStore *RawStore) Push(block *api.Block) {
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

func (rawStore *RawStore) Store(block *api.Block, height int32) {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	data, _ := proto.Marshal(block)
	ioutil.WriteFile(path, data, os.ModePerm)

	rawStore.Current.Height = height
	rawStore.Current.Block = block
}

func (rawStore *RawStore) load(height int32) *api.Block {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	in, err := ioutil.ReadFile(path)
	if err != nil {
		Log.Fatal("Error reading file:", zap.Error(err))
	}
	block := &api.Block{}
	if err := proto.Unmarshal(in, block); err != nil {
		Log.Fatal("Error parse block file:", zap.Error(err))
	}

	return block
}

func (rawStore *RawStore) StoreAndMaybeConsume(blocks []*api.Block) error {
	readyToConsume := false
	waitForHeight := int(rawStore.Current.Height + 1)

	// store
	txn := rawStore.queue.NewTransaction(true)
	for _, block := range blocks {
		if blockPb, err := proto.Marshal(block); err == nil {
			height := int(block.Height)
			heightBs := []byte(strconv.Itoa(height))
			for {
				if err := txn.Set(heightBs, blockPb); err == nil {
					if waitForHeight == height {
						readyToConsume = true
					}
					break
				} else {
					switch err {
					case badger.ErrTxnTooBig:
						// split up big transactions
						_ = txn.Commit(nil)
						txn = rawStore.queue.NewTransaction(true)
						continue
					default:
						panic(err)
					}
				}
			}
		}
	}
	_ = txn.Commit(nil)

	if !readyToConsume {
		return nil
	}

	// consume, cause looks like we just stored at least the height we we're waiting for
	txn = rawStore.queue.NewTransaction(true)
	for {
		height := waitForHeight
		heightBs := []byte(strconv.Itoa(height))
		waitForHeight++

		if blockItem, err := txn.Get(heightBs); err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				// looks like we have no further blocks which matches the required sequence atm
				return nil
			default:
				return err
			}
		} else {
			if blockBs, err := blockItem.Value(); err != nil {
				return err
			} else {
				block := new(api.Block)
				if err := proto.Unmarshal(blockBs, block); err != nil {
					return err
				} else {
					rawStore.Store(block, int32(height))
				RETRY_DELETE:
					if err := txn.Delete(heightBs); err != nil {
						switch err {
						case badger.ErrTxnTooBig:
							// split up big transactions
							_ = txn.Commit(nil)
							txn = rawStore.queue.NewTransaction(true)
							goto RETRY_DELETE
						default:
							panic(err)
						}
					}
				}
			}
		}
	}
}
