package inoltralog

import (
	"context"

	"github.com/go-redis/redis"
)

// REDIS
var remoteRedisServer = "easyapi.westeurope.cloudapp.azure.com:6379"
var remoteRedisServerPass = "pippo"

// TestRemoteRedisServer verifica la raggiungibilità del server Redis remoto.
func TestRemoteRedisServer(ctx context.Context) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     remoteRedisServer,
		Password: remoteRedisServerPass, // no password set
		DB:       0,                     // use default DB
	})
	defer client.Close()

	_, err = client.Ping().Result()
	//fmt.Println(pong, err)

	return err
}

// InviaRecordRedisRemoto invia un record al Redis Server Remoto.
func InviaRecordRedisRemoto(ctx context.Context, record string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     remoteRedisServer,
		Password: remoteRedisServerPass, // no password set
		DB:       0,                     // use default DB
	})
	defer client.Close()

	_, err = client.LPush("records", record).Result()

	return err
}
