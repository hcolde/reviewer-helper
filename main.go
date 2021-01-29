package main

import (
	"context"
	"encoding/json"
	"github.com/hcolde/reviewer-helper/conf"
	"github.com/hcolde/reviewer-helper/database"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
)

var (
	mysql = db.Mysql{}
	redis = db.Redis{}
)

func operate(ctx context.Context, data []string) {
	if len(data) < 2 {
		log.Logger.Info(data)
		return
	}

	payInfo := database.PayInfo{}
	if err := json.Unmarshal([]byte(data[1]), &payInfo); err != nil {
		log.Logger.Info(data)
		log.Logger.Error(err)
		return
	}

	wxid := ""
	if err := redis.Get(ctx, payInfo.TransID, &wxid); err != nil {
		log.Logger.Info("could not get wxid from ", payInfo.TransID)
		log.Logger.Error(err)
		return
	}

	payInfo.Wxid = wxid

	if err := mysql.DB.Create(&payInfo).Error; err != nil {
		log.Logger.Error(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := mysql.New(); err != nil {
		log.Logger.Fatal(err)
	}

	if err := redis.New(ctx); err != nil {
		log.Logger.Fatal(err)
	}
	defer redis.Close()

	if err := mysql.DB.AutoMigrate(database.PayInfo{}); err != nil {
		log.Logger.Fatal(err)
	}

	for {
		var data []string

		if err := redis.BRPop(ctx, conf.Conf.Redis.Key.PayList, &data); err != nil {
			log.Logger.Error(err)
		}

		go operate(ctx, data)
	}
}
