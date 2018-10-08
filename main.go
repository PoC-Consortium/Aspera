package main

import (
	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"
	p2p "github.com/ac0v/aspera/pkg/p2p"
	s "github.com/ac0v/aspera/pkg/store"
	"go.uber.org/zap"
)

func main() {
	c, err := config.Parse("config.yml")
	if err != nil {
		Log.Fatal("parse config", zap.Error(err))
	}

	client := p2p.NewClient(&c.Network.P2P, c.Network.InternetProtocols)
	store := s.Init(c.Storage.Path, c.Network.P2P.Milestones[0])
	defer store.Close()

	p2p.NewSynchronizer(client, store, c.Network.P2P.Milestones)
}
