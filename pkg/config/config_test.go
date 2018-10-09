package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c, err := Parse("config.yml")

	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, []string{"v4"}, c.Network.InternetProtocols)
	assert.Equal(t, []string{"wallet.burst.cryptoguru.org", "burst-a.bayernport.com"}, c.Network.P2P.Peers)
	assert.Equal(t, 5*time.Second, c.Network.P2P.Timeout)
	assert.Equal(t, []Milestone{
		Milestone{
			Height: 0,
			Id:     3444294670862540038,
		},
		Milestone{
			Height: 50000,
			Id:     17169998969130562818,
		},
		Milestone{
			Height: 100000,
			Id:     10851012679396814781,
		},
	}, c.Network.P2P.Milestones)
}
