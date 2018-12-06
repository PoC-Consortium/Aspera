package p2p

import (
	"fmt"
	"strconv"

	api "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/block"
	"github.com/PoC-Consortium/Aspera/pkg/common/math"
	c "github.com/PoC-Consortium/Aspera/pkg/config"
	. "github.com/PoC-Consortium/Aspera/pkg/log"
	s "github.com/PoC-Consortium/Aspera/pkg/store"

	"github.com/golang/protobuf/jsonpb"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
)

var (
	marshaler = jsonpb.Marshaler{}
)

func Serve(config *c.Config, store *s.Store) {
	if len(config.Network.P2P.Listen) < 1 {
		// p2p server disabled
		return
	}
	h := func(ctx *fasthttp.RequestCtx) {
		requestHandler(ctx, config, store)
	}
	h = fasthttp.CompressHandler(h)
	if err := fasthttp.ListenAndServe(config.Network.P2P.Listen, h); err != nil {
		panic(err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx, config *c.Config, store *s.Store) {
	ctx.Response.Header.SetServerBytes([]byte("Aspera"))
	if string(ctx.Path()) != "/burst" {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	var err error
	var src *fastjson.Value
	if src, err = new(fastjson.Parser).ParseBytes(ctx.PostBody()); err != nil {
		Log.Error("bad request", zap.Error(err))
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
	}
	if string(src.GetStringBytes("protocol")) != "B1" {
		Log.Error("invalid protocol", zap.Error(err))
		ctx.Error(`{"error":"Unsupported protocol!"}`, fasthttp.StatusBadRequest)
	}
	switch string(src.GetStringBytes("requestType")) {
	case "addPeers":
		// noop - we do not like to allow peer injection
	case "getCumulativeDifficulty":
		marshaler.Marshal(
			ctx.Response.BodyWriter(),
			&api.GetCumulativeDifficultyResponse{
				BlockchainHeight:     store.RawStore.Current.Height,
				CumulativeDifficulty: math.StringFromBigBytes(store.RawStore.Current.Block.CumulativeDifficulty),
			},
		)
	case "getInfo":
		marshaler.Marshal(
			ctx.Response.BodyWriter(),
			&api.GetInfoResponse{
				AnnouncedAddress: "",
				Version:          c.Version,
				Application:      c.Application,
				Platform:         config.Common.Platform,
				ShareAddress:     true, // only served if listening - and then this is true ... always
			},
		)
	case "getMilestoneBlockIds":
		marshaler.Marshal(
			ctx.Response.BodyWriter(),
			&api.ErrorResponse{
				Error: "Old getMilestoneBlockIds protocol not supported, please upgrade",
			},
		)
	case "getNextBlockIds":
		if blockId, err := strconv.ParseUint(string(src.GetStringBytes("blockId")), 10, 64); err != nil {
			marshaler.Marshal(
				ctx.Response.BodyWriter(),
				&api.ErrorResponse{
					Error: err.Error(),
				},
			)
		} else {
			blocks, _ := store.ChainStore.FindBlocksAfter(
				fmt.Sprintf(block.ById, blockId),
				100,
			)
			var nextBlockIds []uint64
			for _, b := range blocks {
				nextBlockIds = append(nextBlockIds, b.Id)
			}
			marshaler.Marshal(
				ctx.Response.BodyWriter(),
				&api.GetNextBlockIdsResponse{
					NextBlockIds: nextBlockIds,
				},
			)
		}
	case "getBlocksFromHeight":
	case "getNextBlocks":
		if blockId, err := strconv.ParseUint(string(src.GetStringBytes("blockId")), 10, 64); err != nil {
			marshaler.Marshal(
				ctx.Response.BodyWriter(),
				&api.ErrorResponse{
					Error: err.Error(),
				},
			)
		} else {
			blocks, _ := store.ChainStore.FindBlocksAfter(
				fmt.Sprintf(block.ById, blockId),
				100,
			)
			var nextBlocks []*api.Block
			for _, b := range blocks {
				nextBlocks = append(nextBlocks, b.Block)
			}
			marshaler.Marshal(
				ctx.Response.BodyWriter(),
				&api.GetNextBlocksResponse{
					NextBlocks: nextBlocks,
				},
			)
		}
	case "getPeers":
	case "getUnconfirmedTransactions":
	case "processBlock":
	case "processTransactions":
	case "getAccountBalance":
	case "getAccountRecentTransactions":
	default:
		ctx.Error(`{"error":"Unsupported request type!"}`, fasthttp.StatusNotFound)
	}
	//panic(src.Get("requestType").String())
}
