package watcher

import (
	"fmt"

	"github.com/mkramb/mongodb-nats-connector/internal/mongo"
	"github.com/mkramb/mongodb-nats-connector/internal/nats"
)

func (s *Server) watchForChangeEvents() {
	s.MongoClient.StartWatcher()
	s.MongoClient.OnChangeEvent(func(event *mongo.ChangeEvent, json []byte) {
		var opts nats.PublishOptions

		opts.MsgId = event.ResumeToken.Value
		opts.Subject = fmt.Sprintf("%v.%v.%v", event.Ns.Coll, event.OperationType, event.FullDocument.Id.Value)
		opts.Data = json

		s.NatsClient.PublishEvent(&opts)
	})
}
