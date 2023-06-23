package cache

import (
	"context"
	"sql_generate/global"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/23
 * Time: 13:16
 * Description: 缓存
 */

// GetKV 根据获取值
func (c *Cache) GetKV(ctx context.Context, key string) (string, error) {
	result, err := global.CaChe.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

type Cache struct{}

// SetKV 设置键值存入 redis
func (c *Cache) SetKV(ctx context.Context, key string, data []byte, expire time.Duration) error {
	_, err := global.CaChe.Set(ctx, key, data, expire).Result()
	if err != nil {
		return err
	}
	return nil
}

// DeleteKV 删除键
func (c *Cache) DeleteKV(ctx context.Context, keyPattern string) error {
	// Find keys matching the pattern
	keys, err := global.CaChe.Keys(ctx, keyPattern).Result()
	if err != nil {
		return err
	}

	// Ensure there are keys to delete
	if len(keys) == 0 {
		return nil
	}

	// Delete the keys
	_, err = global.CaChe.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}

	return nil
}
