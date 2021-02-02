package main

import (
	"context"
	"github.com/hcolde/reviewer-helper/control"
	db "github.com/hcolde/reviewer-helper/database"
	"github.com/hcolde/reviewer-helper/log"
	"github.com/hcolde/reviewer-helper/pay"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serverNum := 2
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

	quit := make(chan bool, serverNum)
	pubQuit := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())

	pubMsg := make(chan db.PublishMsg, 1024)

	go control.Publisher(ctx, pubMsg, pubQuit) // 发送消息
	go pay.Pay(ctx, pubMsg, quit)              // 转账
	go pay.VIP(ctx, pubMsg, quit)              // 购买会员

	log.Logger.Info("service stared")

	<-ch
	cancel()
	for i := 0; i < serverNum; i++ {
		<-quit
	}
	close(pubMsg)
	<-pubQuit
	log.Logger.Info("service exited")
}
