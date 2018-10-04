package p2p

import (
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Synchronizer struct {
	registry *r.Registry
	store    *s.Store
	client   *Client
	wg       *sync.WaitGroup
}

type blockMeta struct {
	id     uint64
	height int32
}

func NewSynchronizer(client *Client, store *s.Store, registry *r.Registry) *Synchronizer {
	synchronizer := &Synchronizer{client: client, registry: registry, store: store, wg: new(sync.WaitGroup)}

	fetchBlocksChannel := make(chan *blockMeta)

	synchronizer.wg.Add(2)
	go synchronizer.fetchBlockIds(fetchBlocksChannel)
	go synchronizer.fetchBlocks(fetchBlocksChannel)
	synchronizer.wg.Wait()

	return synchronizer
}

func (synchronizer *Synchronizer) fetchBlockIds(fetchBlocksChannel chan *blockMeta) {
	latestBlock := synchronizer.store.RawStore.Current.Block
	previousBlockId := latestBlock.PreviousBlock
	if latestBlock.Height == 0 {
		previousBlockId = uint64(latestBlock.Block)
	}
	height := latestBlock.Height
	for {
		synchronizer.registry.Logger.Info("syncing block meta", zap.Uint64("id", previousBlockId), zap.Int("height", int(height)))
		fetchBlocksChannel <- &blockMeta{id: previousBlockId, height: height}

		var res *pb.GetNextBlockIdsResponse
		for {
			var err error
			res, err = synchronizer.client.GetNextBlockIds(previousBlockId)
			if err != nil {
				continue
			} else if len(res.NextBlockIds) == 0 {
				// wait before asking for fresh blocks - looks like there are no blocks atm around
				time.Sleep(time.Second * 10)
			}

			break
		}

		takeIndex := len(res.NextBlockIds) - 1
		if height != 0 {
			// atm we do not know the blockId, but it's previous
			// - so we ignore the double returned block
			takeIndex--
		}

		height += int32(takeIndex + 1)
		previousBlockId = res.NextBlockIds[takeIndex]
	}

	synchronizer.wg.Done()
}

func (synchronizer *Synchronizer) fetchBlocks(fetchBlocksChannel chan *blockMeta) {
	for blockMeta := range fetchBlocksChannel {
		// try to get the block data - till we have them!
		for {
			height := blockMeta.height

			synchronizer.registry.Logger.Info("syncing block", zap.Uint64("id", blockMeta.id), zap.Int("height", int(blockMeta.height)))
			res, err := synchronizer.client.GetNextBlocks(blockMeta.id)
			// redo on exceptions; may another per has better data
			if err != nil || len(res.NextBlocks) == 0 {
				//synchronizer.registry.Logger.Info("err:", zap.Error(err))
				continue
			}

			for _, block := range res.NextBlocks {
				height++
				//	synchronizer.registry.Logger.Info("syncing block", zap.Uint64("id", block.Block), zap.Uint64("previousBlockId", block.PreviousBlock), zap.Int("height", int(height)))
				synchronizer.store.RawStore.Store(block, height)
			}
			break
		}
	}
	synchronizer.wg.Done()
}
