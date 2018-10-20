package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"
	p2p "github.com/ac0v/aspera/pkg/p2p"
	s "github.com/ac0v/aspera/pkg/store"
)

func main() {
	c, err := config.Parse("config.yml")
	if err != nil {
		Log.Fatal("parse config", zap.Error(err))
	}

	client := p2p.NewClient(&c.Network.P2P, c.Network.InternetProtocols)
	manager := p2p.NewManager(client, c.Network.P2P.Peers, c.Network.InternetProtocols, time.Minute)
	client.SetManager(manager)

	store := s.Init(c.Storage.Path, c.Network.P2P.Milestones[0])
	defer store.Close()

	p2p.NewSynchronizer(client, manager, store, c.Network.P2P.Milestones)
}
