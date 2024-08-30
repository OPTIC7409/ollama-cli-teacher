package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CacheManager struct {
	filePath string
	cache    map[string]string
}

func NewCacheManager(filePath string) *CacheManager {
	manager := &CacheManager{
		filePath: filePath,
		cache:    make(map[string]string),
	}
	manager.loadCache()
	return manager
}

func (c *CacheManager) loadCache() {
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		return
	}

	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		fmt.Printf("Error reading cache file: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &c.cache)
	if err != nil {
		fmt.Printf("Error unmarshalling cache data: %v\n", err)
	}
}

func (c *CacheManager) saveCache() {
	data, err := json.Marshal(c.cache)
	if err != nil {
		fmt.Printf("Error marshalling cache data: %v\n", err)
		return
	}

	err = ioutil.WriteFile(c.filePath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
	}
}

func (c *CacheManager) HasTopic(topic string) bool {
	_, exists := c.cache[topic]
	return exists
}

func (c *CacheManager) GetCategories(topic string) (string, error) {
	categories, exists := c.cache[topic]
	if !exists {
		return "", fmt.Errorf("categories not found for topic '%s'", topic)
	}
	return categories, nil
}

func (c *CacheManager) SaveCategories(topic, categories string) {
	c.cache[topic] = categories
	c.saveCache()
}
