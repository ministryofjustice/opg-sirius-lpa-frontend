package sirius

import (
	"sync"
	"time"
)

type cacheItem struct {
	time  time.Time
	value []RefDataItem
}

var mutex = &sync.RWMutex{}

var cache map[string]cacheItem

func getCached(category string) ([]RefDataItem, bool) {
	var v []RefDataItem
	found := false

	oneHourAgo := time.Now().Add(-1 * time.Hour)

	mutex.RLock()

	if cache[category].time.After(oneHourAgo) && len(cache[category].value) > 0 {
		v = cache[category].value
		found = true
	}

	mutex.RUnlock()

	return v, found
}

func setCached(category string, value []RefDataItem) {
	if cache == nil {
		cache = map[string]cacheItem{}
	}

	mutex.Lock()

	cache[category] = cacheItem{
		time:  time.Now(),
		value: value,
	}

	mutex.Unlock()
}
