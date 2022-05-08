package cfg

import (
	"context"
	"fmt"

	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/radix/v4"
)

// RadixClient is a single redis client which can be configured.
type RadixClient struct {
	radix.Client

	proto, addr string
	poolSize    int
}

// SetupCfg implement the cfg.Cfger interface.
func (c *RadixClient) SetupCfg(cfg *Cfg) {

	cfg.StringVar(&c.proto, "redis-proto", "tcp", "Network protocol to connect to redis over, can be tcp or unix")
	cfg.StringVar(&c.addr, "redis-addr", "127.0.0.1:6379", "Address redis is expected to listen on")
	cfg.IntVar(&c.poolSize, "redis-pool-size", 5, "Number of connections in the redis pool to keep")

	cfg.OnInit(func(ctx context.Context) error {
		client, err := (radix.PoolConfig{
			Size: c.poolSize,
		}).New(
			ctx, c.proto, c.addr,
		)

		if err != nil {
			return fmt.Errorf(
				"initializing redis pool of size %d at %s://%s: %w",
				c.poolSize, c.proto, c.addr, err,
			)
		}

		c.Client = client
		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (c *RadixClient) Annotate(a mctx.Annotations) {
	a["redisProto"] = c.proto
	a["redisAddr"] = c.addr
	a["redisPoolSize"] = c.poolSize
}

// Close cleans up the radix client.
func (c *RadixClient) Close() error {
	return c.Client.Close()
}
