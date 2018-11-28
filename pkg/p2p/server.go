package p2p

import (
	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/common/math"
	"github.com/ac0v/aspera/pkg/config"
	. "github.com/ac0v/aspera/pkg/log"
	s "github.com/ac0v/aspera/pkg/store"

	"github.com/golang/protobuf/jsonpb"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
)

var (
	marshaler = jsonpb.Marshaler{}
)

func Serve(config *config.P2P, store *s.Store) {
	h := func(ctx *fasthttp.RequestCtx) {
		requestHandler(ctx, store)
	}
	h = fasthttp.CompressHandler(h)
	if err := fasthttp.ListenAndServe(config.Listen, h); err != nil {
		panic(err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx, store *s.Store) {
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
		// {"announcedAddress":"195.201.124.43","application":"BRS","version":"2.2.5","platform":"PC","shareAddress":true}
	case "getMilestoneBlockIds":
		// {"error":"Old getMilestoneBlockIds protocol not supported, please upgrade"}
	case "getNextBlockIds":
		
	case "getBlocksFromHeight":
	case "getNextBlocks":
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
