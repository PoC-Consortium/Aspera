package p2p

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/PoC-Consortium/aspera/pkg/block"
	"github.com/PoC-Consortium/aspera/pkg/config"
	. "github.com/PoC-Consortium/aspera/pkg/log"
	"github.com/PoC-Consortium/aspera/pkg/p2p/manager"
	s "github.com/PoC-Consortium/aspera/pkg/store"
)

type Synchronizer struct {
	store  *s.Store
	client Client

	blockBatchesEmpty  chan *blockBatch // range information
	blockBatchesIDs    chan *blockBatch // ids for fetching blocks
	blockBatchesBlocks chan *blockBatch // blocks for validation

	glueBlocks   map[int32]*block.Block
	glueBlocksMu sync.Mutex

	statistic   *statistic
	statisticMu sync.Mutex
}

type statistic struct {
	syncSpeed float64
	processed int
	tick      time.Time
}

type blockMeta struct {
	id     uint64
	height int32
}

type blockBatch struct {
	from blockMeta
	to   blockMeta

	peers []manager.Peer

	ids    []uint64
	blocks []*block.Block
}

func (s *statistic) update(blockCount int) (float64, int) {
	s.processed += blockCount
	now := time.Now()
	d := now.Sub(s.tick).Seconds()
	s.syncSpeed = 0.95*s.syncSpeed + 0.05*float64(blockCount)/d
	s.tick = now
	return s.syncSpeed, s.processed
}

func NewSynchronizer(client Client, store *s.Store, milestones []config.Milestone) *Synchronizer {
	s := &Synchronizer{
		statistic: &statistic{tick: time.Now()},
		client:    client,
		store:     store,

		blockBatchesEmpty:  make(chan *blockBatch, len(milestones)),
		blockBatchesIDs:    make(chan *blockBatch, len(milestones)),
		blockBatchesBlocks: make(chan *blockBatch, len(milestones)),

		glueBlocks: make(map[int32]*block.Block),
	}

	// the current block is a predecessor for the first glue action
	if currentBlock, err := block.NewBlock(store.RawStore.Current.Block); err != nil {
		panic(err)
	} else {
		s.glueBlocks[currentBlock.Height] = currentBlock

		for i := 0; i < len(milestones); i++ {
			go s.fetchBlockIDs()
			go s.fetchBlocks()
			go s.validateBlocks()
		}

		bms := make([]*blockMeta, len(milestones))
		for i, m := range milestones {
			bms[i] = &blockMeta{id: m.Id, height: m.Height}
		}

		blockBatches := alignMilestonesWithCurrent(bms, &blockMeta{id: currentBlock.Id, height: currentBlock.Height})
		for _, blockBatch := range blockBatches {
			s.blockBatchesEmpty <- blockBatch
		}
	}

	select {}

	return s
}

func alignMilestonesWithCurrent(milestones []*blockMeta, current *blockMeta) []*blockBatch {
	var blockBatches []*blockBatch
	for i, milestone := range milestones[:len(milestones)-1] {
		if milestones[i+1].height <= current.height {
			continue
		}

		// milestone already handled (partialy) - adjust it's start
		if current.height > milestone.height {
			milestone.id = current.id
			milestone.height = current.height
		}

		blockBatches = append(
			blockBatches,
			&blockBatch{
				from: blockMeta{
					id:     milestone.id,
					height: milestone.height,
				},
				to: blockMeta{
					id:     milestones[i+1].id,
					height: milestones[i+1].height,
				},
			},
		)
	}
	return blockBatches
}

func (s *Synchronizer) refetchBlockRangeAndBlockPeers(blockBatch *blockBatch, err error) {
	Log.Error("got invalid blocks", zap.Error(err))
	for _, p := range blockBatch.peers {
		p.Throttle()
	}
	s.blockBatchesEmpty <- blockBatch
}

func (s *Synchronizer) fetchBlockIDs() {
	for blockBatch := range s.blockBatchesEmpty {
		currentID := blockBatch.from.id
		currentHeight := blockBatch.from.height

		// TODO: if we can't get any blocks we probably have to go back to an earlier id
	fetchBlockIDs:
		for {
			res, peers, err := s.client.GetNextBlockIDs(currentID, currentHeight)
			if err != nil {
				Log.Error("get next blocks ids", zap.Error(err))
				continue fetchBlockIDs
			}

			ids := res.NextBlockIds
			if len(ids) == 0 {
				time.Sleep(10 * time.Second)
				Log.Error("get next blocks ids", zap.Error(errors.New("empty response")))
				continue fetchBlockIDs
			}
			blockBatch.ids = ids
			blockBatch.peers = peers
			s.blockBatchesIDs <- blockBatch
			break fetchBlockIDs
		}
	}
}

