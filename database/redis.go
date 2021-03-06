package database

import (
	"context"
	"github.com/hcolde/reviewer-helper/conf"
	redix "github.com/mediocregopher/radix/v4"
)

type Redis struct {
	client redix.Client
}

func (redis *Redis) New(ctx context.Context) error {
	client, err := (redix.PoolConfig{}).New(ctx, "tcp", conf.Conf.Redis.Host)
	if err != nil {
		return err
	}

	if conf.Conf.Redis.Password != "" {
		if err := client.Do(ctx, redix.Cmd(nil, "AUTH", conf.Conf.Redis.Password)); err != nil {
			return err
		}
	}

	redis.client = client
	return nil
}

func (redis *Redis) LPush(ctx context.Context, key, value string) error {
	return redis.client.Do(ctx, redix.Cmd(nil, "LPUSH", key, value))
}

func (redis *Redis) BRPop(ctx context.Context, key string, rcv *[]string) error {
	return redis.client.Do(ctx, redix.Cmd(&rcv, "BRPOP", key, "0"))
}

func (redis *Redis) Get(ctx context.Context, key string, rcv *string) error {
	return redis.client.Do(ctx, redix.Cmd(&rcv, "GET", key))
}

func (redis *Redis) Hset(ctx context.Context, key, field, value string) error {
	return redis.client.Do(ctx, redix.Cmd(nil, "HSET", key, field, value))
}

func (redis *Redis) Hget(ctx context.Context, key, field string, rcv *string) error {
	return redis.client.Do(ctx, redix.Cmd(&rcv, "HGET", key, field))
}

func (redis *Redis) Del(ctx context.Context, key string) error {
	return redis.client.Do(ctx, redix.Cmd(nil, "DEL", key))
}

func (redis *Redis) Close() {
	_ = redis.client.Close()
}
