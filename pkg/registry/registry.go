package registry

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Network Network `yaml:"network"`
	Storage Storage `yaml:"storage"`
}

type Network struct {
	InternetProtocols []string `yaml:"internetProtocols"`
	P2P               P2P      `yaml:"p2p"`
}

type P2P struct {
	Peers      []string    `yaml:"peers"`
	Milestones []Milestone `yaml:"milestones"`
}

type Milestone struct {
	Height int32  `yaml:"height"`
	Id     uint64 `yaml:"id"`
}

type Storage struct {
	Path string `yaml:"path"`
}

type Registry struct {
	Logger zap.Logger
	Config Config
}

var Context Registry

func Init() {
	logger, _ := zap.NewProduction()
	Context.Logger = *logger

	raw, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		logger.Fatal("reading config failed", zap.Error(err))
	}

	err = yaml.Unmarshal(raw, &Context.Config)
	if err != nil {
		logger.Fatal("unpacking config failed", zap.Error(err))
	}

	if len(Context.Config.Network.InternetProtocols) > 2 {
		logger.Fatal("invalid amount of internetProtocols")
	} else if len(Context.Config.Network.InternetProtocols) > 0 {
		for _, protocol := range Context.Config.Network.InternetProtocols {
			if protocol != "v4" && protocol != "v6" {
				logger.Fatal("invalid internetProtocol", zap.String("protocol", protocol))
			}
		}
	} else {
		Context.Config.Network.InternetProtocols = []string{"v4", "v6"}
	}

	// our genesis block should be always the first milestone
	// - that's also important for our current raw storage implementation
	if len(Context.Config.Network.P2P.Milestones) == 0 || Context.Config.Network.P2P.Milestones[0].Height != 0 {
		Context.Config.Network.P2P.Milestones = append(
			[]Milestone{
				Milestone{
					Height: 0,
					Id:     3444294670862540038,
				},
			},
			Context.Config.Network.P2P.Milestones...,
		)
		// ToDo: sort milestones by height !!
	}
}
