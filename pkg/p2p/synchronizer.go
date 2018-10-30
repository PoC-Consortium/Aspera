package p2p

import (
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/block"
	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"
	s "github.com/ac0v/aspera/pkg/store"
)

type Synchronizer struct {
	store     *s.Store
	client    Client
	manager   Manager
	statistic *statistic

	blockRanges        chan *blockRange
	blockBatchesEmpty  chan *blockBatch
	blockBatchesFilled chan *blockBatch
	blockBatchesGlue   chan []*api.Block

	glueBlockOf *sync.Map
}

type statistic struct {
	start     time.Time
	processed int32
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

	peers  []string
	ids    []uint64
	blocks []*api.Block

	isGlueResult bool
}

func NewSynchronizer(client Client, manager Manager, store *s.Store, milestones []config.Milestone) *Synchronizer {
	s := &Synchronizer{
		statistic: &statistic{
			start:     time.Now(),
			processed: 0,
		},
		client:             client,
		manager:            manager,
		store:              store,
		blockRanges:        make(chan *blockRange, len(milestones)),
		blockBatchesEmpty:  make(chan *blockBatch, len(milestones)),
		blockBatchesFilled: make(chan *blockBatch),
		blockBatchesGlue:   make(chan []*api.Block),
		glueBlockOf:        &sync.Map{},
	}

	for i := 0; i < 20; i++ {
		go s.fetchBlocks()
		go s.validateBlocks()
	}

	for i, milestone := range milestones {
		// milestone already handled (partialy) - adjust it's start
		if store.RawStore.Current.Block.Height > milestone.Height {
			milestone.Id = store.RawStore.Current.Block.Id
			milestone.Height = store.RawStore.Current.Block.Height
		}

		var toBlockMeta *blockMeta
		if len(milestones) > i+1 {
			toBlockMeta = &blockMeta{
				id:     milestones[i+1].Id,
				height: milestones[i+1].Height,
			}
			if toBlockMeta.height <= store.RawStore.Current.Block.Height {
				// milestone already done; next one partially handled
				// so we continue with the next one
				continue
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
	//var processed int
	//	start := time.Now()

	/*
		for blockBatch := range s.blockBatchesFilled {
			//panic(blockBatch.ToBytes())
			processed += len(blockBatch.blocks)
			Log.Info(
				"syncing with",
				zap.Float64("blocks/s", float64(processed)/time.Since(start).Seconds()))
		}
	*/
	// end
	select {}

	return s
}

func (s *Synchronizer) refetchBlockRangeAndBlockPeers(blockBatch *blockBatch, err error) {
	Log.Error("got invalid blocks", zap.Error(err))
	for _, p := range blockBatch.peers {
		s.manager.BlockPeer(p, PeerDataIntegrityValidation)
	}
	s.blockRanges <- blockBatch.blockRange
}

func (s *Synchronizer) validateBlocks() {
ValidateBlocks:
	for {
		select {
		case blockBatch := <-s.blockBatchesFilled:
			// if we have more than one block, we can do the basic validation (eg. hashing)
			// for all blocks except of the first one where the previous is unknown
			blocks := blockBatch.blocks

			if len(blocks) < 2 {
				// can't validate a single block without any previous
				s.blockBatchesGlue <- []*api.Block{blockBatch.blocks[0]}
				continue
			}

			blockWrappers := make([]*block.Block, len(blocks))
			var err error
			for i, b := range blocks {
				if blockWrappers[i], err = block.NewBlock(b); err != nil {
					s.refetchBlockRangeAndBlockPeers(blockBatch, err)
					continue ValidateBlocks
				}
			}

			for i, b := range blockWrappers[1 : len(blockWrappers)-1] {
				// blocks[i] is the previousBlock .. cause that's a slice above
				// - starting with element no. 2
				if err = b.Validate(blockWrappers[i]); err != nil {
					s.refetchBlockRangeAndBlockPeers(blockBatch, err)
					continue ValidateBlocks
				} else {
					blockWrappers[i].Id = b.PreviousBlock
					Log.Info("syncing block", zap.Int32("height", b.Height), zap.Uint64("id", b.Id))
				}
			}

			// ToDo: may we should hand over the block batch to allow blocking bad peers?
			s.store.RawStore.StoreAndMaybeConsume(blockBatch.blocks)
			storedCount := int32(len(blocks) - 1)

			// we need to store only the successor (last block) of a glue result
			// .. has already been done by the loop around blocks[1:len(blocks)-1]
			if !blockBatch.isGlueResult {
				if blocks[0].Height == 0 {
					// ToDo: may we should hand over the block batch to allow blocking bad peers?
					s.store.RawStore.StoreAndMaybeConsume(blockBatch.blocks)
					storedCount++
				} else {
					s.blockBatchesGlue <- []*api.Block{
						blocks[0],
						blocks[len(blocks)-1],
					}
				}
			}

			processedCount := atomic.AddInt32(&s.statistic.processed, storedCount)
			Log.Info(
				"syncing with",
				zap.Int32(
					"blocks/s",
					int32(float64(processedCount)/time.Since(s.statistic.start).Seconds()),
				),
			)
		case blocks := <-s.blockBatchesGlue:
			orphanedBlock := blocks[0]
			// think we do not need any lock between load and delete,
			// cause this matching pair will not be handled by someone else
			if previousBlock, ok := s.glueBlockOf.Load(orphanedBlock.Height - 1); ok {
				s.glueBlockOf.Delete(orphanedBlock.Height)
				previousBlock := previousBlock.(*api.Block)
				s.blockBatchesFilled <- &blockBatch{
					blockRange: &blockRange{
						from: blockMeta{
							id:     previousBlock.Id,
							height: previousBlock.Height,
						},
						to: &blockMeta{
							id:     orphanedBlock.Id,
							height: orphanedBlock.Height,
						},
					},
					blocks:       []*api.Block{previousBlock, orphanedBlock},
					ids:          []uint64{previousBlock.Id, orphanedBlock.Id},
					isGlueResult: true,
				}
				Log.Info("glue pair found", zap.Int32("leftHeight", previousBlock.Height), zap.Int32("rightHeight", orphanedBlock.Height))
			} else {
				// keep track of possible successor
				s.glueBlockOf.Store(orphanedBlock.Height, orphanedBlock)
			}
			if len(blocks) > 1 {
				lastBlock := blocks[len(blocks)-1]
				// keep track of - possible predecessor
				s.glueBlockOf.Store(lastBlock.Height, lastBlock)
			}
		}
	}
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
				ids:        append([]uint64{currentID}, ids...),
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

			validBlocks := countValidBlocksAndSetHeight(blockBatch)

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

func countValidBlocksAndSetHeight(blockBatch *blockBatch) int {
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
		blocks[validBlocks].Height = blockBatch.blockRange.from.height + int32(validBlocks)
		if blocks[validBlocks].PreviousBlock != ids[validBlocks] {
			return validBlocks
		}
	}

	return validBlocks
}
