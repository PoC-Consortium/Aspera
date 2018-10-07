package p2p

import (
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	"go.uber.org/zap"
)

type Synchronizer struct {
	registry *r.Registry
	store    *s.Store
	client   *Client

	blockRanges        chan *blockRange
	blockBatchesEmpty  chan *blockBatch
	blockBatchesFilled chan *blockBatch
}

type blockRange struct {
	from *blockMeta
	to   *blockMeta
}

type blockMeta struct {
	id     uint64
	height int32
}

type blockBatch struct {
	requestedID uint64

	peers  []*Peer
	ids    []uint64
	blocks []*pb.Block
}

func NewSynchronizer(client *Client, store *s.Store, registry *r.Registry) *Synchronizer {
	milestones := registry.Config.Network.P2P.Milestones

	s := &Synchronizer{
		client:             client,
		registry:           registry,
		store:              store,
		blockRanges:        make(chan *blockRange, len(milestones)),
		blockBatchesEmpty:  make(chan *blockBatch),
		blockBatchesFilled: make(chan *blockBatch),
	}

	go s.fetchBlockIds()
	for i := 0; i < 4; i++ {
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
			from: &blockMeta{
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
		s.registry.Logger.Info(
			"syncing with",
			zap.Float64("blocks/s", float64(processed)/time.Since(start).Seconds()))
	}
	// end

	return s
}

func (s *Synchronizer) fetchBlockIds() {
	for blockRange := range s.blockRanges {
		currentHeight := blockRange.from.height
		currentID := blockRange.from.id

		for currentHeight <= blockRange.to.height {
			res, peers, err := s.client.GetNextBlockIds(currentID)
			if err != nil {
				s.registry.Logger.Error("get next blocks", zap.Error(err))
				continue
			}

			ids := res.NextBlockIds

			idCount := len(ids)
			if idCount == 0 {
				continue
			}

			s.registry.Logger.Info("syncing block", zap.Uint64("id", currentID))

			currentHeight += int32(idCount)
			currentID = ids[idCount-1]

			s.blockBatchesEmpty <- &blockBatch{
				requestedID: currentID,
				ids:         ids,
				peers:       peers,
			}
		}
	}
}

func (s *Synchronizer) fetchBlocks() {
	for blockBatch := range s.blockBatchesEmpty {
		currentID := blockBatch.ids[0]

	FETCH_BLOCKS_AGAIN:
		res, peers, err := s.client.GetNextBlocks(currentID)
		if err != nil {
			s.registry.Logger.Error("get next blocks", zap.Error(err))
			goto FETCH_BLOCKS_AGAIN
		}

		blockBatch.blocks = res.NextBlocks

		if !idsMatchBlocks(blockBatch) {
			goto FETCH_BLOCKS_AGAIN
		}

		blockBatch.peers = append(blockBatch.peers, peers...)

		s.blockBatchesFilled <- blockBatch
	}
}

func idsMatchBlocks(blockBatch *blockBatch) (valid bool) {
	ids := blockBatch.ids
	blocks := blockBatch.blocks

	if len(ids) != len(blocks) {
		return false
	}

	for i := 0; i < len(blocks); i++ {
		if blocks[i].PreviousBlock != ids[i] {
			return false
		}
	}

	return true
}
