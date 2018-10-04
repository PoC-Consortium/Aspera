package main

import (
	//	pb "github.com/ac0v/ce-stock/internal/api/protobuf-spec"
	p2p "github.com/ac0v/aspera/pkg/p2p"
	r "github.com/ac0v/aspera/pkg/registry"
	s "github.com/ac0v/aspera/pkg/store"
	//"go.uber.org/zap"
	//"github.com/dgraph-io/badger"
	//"github.com/golang/protobuf/proto"
	// "gopkg.in/resty.v1"
	//"encoding/binary"
	//"log"
	//"path/filepath"
	//"fmt"
	//"strconv"
)

func main() {
	r.Init()
	registry := &r.Context
	client := p2p.NewClient(registry)
	store := s.Init(registry)
	defer store.Close()

	p2p.NewSynchronizer(client, store, registry)
}
