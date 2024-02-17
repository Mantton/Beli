package cache

import (
	"context"
	"fmt"

	"github.com/mantton/beli/internal/env"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Connect() error {
	url := env.Redis()
	opts, err := redis.ParseURL(url)

	if err != nil {
		return err
	}

	rdb := redis.NewClient(opts)

	c.client = rdb

	return nil
}

func (c *Cache) SetTile(ctx context.Context, x, y int64, color int) error {

	offset := (y * 10) + x
	q := c.client.BitField(
		ctx,                        // Context
		"CURRENT_BOARD",            // Key / Board Name
		"SET",                      // Instruction
		"u8",                       // Variable Type (8 bit Integer) (0 - 256)
		fmt.Sprintf("#%d", offset), // Offset
		color,                      // Value / Data
	)

	_, err := q.Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetTile(ctx context.Context, x, y int64) (int, error) {
	offset := (y * 10) + x

	res, err := c.client.BitField(
		ctx,             // Current Context
		"CURRENT_BOARD", // Board Key
		"GET",           // Command
		"u8",            // Return Type (8 Bit Integer) (0 - 256)
		offset,          // The Tile offset
	).Result()

	if err != nil {
		return 0, err
	}

	return int(res[0]), nil
}

func (c *Cache) GetBoard(ctx context.Context, key string) (string, error) {
	res, err := c.client.Get(ctx, "CURRENT_BOARD").Result()

	if err != nil {
		return "", err
	}

	return res, nil
}
