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

type blockRange struct {
	from blockMeta
	to   blockMeta
	ids  []uint64
}

type blockMeta struct {
	id     uint64
	height int32
}

func NewSynchronizer(client *Client, store *s.Store, registry *r.Registry) *Synchronizer {
	synchronizer := &Synchronizer{client: client, registry: registry, store: store, wg: new(sync.WaitGroup)}

	blockRanges := []blockRange{}

	// fetch block Ids after each milestone - in parallel
	for milestoneIndex, milestone := range registry.Config.Network.P2P.Milestones {
		if store.RawStore.Current.Block.Height <= milestone.Height {
			var toBlockMeta blockMeta
			if len(registry.Config.Network.P2P.Milestones) > milestoneIndex+1 {
				toBlockMeta = blockMeta{
					id:     registry.Config.Network.P2P.Milestones[milestoneIndex+1].Id,
					height: registry.Config.Network.P2P.Milestones[milestoneIndex+1].Height,
				}
			}
			blockRanges = append(
				blockRanges,
				blockRange{
					from: blockMeta{id: milestone.Id, height: milestone.Height},
					to:   toBlockMeta,
				},
			)
		}
	}

	fetchBlocksChannel := make(chan *blockRange)
	synchronizer.wg.Add(len(blockRanges) * 2)
	for _, blockRange := range blockRanges {
		go synchronizer.fetchBlockIds(fetchBlocksChannel, blockRange)
		//go synchronizer.fetchBlocks(fetchBlocksChannel)
		go synchronizer.fetchBlocks(fetchBlocksChannel)
	}
	synchronizer.wg.Wait()

	return synchronizer
}

func (synchronizer *Synchronizer) fetchBlockIds(fetchBlocksChannel chan *blockRange, blockRange blockRange) {
	currentHeight := blockRange.from.height
	for {
		var res *pb.GetNextBlockIdsResponse
		for {
			var err error
			res, err = synchronizer.client.GetNextBlockIds(blockRange.from.id)
			if err != nil {
				synchronizer.registry.Logger.Error("getNextBlockIds", zap.Error(err))
				continue
			} else if len(res.NextBlockIds) == 0 {
				if currentHeight < synchronizer.registry.Config.Network.P2P.Milestones[len(synchronizer.registry.Config.Network.P2P.Milestones)-1].Height {
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

		if (blockRange.to.id == 0 && blockRange.to.height == 0) || ((currentHeight + int32(len(res.NextBlockIds))) < blockRange.to.height) {
			blockRange.ids = res.NextBlockIds
			fetchBlocksChannel <- &blockRange
			currentHeight += int32(len(res.NextBlockIds))
			blockRange.from.id = res.NextBlockIds[len(res.NextBlockIds)-1]
			blockRange.from.height = currentHeight
		} else if (currentHeight + int32(len(res.NextBlockIds))) >= blockRange.to.height {
			// take only all elements before and the to-element itself
			blockRange.ids = res.NextBlockIds[0 : 1+(int32(len(res.NextBlockIds))-blockRange.to.height-currentHeight)]
			fetchBlocksChannel <- &blockRange
			break
		}
	}

	synchronizer.wg.Done()
}

func (synchronizer *Synchronizer) fetchBlocks(fetchBlocksChannel chan *blockRange) {
	for blockRange := range fetchBlocksChannel {
		currentHeight := blockRange.from.height
		ids := blockRange.ids

		// try to get the block data - till we have them!
		for {
			res, err := synchronizer.client.GetNextBlocks(ids[0])
			// redo on exceptions; may another peer has better data
			if err != nil || len(res.NextBlocks) == 0 {
				synchronizer.registry.Logger.Error("getNextBlocks", zap.Error(err))
				continue
			}
			synchronizer.registry.Logger.Info("syncing block", zap.Uint64("id", ids[0]), zap.Int("height", int(currentHeight)))

			handledIds := ids[0:len(res.NextBlocks)]

			for idx, _ := range handledIds {
				// ToDo: add BlockID calc here
				// if calculatedBlockId of res.NextBlocks[idx] != id {
				//         continue
				// }
				currentHeight++
				// handled the requested range - return
				if blockRange.to.height != 0 && blockRange.to.id != 0 && currentHeight > blockRange.to.height {
					goto DONE
				}
				synchronizer.store.RawStore.Store(res.NextBlocks[idx], currentHeight)
			}
			// looks like we did not get data for all blocks we wanted to get
			// this could happen to a response size limit in the old java wallet
			if len(handledIds) != len(ids) {
				ids = ids[len(res.NextBlocks)-1 : len(ids)]
				continue
			}

			break
		}
	DONE:
	}

	synchronizer.wg.Done()
}
