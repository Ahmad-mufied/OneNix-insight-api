package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"google-custom-search/model"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedClient struct {
	client *memcache.Client
}

func NewMemcachedClient(host string) *MemcachedClient {
	return &MemcachedClient{
		client: memcache.New(host),
	}
}

func generateCacheKey(filters map[string]string) string {
	hash := sha256.New()
	for k, v := range filters {
		if v == "" {
			// If the value is empty, used 'all' as the default value
			v = "all"

		}
		hash.Write([]byte(k + "=" + v + "&"))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (r *MemcachedClient) GetCachedList(filters map[string]string) ([]model.News, error) {
	cacheKey := generateCacheKey(filters)
	log.Println("Cache key:", cacheKey)
	item, err := r.client.Get(cacheKey)
	if err != nil {
		return nil, err
	}

	var newsList []model.News
	err = json.Unmarshal(item.Value, &newsList)
	if err != nil {
		return nil, err
	}
	return newsList, nil
}

func (r *MemcachedClient) SetCachedList(filters map[string]string, newsList []model.News) error {
	cacheKey := generateCacheKey(filters)
	data, err := json.Marshal(newsList)
	if err != nil {
		return err
	}
	return r.client.Set(&memcache.Item{
		Key:        cacheKey,
		Value:      data,
		Expiration: int32(10 * time.Minute.Seconds()), // Cache for 10 minutes
	})
}
