package redis

import "github.com/go-redis/redis"

func NewClient() redis.UniversalClient {
	config := &redis.ClusterOptions{
		Addrs:        []string{
			"localhost:7000",
			"localhost:7001",
			"localhost:7002",
		},
		MaxRetries:   1,
		PoolSize:     10,
		MinIdleConns: 5,
	}

	return NewProxyClient(config)

	//return redis.NewClusterClient(config)
}