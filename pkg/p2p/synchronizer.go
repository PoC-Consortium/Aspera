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
	blockBatchesGlue   chan []*block.Block

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
		blockBatchesGlue:   make(chan []*block.Block),
		glueBlockOf:        &sync.Map{},
	}

	// the current block is a predecessor for the first glue action
	if currentBlock, err := block.NewBlock(store.RawStore.Current.Block); err != nil {
		panic(err)
	} else {
		go func() {
			s.blockBatchesGlue <- []*block.Block{currentBlock}
		}()

		for i := 0; i < 20; i++ {
			go s.fetchBlocks()
			go s.validateBlocks()
		}

		var bm []*blockMeta
		for _, m := range milestones {
			bm = append(bm, &blockMeta{id: m.Id, height: m.Height})
		}

		for _, blockRange := range alignMilestonesWithCurrent(bm, &blockMeta{id: currentBlock.Id, height: currentBlock.Height}) {
			s.blockRanges <- blockRange
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

func alignMilestonesWithCurrent(milestones []*blockMeta, current *blockMeta) []*blockRange {
	var blockRanges []*blockRange
	for i, milestone := range milestones {
		// milestone already handled (partialy) - adjust it's start
		if current.height > milestone.height {
			milestone.id = current.id
			milestone.height = current.height
		}

		var toBlockMeta *blockMeta
		if len(milestones) > i+1 {
			toBlockMeta = &blockMeta{
				id:     milestones[i+1].id,
				height: milestones[i+1].height,
			}
			if toBlockMeta.height <= current.height {
				// milestone already done; next one partially handled
				// so we continue with the next one
				continue
			}
		}

		blockRanges = append(
			blockRanges,
			&blockRange{
				from: blockMeta{
					id:     milestone.id,
					height: milestone.height,
				},
				to: toBlockMeta,
			},
		)
	}

	return blockRanges
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

			blockWrappers := make([]*block.Block, len(blocks))
			var err error
			for i, b := range blocks {
				if blockWrappers[i], err = block.NewBlock(b); err != nil {
					s.refetchBlockRangeAndBlockPeers(blockBatch, err)
					continue ValidateBlocks
				}
			}

			storedCount := 0
			if len(blockWrappers) > 1 {
				// can only validate a block series with at least 2 elements (aka. block with previousBlock)
				for i, b := range blockWrappers[1:len(blockWrappers)] {
					// blocks[i] is the previousBlock .. cause that's a slice above
					// - starting with element no. 2
					if err = b.Validate(blockWrappers[i]); err != nil {
						s.refetchBlockRangeAndBlockPeers(blockBatch, err)
						continue ValidateBlocks
					} else {
						Log.Info("syncing block", zap.Int32("height", b.Height), zap.Uint64("id", b.Id))
					}
				}

				// ToDo: may we should hand over the block batch to allow blocking bad peers?
				s.store.RawStore.StoreAndMaybeConsume(blocks[1:len(blocks)])
				storedCount = len(blocks) - 1

				// if this was no glue action we need to store the possible successor and predecessor
				if !blockBatch.isGlueResult {
					s.blockBatchesGlue <- []*block.Block{
						blockWrappers[0],
						blockWrappers[len(blocks)-1],
					}
				}
			} else {
				// the first element is a possible predecessor - except for a handled glue result
				s.blockBatchesGlue <- []*block.Block{blockWrappers[0]}
			}

			//if blocks[0].Height == 0 {
			// ToDo: may we should hand over the block batch to allow blocking bad peers?
			//s.store.RawStore.StoreAndMaybeConsume(blockBatch.blocks)
			//storedCount++

			processedCount := atomic.AddInt32(&s.statistic.processed, int32(storedCount))
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
				previousBlock := previousBlock.(*block.Block)
				blockBatch := &blockBatch{
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
					blocks: []*api.Block{previousBlock.Block, orphanedBlock.Block},
					ids:    []uint64{previousBlock.Id, orphanedBlock.Id},
				}

				s.glueBlockOf.Delete(orphanedBlock.Height)
				blockBatch.isGlueResult = true
				s.blockBatchesFilled <- blockBatch

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
				ids:        ids,
				peers:      peers,
			}
		case blockBatch := <-s.blockBatchesEmpty:
			currentID := blockBatch.blockRange.from.id

		FETCH_BLOCKS_AGAIN:
			res, peers, err := s.client.GetNextBlocks(currentID)
			if err != nil {
				Log.Error("get next blocks", zap.Error(err))
				goto FETCH_BLOCKS_AGAIN
			}

			blockBatch.blocks = res.NextBlocks

			validBlockCount := countValidBlocksAndSetHeight(currentID, blockBatch)

			if validBlockCount == 0 {
				s.blockRanges <- blockBatch.blockRange
				continue
			}

			from := blockBatch.blockRange.from
			to := blockBatch.blockRange.to

			from.height += int32(validBlockCount)
			from.id = blockBatch.ids[validBlockCount-1]

			blockBatch.blocks = blockBatch.blocks[:validBlockCount]

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

func countValidBlocksAndSetHeight(currentId uint64, blockBatch *blockBatch) int {
	ids := blockBatch.ids
	blocks := blockBatch.blocks

	for i, id := range append([]uint64{currentId}, ids...) {
		if i > len(blocks)-1 || blocks[i].PreviousBlock != id {
			return i
		}
		blocks[i].Height = blockBatch.blockRange.from.height + int32(i+1)
	}

	return len(ids)
}
