package control

import (
	"context"
	"encoding/json"
	"github.com/hcolde/reviewer-helper/conf"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
)

var redis db.Redis

func init() {
	if err := redis.New(context.Background()); err != nil {
		log.Logger.Fatal(err)
	}
}

func Publisher(ctx context.Context, pubMsgPay chan db.PublishMsg, quit chan bool) {
	defer func() {
		log.Logger.Info("publisher exited")
		quit <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case pm, ok := <-pubMsgPay:
			if !ok {return}

			msg, err := json.Marshal(pm)
			if err != nil {
				log.Logger.Error(err)
				continue
			}
			if err := redis.LPush(ctx, conf.Conf.Redis.Key.Publisher, string(msg)); err != nil {
				log.Logger.Error(err)
			} else {
				log.Logger.Info(string(msg))
			}
		}
	}
}
