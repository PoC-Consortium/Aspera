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
	stopAt r.Milestone
}

func NewSynchronizer(client *Client, store *s.Store, registry *r.Registry) *Synchronizer {
	synchronizer := &Synchronizer{client: client, registry: registry, store: store, wg: new(sync.WaitGroup)}

	fetchBlocksChannel := make(chan *blockMeta)

	synchronizer.wg.Add(len(registry.Config.Network.P2P.Milestones) + 1)

	// fetch block Ids after each milestone - in parallel
	for milestoneIndex, milestone := range registry.Config.Network.P2P.Milestones {
		var stopAt r.Milestone
		if len(registry.Config.Network.P2P.Milestones) > milestoneIndex+1 {
			stopAt = registry.Config.Network.P2P.Milestones[milestoneIndex+1]
		}
		go synchronizer.fetchBlockIds(fetchBlocksChannel, milestone, stopAt)
	}
	go synchronizer.fetchBlocks(fetchBlocksChannel)
	synchronizer.wg.Wait()

	return synchronizer
}

func (synchronizer *Synchronizer) fetchBlockIds(fetchBlocksChannel chan *blockMeta, milestone r.Milestone, stopAt r.Milestone) {
	previousBlockId := milestone.Id
	height := milestone.Height
	for {
		synchronizer.registry.Logger.Info("syncing block meta", zap.Uint64("id", previousBlockId), zap.Int("height", int(height)), zap.Int("pending", len(fetchBlocksChannel)))
		fetchBlocksChannel <- &blockMeta{id: previousBlockId, height: height, stopAt: stopAt}

		// if we got IDs after the next milestone - end this block ID fetcher
		if &stopAt != nil && height > stopAt.Height {
			break
		}

		var res *pb.GetNextBlockIdsResponse
		for {
			var err error
			res, err = synchronizer.client.GetNextBlockIds(previousBlockId)
			if err != nil {
				synchronizer.registry.Logger.Error("getNextBlocks", zap.Error(err))
				continue
			} else if len(res.NextBlockIds) == 0 {
				if height < synchronizer.registry.Config.Network.P2P.Milestones[len(synchronizer.registry.Config.Network.P2P.Milestones)-1].Height {
					// did not reach the height of the last milestone block
					// - so it looks like a temporary network or peer issue
					continue
				} else {
					// wait before asking for fresh blocks - looks like there are no blocks atm around
					time.Sleep(time.Second * 10)
				}
			}

			break
		}

		takeIndex := len(res.NextBlockIds) - 1
		if height != milestone.Height {
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
				synchronizer.registry.Logger.Error("getNextBlocks", zap.Error(err))
				continue
			}

			for _, block := range res.NextBlocks {
				height++
				//	synchronizer.registry.Logger.Info("syncing block", zap.Uint64("id", block.Block), zap.Uint64("previousBlockId", block.PreviousBlock), zap.Int("height", int(height)))
				synchronizer.store.RawStore.Store(block, height)

				// if there is a further milestone, we stop this processing
				if &blockMeta.stopAt != nil && height >= blockMeta.stopAt.Height {
					break
				}
			}
			break
		}
	}
	synchronizer.wg.Done()
}
