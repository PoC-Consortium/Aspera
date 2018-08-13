package store

import (
	"fmt"
	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type RawStore struct {
	BasePath string
	registry *r.Registry
	Current  *RawCurrent
}

type RawCurrent struct {
	Height int32
	Block  *pb.Block
}

func NewRawStore(registry *r.Registry) *RawStore {
	var rawStore RawStore
	rawStore.registry = registry

	rawStore.BasePath = filepath.Join(registry.Config.StoragePath, "raw")
	if _, err := os.Stat(rawStore.BasePath); os.IsNotExist(err) {
		os.MkdirAll(rawStore.BasePath, os.ModePerm)
	}

	// get most recent block
	heightString := "0"
	numericRegexp := regexp.MustCompile("^[0-9]+")

	currentPath := rawStore.BasePath
	for {
		items, err := ioutil.ReadDir(currentPath)
		if err != nil {
			rawStore.registry.Logger.Fatal("Fatal", zap.Error(err))
		}
		if len(items) == 0 {
			break
		}

		currentItem := items[len(items)-1]
		for _, item := range items {
			if item.IsDir() {
				currentItem = item
			}
		}

		currentPath = filepath.Join(currentPath, currentItem.Name())
		heightString += strings.Join(numericRegexp.FindAllString(currentItem.Name(), -1), "")
		if !currentItem.IsDir() {
			break
		}
	}

	// update current; create genesis or load most recent block
	height, err := strconv.ParseInt(heightString, 10, 32)
	if err != nil {
		rawStore.registry.Logger.Fatal("Fatal", zap.Error(err))
	}
	rawStore.Current = &RawCurrent{Height: int32(height)}

	if heightString == "0" {
		block := &pb.Block{Block: 3444294670862540038}
		rawStore.store(block, 0)
	} else {
		rawStore.Current.Block = rawStore.load(rawStore.Current.Height)
	}
	rawStore.registry.Logger.Info("loaded Raw Storage", zap.Int("height", int(height)))

	return &rawStore
}

func (rawStore *RawStore) Push(block *pb.Block) {
	rawStore.store(block, rawStore.Current.Height+1)
}

func (rawStore *RawStore) convertHeightToPathInfo(height int32) string {
	parts := []rune(fmt.Sprintf("%010d", int(height)))
	var path string
	for _, part := range parts {
		path = filepath.Join(path, string(part))
	}
	return filepath.Join(rawStore.BasePath, path+".bin")
}

func (rawStore *RawStore) store(block *pb.Block, height int32) {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	block.Height = height
	/* ToDo:
	block.PayloadLength
	block.TotalAmountNQT
	block.Generator
	block.BaseTarget
	block.generatorRS
	block.BlockReward
	block.NextBlock
	block.ScoopNum
	block.NumberOfTransactions
	block.Transactions
	block.TotalFeeNQT
	block.Block
	*/

	data, _ := proto.Marshal(block)
	ioutil.WriteFile(path, data, os.ModePerm)

	rawStore.Current.Height = height
	rawStore.Current.Block = block
}

func (rawStore *RawStore) load(height int32) *pb.Block {
	path := rawStore.convertHeightToPathInfo(height)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	in, err := ioutil.ReadFile(path)
	if err != nil {
		rawStore.registry.Logger.Fatal("Error reading file:", zap.Error(err))
	}
	block := &pb.Block{}
	if err := proto.Unmarshal(in, block); err != nil {
		rawStore.registry.Logger.Fatal("Error parse block file:", zap.Error(err))
	}

	return block
}
