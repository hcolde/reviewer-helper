package pay

import (
	"context"
	"encoding/json"
	"github.com/hcolde/reviewer-helper/conf"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
	"strconv"
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

	payInfo := db.PayInfo{}
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
		return
	}

	tx := mysql.DB.Exec(`INSERT INTO members (wxid, money) VALUES (?, ?) ON DUPLICATE KEY UPDATE money=money+?`,
		wxid,
		payInfo.Money,
		payInfo.Money,
	)

	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return
	}

	moneyStr := ""
	if err := redis.Hget(ctx, "member", wxid, &moneyStr); err != nil {
		log.Logger.Error(err)
		return
	}

	var money float64
	if moneyStr != "" {
		m, err := strconv.ParseFloat(moneyStr, 64)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		money = m
	}
	money += payInfo.Money

	moneyStr = strconv.FormatFloat(money, 'f', 2, 64)
	if err := redis.Hset(ctx, "member", wxid, moneyStr); err != nil {
		log.Logger.Error(err)
	}
}

func Pay(ctx context.Context) {
	if err := mysql.New(); err != nil {
		log.Logger.Fatal(err)
	}

	if err := redis.New(ctx); err != nil {
		log.Logger.Fatal(err)
	}
	defer redis.Close()

	if err := mysql.DB.AutoMigrate(db.Member{}); err != nil {
		log.Logger.Fatal(err)
	}

	if err := mysql.DB.AutoMigrate(db.PayInfo{}); err != nil {
		log.Logger.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var data []string
			if err := redis.BRPop(ctx, conf.Conf.Redis.Key.PayList, &data); err != nil {
				log.Logger.Error(err)
			}

			go operate(ctx, data)
		}
	}
}
