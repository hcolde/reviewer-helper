package pay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hcolde/reviewer-helper/conf"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
	"strconv"
)

var (
	mysql = db.Mysql{}
	redis = db.Redis{}
)

func init() {
	if err := mysql.New(); err != nil {
		log.Logger.Fatal(err)
	}

	if err := redis.New(context.Background()); err != nil {
		log.Logger.Fatal(err)
	}
}

func getMoney(ctx context.Context, wxid string) (money float64) {
	moneyStr := ""
	if err := redis.Hget(ctx, conf.Conf.Redis.Key.Member, wxid, &moneyStr); err != nil {
		log.Logger.Error(err)
		return 0
	}

	if moneyStr != "" {
		m, err := strconv.ParseFloat(moneyStr, 64)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		money = m
	}
	return
}

func operate(ctx context.Context, data []string, pubMsg chan db.PublishMsg) {
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
	if err := redis.Hget(ctx, conf.Conf.Redis.Key.Member, wxid, &moneyStr); err != nil {
		log.Logger.Error(err)
		return
	}

	money := getMoney(ctx, wxid)
	money += payInfo.Money

	moneyStr = strconv.FormatFloat(money, 'f', 2, 64)
	if err := redis.Hset(ctx, conf.Conf.Redis.Key.Member, wxid, moneyStr); err != nil {
		log.Logger.Error(err)
	} else {
		pubMsg <- db.PublishMsg{
			Wxid: payInfo.Wxid,
			Msg:  fmt.Sprintf("收到转账：%.2f元[好的]", payInfo.Money),
		}
	}
}

func Pay(ctx context.Context, pubMsg chan db.PublishMsg, quit chan bool) {
	defer func() {
		log.Logger.Info("pay exited")
		quit <- true
	}()

	if err := mysql.DB.AutoMigrate(db.Member{}); err != nil {
		log.Logger.Fatal(err)
	}

	if err := mysql.DB.AutoMigrate(db.PayInfo{}); err != nil {
		log.Logger.Fatal(err)
	}

	for {
		var data []string
		select {
		case <-ctx.Done():
			return
		default:
			if err := redis.BRPop(ctx, conf.Conf.Redis.Key.PayList, &data); err != nil {
				log.Logger.Error(err)
				continue
			}

			go operate(ctx, data, pubMsg)
		}
	}
}
