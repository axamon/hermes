package inoltralog

import (
	"context"

	"github.com/go-redis/redis"
)

// TestRemoteRedisServer verifica la raggiungibilità del server Redis remoto.
func TestRemoteRedisServer(ctx context.Context, remoteRedisServer, remoteRedisServerPass string) (ok bool, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     remoteRedisServer,
		Password: remoteRedisServerPass, // no password set
		DB:       0,                     // use default DB
	})
	defer client.Close()

	pong, err := client.Ping().Result()
	//fmt.Println(pong, err)
	if pong == "pong" {
		return true, nil
	}
	return false, err
}
