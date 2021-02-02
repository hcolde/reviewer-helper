package pay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hcolde/reviewer-helper/conf"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
	"strconv"
	"time"
)

func joinVIP(ctx context.Context, data []string, pubMsg chan db.PublishMsg) {
	if len(data) < 2 {
		log.Logger.Info(data)
		return
	}

	vipInfo := db.VipInfo{}
	if err := json.Unmarshal([]byte(data[1]), &vipInfo); err != nil {
		log.Logger.Info(data)
		log.Logger.Error(err)
		return
	}

	msg := db.PublishMsg{Wxid: vipInfo.Wxid}

	money := getMoney(ctx, vipInfo.Wxid)
	if money < vipInfo.Money {
		msg.Msg = "抱歉，您的余额不足，请充值[嘿哈]"
		pubMsg <- msg
		return
	}

	if err := mysql.DB.Create(&vipInfo).Error; err != nil {
		log.Logger.Error(err)
		return
	}

	money = money - vipInfo.Money
	member := db.Member{
		Wxid:  vipInfo.Wxid,
		Money: money,
	}

	var day int64
	if vipInfo.Money == 1 {
		day = 1
	} else if vipInfo.Money == 25 {
		day = 31
	} else if vipInfo.Money == 60 {
		day = 93
	} else if vipInfo.Money == 200 {
		day = 372
	}

	vipStr := "0"
	if err := redis.Hget(ctx, conf.Conf.Redis.Key.VIP, vipInfo.Wxid, &vipStr); err != nil {
		log.Logger.Error(err)
		return
	}
	vip, err := strconv.ParseInt(vipStr, 10, 64)
	if err != nil {
		log.Logger.Error(err)
		return
	}

	loc := time.FixedZone("UTC", 8 * 3600) // 东八区
	now := time.Now().In(loc).Unix()

	if vip < now {
		member.Vip = now + day * 24 * 60 * 60
	} else {
		member.Vip = vip + day * 24 * 60 * 60
	}

	if err := mysql.DB.Save(&member).Error; err != nil {
		log.Logger.Error(err)
		return
	}

	moneyStr := strconv.FormatFloat(money, 'f', 2, 64)
	if err := redis.Hset(ctx, conf.Conf.Redis.Key.Member, vipInfo.Wxid, moneyStr); err != nil {
		log.Logger.Error(err)
		return
	}

	vipStr = strconv.FormatInt(member.Vip, 10)

	if err := redis.Hset(ctx, conf.Conf.Redis.Key.VIP, vipInfo.Wxid, vipStr); err != nil {
		log.Logger.Error(err)
		return
	}

	end := time.Unix(member.Vip, 0).In(loc).Format("2006年01月02日 15:04:05")
	msg.Msg = fmt.Sprintf("[庆祝]您已成功开通%d天会员[庆祝]\n会员到期时间为：\n%s", day, end)

	if err := redis.Del(ctx, fmt.Sprintf("%s_vip", vipInfo.Wxid)); err != nil {
		log.Logger.Error(err)
	}

	pubMsg <- msg
}

func VIP(ctx context.Context, pubMsg chan db.PublishMsg, quit chan bool) {
	defer func() {
		log.Logger.Info("vip exited")
		quit <- true
	}()

	if err := mysql.DB.AutoMigrate(db.VipInfo{}); err != nil {
		log.Logger.Fatal(err)
	}

	for {
		var data []string
		select {
		case <-ctx.Done():
			return
		default:
			if err := redis.BRPop(ctx, conf.Conf.Redis.Key.VIPList, &data); err != nil {
				log.Logger.Error(err)
				continue
			}

			go joinVIP(ctx, data, pubMsg)
		}
	}
}
