package main

import (
	"go.uber.org/zap"

	"github.com/ac0v/aspera/pkg/blockchain"
	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"
	p2p "github.com/ac0v/aspera/pkg/p2p"
	"github.com/ac0v/aspera/pkg/p2p/manager"
	s "github.com/ac0v/aspera/pkg/store"
)

func main() {
	blockchain.Init()

	c, err := config.Parse("config.yml")
	if err != nil {
		Log.Fatal("parse config", zap.Error(err))
	}

	var minHeights []int32
	for _, m := range c.Network.P2P.Milestones {
		minHeights = append(minHeights, m.Height)
	}
	manager := manager.NewManager(c.Network.P2P.Peers, c.Network.InternetProtocols)
	manager.SetIterators(minHeights)

	client := p2p.NewClient(&c.Network.P2P, manager)

	store := s.Init(c.Storage.Path, c.Network.P2P.Milestones[0])
	defer store.Close()

	p2p.NewSynchronizer(client, store, c.Network.P2P.Milestones)
}
