package cache

import (
	"encoding/json"
	"os"
	"time"
)

type CacheEntry struct {
	Value      map[string]map[string]string
	Expiration time.Time
}

type Cache struct {
	Items map[string]CacheEntry
}

func NewCache() *Cache {
	return &Cache{
		Items: make(map[string]CacheEntry),
	}
}

func (c *Cache) Set(key string, value map[string]map[string]string, expiration time.Duration) {
	expirationTime := time.Now().Add(expiration)
	c.Items[key] = CacheEntry{Value: value, Expiration: expirationTime}
}

func (c *Cache) Get(key string) (map[string]map[string]string, bool) {
	entry, exists := c.Items[key]
	if !exists || time.Now().After(entry.Expiration) {
		return nil, false // Cache miss or expired
	}

	return entry.Value, true
}

func (c *Cache) SaveToFile(filename string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *Cache) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