func (s *Synchronizer) fetchBlocks() {
	for blockBatch := range s.blockBatchesIDs {
		currentID := blockBatch.from.id
		currentHeight := blockBatch.from.height

		// TODO: if we can't get any blocks we probably have to go back to an earlier id
	fetchBlocks:
		for {
			res, peers, err := s.client.GetNextBlocks(currentID, currentHeight)
			if err != nil {
				for _, p := range peers {
					p.Throttle()
				}
				Log.Error("get next blocks", zap.Error(err))
				continue fetchBlocks
			}
			if len(res.NextBlocks) == 0 {
				for _, p := range peers {
					p.Throttle()
				}
				Log.Error("get next blocks", zap.Error(errors.New("empty response")))
				continue fetchBlocks
			}

			blocks := make([]*block.Block, len(res.NextBlocks))
			for i, b := range res.NextBlocks {
				if blocks[i], err = block.NewBlock(b); err != nil {
					Log.Error("put api block into wrapper", zap.Error(err))
					continue fetchBlocks
				}
			}
			blockBatch.blocks = blocks
			validBlockCount := countValidBlocksAndSetHeight(currentID, blockBatch)
			if validBlockCount == 0 {
				Log.Error("got not id overlaps for blocks and ids", zap.Uint64("block id", currentID))
				continue fetchBlocks
			} else {
				blockBatch.blocks = blockBatch.blocks[:validBlockCount]
			}

			blockBatch.peers = append(blockBatch.peers, peers...)
			s.blockBatchesBlocks <- blockBatch
			break
		}
	}
}

func (s *Synchronizer) validateBlocks() {
validateBlocks:
	for blockBatch := range s.blockBatchesBlocks {
		blocks := blockBatch.blocks

		numNewBlocks := len(blocks)
		firstBlock := blocks[0]
		lastBlock := blocks[len(blocks)-1]

		s.glueBlocksMu.Lock()
		// we need the predecessor of the first block of this batch
		// if we have it cached we can simply prepend it
		// otherwise we need to cache the first block for later validation
		if b, exists := s.glueBlocks[firstBlock.Height-1]; exists {
			blocks = append([]*block.Block{b}, blocks...)
		} else {
			s.glueBlocks[firstBlock.Height] = firstBlock
		}

		// the last block of this batch could be the predecessor of another block
		// that still needs to be validated
		if b, exists := s.glueBlocks[lastBlock.Height+1]; exists {
			blocks = append(blocks, b)
			delete(s.glueBlocks, b.Height)
		} else {
			s.glueBlocks[lastBlock.Height] = lastBlock
		}
		s.glueBlocksMu.Unlock()

		for i, b := range blocks[1:len(blocks)] {
			previousBlock := blocks[i]
			// ToDo: after validation fails we should probably refetch the block ids
			if err := b.PreValidate(previousBlock); err != nil {
				s.refetchBlockRangeAndBlockPeers(blockBatch, err)
				continue validateBlocks
			}
		}

		// ToDo: may we should hand over the block batch to allow blocking bad peers?
		// sync speed

		s.store.RawStore.StoreAndMaybeConsume(blocks[1:len(blocks)])
		blockBatch.from.height += int32(numNewBlocks)
		blockBatch.from.id = blocks[len(blocks)-1].Id

		s.statisticMu.Lock()
		syncSpeed, processed := s.statistic.update(numNewBlocks)
		s.statisticMu.Unlock()
		Log.Info("syncing with", zap.Float64("blocks/s", syncSpeed), zap.Int("processed", processed))

		if blockBatch.from.height < blockBatch.to.height {
			s.blockBatchesEmpty <- blockBatch
		}
	}
}

func countValidBlocksAndSetHeight(currentId uint64, blockBatch *blockBatch) int {
	ids := blockBatch.ids
	blocks := blockBatch.blocks
	for i, id := range append([]uint64{currentId}, ids...) {
		if i > len(blocks)-1 || blocks[i].PreviousBlock != id {
			return i
		}
		blocks[i].Height = blockBatch.from.height + int32(i+1)
	}
	return len(ids)
}
