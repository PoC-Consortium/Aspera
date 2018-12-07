package main

import (
	"flag"

	"github.com/PoC-Consortium/Aspera/pkg/blockchain"
	"github.com/PoC-Consortium/Aspera/pkg/config"
	p2p "github.com/PoC-Consortium/Aspera/pkg/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/p2p/manager"
	"github.com/PoC-Consortium/Aspera/pkg/shell"
	s "github.com/PoC-Consortium/Aspera/pkg/store"
)

func main() {
	enablePrompt := flag.Bool("interactive", true, "enable interactive shell prompt")
	flag.Parse()

	blockchain.Init()

	c, err := config.Parse("config.yml")
	if err != nil {
		panic("parse config" + err.Error())
	}

	store := s.Init(c.Storage.Path, c.Network.P2P.Milestones[0])
	defer store.Close()
	p2p.Serve(c, store)

	var minHeights []int32
	for _, m := range c.Network.P2P.Milestones {
		minHeights = append(minHeights, m.Height)
	}
	manager := manager.NewManager(c.Network.P2P.Peers, c.Network.InternetProtocols)
	manager.SetIterators(minHeights)

	client := p2p.NewClient(&c.Network.P2P, manager)

	if *enablePrompt {
		shell.Prompt(store)
	}

	p2p.NewSynchronizer(client, store, c.Network.P2P.Milestones)
}
