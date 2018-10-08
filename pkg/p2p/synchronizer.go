package p2p

import (
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	. "github.com/ac0v/aspera/pkg/log"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	"go.uber.org/zap"
)

type Synchronizer struct {
	registry *r.Registry
	store    *s.Store
	client   Client

	blockRanges        chan *blockRange
	blockBatchesEmpty  chan *blockBatch
	blockBatchesFilled chan *blockBatch
}

type blockRange struct {
	from blockMeta
	to   *blockMeta
}

type blockMeta struct {
	id     uint64
	height int32
}

type blockBatch struct {
	blockRange *blockRange

	peers  []*Peer
	ids    []uint64
	blocks []*pb.Block
}

func NewSynchronizer(client Client, store *s.Store, registry *r.Registry) *Synchronizer {
	milestones := registry.Config.Network.P2P.Milestones

	s := &Synchronizer{
		client:             client,
		registry:           registry,
		store:              store,
		blockRanges:        make(chan *blockRange, len(milestones)),
		blockBatchesEmpty:  make(chan *blockBatch, len(milestones)),
		blockBatchesFilled: make(chan *blockBatch),
	}

	for i := 0; i < 20; i++ {
		go s.fetchBlocks()
	}

	for i, milestone := range milestones {
		if store.RawStore.Current.Block.Height > milestone.Height {
			continue
		}

		var toBlockMeta *blockMeta
		if len(milestones) > i+1 {
			toBlockMeta = &blockMeta{
				id:     milestones[i+1].Id,
				height: milestones[i+1].Height,
			}
		}

		s.blockRanges <- &blockRange{
			from: blockMeta{
				id:     milestone.Id,
				height: milestone.Height,
			},
			to: toBlockMeta,
		}
	}

	// debug
	var processed int
	start := time.Now()
	for blockBatch := range s.blockBatchesFilled {
		processed += len(blockBatch.blocks)
		Log.Info(
			"syncing with",
			zap.Float64("blocks/s", float64(processed)/time.Since(start).Seconds()))
	}
	// end

	return s
}

func (s *Synchronizer) fetchBlocks() {
	for {
		select {
		case blockRange := <-s.blockRanges:
			currentID := blockRange.from.id

			res, peers, err := s.client.GetNextBlockIDs(currentID)
			if err != nil {
				Log.Error("get next blocks ids", zap.Error(err))
				continue
			}

			ids := res.NextBlockIds

			idCount := len(ids)
			if idCount == 0 {
				time.Sleep(10 * time.Second)
				s.blockRanges <- blockRange
				continue
			}

			Log.Info("syncing block", zap.Uint64("id", currentID))

			s.blockBatchesEmpty <- &blockBatch{
				blockRange: blockRange,
				ids:        ids,
				peers:      peers,
			}
		case blockBatch := <-s.blockBatchesEmpty:
			currentID := blockBatch.ids[0]

		FETCH_BLOCKS_AGAIN:
			res, peers, err := s.client.GetNextBlocks(currentID)
			if err != nil {
				Log.Error("get next blocks", zap.Error(err))
				goto FETCH_BLOCKS_AGAIN
			}

			blockBatch.blocks = res.NextBlocks

			validBlocks := countValidBlocks(blockBatch)

			if validBlocks == 0 {
				s.blockRanges <- blockBatch.blockRange
				continue
			}

			from := blockBatch.blockRange.from
			to := blockBatch.blockRange.to

			from.height += int32(validBlocks)
			from.id = blockBatch.ids[validBlocks-1]
			blockBatch.blocks = blockBatch.blocks[:validBlocks]

			// TODO: kill some sync threads when we reached end of a block range?
			if to == nil || from.height <= to.height {
				s.blockRanges <- &blockRange{
					from: from,
					to:   to,
				}
			}

			blockBatch.peers = append(blockBatch.peers, peers...)

			s.blockBatchesFilled <- blockBatch
		}
	}
}

func countValidBlocks(blockBatch *blockBatch) int {
	ids := blockBatch.ids
	blocks := blockBatch.blocks

	idCount := len(ids)
	blockCount := len(blocks)

	var iEnd int
	if idCount < blockCount {
		iEnd = idCount
	} else {
		iEnd = blockCount
	}

	var validBlocks int
	for ; validBlocks < iEnd; validBlocks++ {
		if blocks[validBlocks].PreviousBlock != ids[validBlocks] {
			return validBlocks
		}
	}

	return validBlocks
}
