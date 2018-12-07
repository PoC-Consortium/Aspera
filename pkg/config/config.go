package config

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	. "github.com/PoC-Consortium/Aspera/pkg/log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrInvalidInternetProtocol = errors.New("invalid internet protocol")
	ErrEmptyPeers              = errors.New("no peers")
)

const (
	Application    = "Aspera"
	Version        = "0.0.1"
	DefaultP2PPort = "8123"
)

type Config struct {
	Common  Common
	Network Network
	Storage Storage
	Log     zap.Config
}

type Common struct {
	Platform string
}

type Network struct {
	InternetProtocols []string
	P2P               P2P
}

type P2P struct {
	Timeout    time.Duration
	Debug      bool
	Listen     string
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

	c.Log = zap.NewProductionConfig()
	c.Log.OutputPaths = append(c.Log.OutputPaths, "var/log/aspera.log")
	c.Network.P2P.Timeout = 5 * time.Second
	c.Network.P2P.Debug = false

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	for _, logfile := range c.Log.OutputPaths {
		if logfile != "stdout" && logfile != "stderr" {
			if _, err := os.Stat(filepath.Dir(logfile)); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(logfile), os.ModePerm)
			}
		}
	}

	var err error
	if Log, err = c.Log.Build(); err != nil {
		return nil, err
	}
	// .. let's wait for the next badger release ...
	// badger.SetLogger(Log.Sugar())

	if err := validate(&c); err != nil {
		return nil, err
	}

	maybeInsertGenesisMilestone(&c)

	if len(c.Common.Platform) < 1 {
		c.Common.Platform = runtime.GOOS
	}

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

func parseIP(c *Config, s string) (net.IP, string, string, error) {
	var ip net.IP
	var port string
	var err error
	if ip = net.ParseIP(s); ip == nil {
		var host string
		if host, port, err = net.SplitHostPort(s); err != nil {
			if len(port) > 0 {
				if _, err = strconv.ParseUint(port, 10, 16); err != nil {
					return nil, "", "", err
				}
			}
			if ip = net.ParseIP(host); ip == nil {
				return nil, "", "", errors.New("invalid IP")
			}
		} else {
			return nil, "", "", err
		}
	}

	space := "IPv6"
	if ip4 := ip.To4(); ip4 != nil {
		space = "IPv4"
		ip = ip4
	}

	return ip, port, space, nil
}

/*
func detectPublicIP(c *Config) error {
	if len(c.Network.P2P.Listen) == 0 {

	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
http: //whatismyip.akamai.com/
}
*/
