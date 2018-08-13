package main

import (
	//	pb "github.com/ac0v/ce-stock/internal/api/protobuf-spec"
	"github.com/ac0v/aspera/pkg/p2p"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	//"github.com/dgraph-io/badger"
	//"github.com/golang/protobuf/proto"
	// "gopkg.in/resty.v1"
	//"encoding/binary"
	//"log"
	//"path/filepath"
	//"strconv"
)

func main() {
	r.Init()
	registry := &r.Context
	store := s.Init(registry)
	defer store.Close()

	latestBlock := store.RawStore.Current.Block
	for {
		res, _ := p2p.GetNextBlocks(uint64(latestBlock.Block))
		if len(res.NextBlocks) == 0 {
			break
		}
		for _, block := range res.NextBlocks {
			store.RawStore.Push(block)
			latestBlock = block
		}
	}

}
