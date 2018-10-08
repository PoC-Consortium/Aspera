package registry

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	. "github.com/ac0v/aspera/pkg/log"
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
	Timeout    int         `yaml:"timeout"`
	Debug      bool        `yaml:"debug"`
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
	Config Config
}

var Context Registry

func Init() {
	raw, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		Log.Fatal("reading config failed", zap.Error(err))
	}

	err = yaml.Unmarshal(raw, &Context.Config)
	if err != nil {
		Log.Fatal("unpacking config failed", zap.Error(err))
	}

	if len(Context.Config.Network.InternetProtocols) > 2 {
		Log.Fatal("invalid amount of internetProtocols")
	} else if len(Context.Config.Network.InternetProtocols) > 0 {
		for _, protocol := range Context.Config.Network.InternetProtocols {
			if protocol != "v4" && protocol != "v6" {
				Log.Fatal("invalid internetProtocol", zap.String("protocol", protocol))
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

	if Context.Config.Network.P2P.Timeout == 0 {
		Context.Config.Network.P2P.Timeout = 5
	}
}
