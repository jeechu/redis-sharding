package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/kkdai/consistent"
	"time"
)

func NewProxyClient(options *redis.ClusterOptions) *ProxyClient {
	if options == nil {
		return nil
	}

	proxyClient := &ProxyClient{
		ring:    consistent.NewConsistentHashing(),
		nodeMap: map[string]*redis.Client{},
	}
	for index, addr := range options.Addrs {
		options := redis.Options{
			Addr:               addr,
			Dialer:             options.Dialer,
			OnConnect:          options.OnConnect,
			Username:           options.Username,
			Password:           options.Password,
			DB:                 0,
			MaxRetries:         options.MaxRetries,
			MinRetryBackoff:    options.MinRetryBackoff,
			MaxRetryBackoff:    options.MaxRetryBackoff,
			DialTimeout:        options.DialTimeout,
			ReadTimeout:        options.ReadTimeout,
			WriteTimeout:       options.WriteTimeout,
			PoolSize:           options.PoolSize,
			MinIdleConns:       options.MinIdleConns,
			MaxConnAge:         options.MaxConnAge,
			PoolTimeout:        options.PoolTimeout,
			IdleTimeout:        options.IdleTimeout,
			IdleCheckFrequency: options.IdleCheckFrequency,
			TLSConfig:          options.TLSConfig,
		}
		node := redis.NewClient(&options)
		//proxyClient.nodes = append(proxyClient.nodes, node)
		proxyClient.ring.Add(addr)
		proxyClient.nodeMap[addr] = node
		if index == 0 {
			proxyClient.nodeMap["default"] = node
		}
	}

	//proxyClient.ring = hashring.New(options.Addrs)

	return proxyClient
}

type ProxyClient struct {
	redis.Cmdable
	ring    *consistent.ConsistentHashing
	nodeMap map[string]*redis.Client
}

func (c *ProxyClient) selectNode(keys ...string) *redis.Client {
	// if no key is provided, use default node
	if len(keys) == 0 {
		return c.nodeMap["default"]
	}

	serverAddr, err := c.ring.Get(keys[0])
	if err != nil {
		return c.nodeMap["default"]
	}

	return c.nodeMap[serverAddr]
}

func (c *ProxyClient) Context() context.Context {
	return c.selectNode().Context()
}

func (c *ProxyClient) AddHook(hook redis.Hook) {
	panic("not implemented")
}

func (c *ProxyClient) Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error {
	panic("not implemented")
}

func (c *ProxyClient) Do(ctx context.Context, args ...interface{}) *redis.Cmd {
	panic("not implemented")
}

func (c *ProxyClient) Process(ctx context.Context, cmd redis.Cmder) error {
	panic("not implemented")
}

func (c *ProxyClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	panic("not implemented")
}

func (c *ProxyClient) PSubscribe(ctx context.Context, channels ...string) *redis.PubSub {
	panic("not implemented")
}

func (c *ProxyClient) Close() error {
	var err error
	for _, node := range c.nodeMap {
		nodeErr := node.Close()
		if err == nil {
			err = nodeErr
		}
	}

	return err
}

func (c *ProxyClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.selectNode(key).Get(ctx, key)
}

func (c *ProxyClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.selectNode(key).Set(ctx, key, value, expiration)
}

