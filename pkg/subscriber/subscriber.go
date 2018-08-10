package subscriber

import (
	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"go.uber.org/zap"
	"reflect"
	"runtime"
)

type fn func()

func Subscribe(db *badger.DB, subject string) { // , unmarshal fn) {
	logger, _ := zap.NewProduction()

	connectionLogger := logger.With(zap.String("url", nats.DefaultURL), zap.String("subject", subject))
	connectionLogger.Debug("Subscribing to Queue")

	natsConnection, err := nats.Connect(nats.DefaultURL)
	if err == nil {
		natsConnection.Subscribe(subject, func(msg *nats.Msg) {
			connectionLogger.Debug("received message")

			parsedMessage := defaultMessage
			err = proto.Unmarshal(msg.Data, parsedMessage)

			if err == nil {
				cidBytes := reflect.Indirect(reflect.ValueOf(parsedMessage)).FieldByName("cid").Bytes()
				cid := string(cidBytes)
				cidLogger := logger.With(zap.String("cid", cid))

				cidLogger.Debug("message unmarshaled")
				err = db.Update(func(txn *badger.Txn) error {
					err := txn.Set(cidBytes, msg.Data)
					return err
				})
				if err == nil {
					cidLogger.Debug("inserted message")
				} else {
					cidLogger.Error("can't unmarshal message", zap.ByteString("messageData", msg.Data), zap.Error(err))
				}
			} else {
				connectionLogger.Error("can't unmarshal message", zap.ByteString("messageData", msg.Data), zap.Error(err))
			}
		})

		// Keep the connection alive
		runtime.Goexit()
	} else {
		connectionLogger.Fatal("can't connect to queue server", zap.Error(err))
	}
}
