package databases

import (
	"context"
	"data_recover_14_nov/globals"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func CheckConnection(ctx context.Context, wg *sync.WaitGroup) bool {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Connection checking stopped due to shutdown signal.")
			if RedisClient != nil {
				if err := RedisClient.Close(); err != nil {
					log.Printf("Error closing Redis connection: %v", err)
				} else {
					fmt.Println("Redis connection closed successfully.")
				}
			}
			return false
		default:
			if RedisClient == nil {
				if err := EstablishRedisQueueConnection(); err != nil {
					log.Printf("Error establishing Redis connection: %v", err)
					return false
				}
			}

			if err := RedisClient.Ping(ctx).Err(); err != nil {
				log.Printf("Error pinging Redis: %v", err)
				return false
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func EstablishRedisQueueConnection() error {
	fmt.Println("Establishing Redis Queue Connection")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v", globals.ApplicationConfig.RedisConn.Host+":"+globals.ApplicationConfig.RedisConn.Port),
		Username: globals.ApplicationConfig.RedisConn.Username,
		Password: globals.ApplicationConfig.RedisConn.Password,
		DB:       0,
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("unable to connect to Redis: %v", err)
	}

	RedisClient = rdb
	fmt.Println("Redis Queue Connection Established....")
	return nil
}
