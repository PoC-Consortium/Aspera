package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

var (
	ErrInvalidInternetProtocol = errors.New("invalid internet protocol")
	ErrEmptyPeers              = errors.New("no peers")
)

type Config struct {
	Network Network
	Storage Storage
}

type Network struct {
	InternetProtocols []string
	P2P               P2P
}

type P2P struct {
	Timeout    time.Duration
	Debug      bool
	Peers      []string
	Milestones []Milestone
}

type Milestone struct {
	Height int32
	Id     uint64
}

type Storage struct {
	Path string
}

func Parse(cfgPath string) (*Config, error) {
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	c.Network.P2P.Timeout = 5 * time.Second
	c.Network.P2P.Debug = false

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	if err := validate(&c); err != nil {
		return nil, err
	}

	maybeInsertGenesisMilestone(&c)

	return &c, nil
}

func maybeInsertGenesisMilestone(c *Config) {
	milestones := c.Network.P2P.Milestones
	if len(milestones) == 0 || milestones[0].Height != 0 {
		c.Network.P2P.Milestones = append([]Milestone{
			Milestone{
				Height: 0,
				Id:     3444294670862540038,
			},
		}, milestones...)
	}
}

func validate(c *Config) error {
	if len(c.Network.P2P.Peers) == 0 {
		return ErrEmptyPeers
	}

	protocols := c.Network.InternetProtocols
	if len(protocols) > 2 {
		return ErrInvalidInternetProtocol
	}
	for _, protocol := range protocols {
		switch protocol {
		case "v4":
		case "v6":
		default:
			return ErrInvalidInternetProtocol
		}
	}
	if len(protocols) == 0 {
		c.Network.InternetProtocols = []string{"v6", "v4"}
	}

	return nil
}
