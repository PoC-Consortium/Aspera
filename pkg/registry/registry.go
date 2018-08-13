package registry

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	StoragePath string `yaml:"storagePath"`
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
}
