package main

import (
	//	pb "github.com/ac0v/ce-stock/internal/api/protobuf-spec"
	p2p "github.com/ac0v/aspera/pkg/p2p"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	"go.uber.org/zap"
	//"github.com/dgraph-io/badger"
	//"github.com/golang/protobuf/proto"
	// "gopkg.in/resty.v1"
	//"encoding/binary"
	//"log"
	//"path/filepath"
	//"fmt"
	"strconv"
)

func main() {
	r.Init()
	registry := &r.Context
	client := p2p.NewClient(registry)
	store := s.Init(registry)
	defer store.Close()

	for {
		latestBlock := store.RawStore.Current.Block
		registry.Logger.Info("syncing", zap.Int("height", int(store.RawStore.Current.Height)), zap.String("previousBlock", latestBlock.PreviousBlock))

		nextBlockId, _ := strconv.ParseUint(latestBlock.PreviousBlock, 10, 64)
		if latestBlock.Height == 0 {
			nextBlockId = uint64(latestBlock.Block)
		}
		res, _ := client.GetNextBlocksByMajority(nextBlockId)

		if len(res.NextBlocks) == 0 {
			break
		}
		if latestBlock.Height != 0 {
			// atm we do not know the blockId, but it's previous
			// - so we ignore the double returned block
			res.NextBlocks = res.NextBlocks[1:]
		}
		for _, block := range res.NextBlocks {
			store.RawStore.Push(block)
			latestBlock = block
		}
	}
}
