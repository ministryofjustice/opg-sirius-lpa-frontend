package sirius

import (
	"time"
)

type cacheItem struct {
	time  time.Time
	value []RefDataItem
}

var cache map[string]cacheItem

func getCached(category string) ([]RefDataItem, bool) {
	var v []RefDataItem
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	if cache[category].time.After(oneHourAgo) && len(cache[category].value) > 0 {
		return cache[category].value, true
	}

	return v, false
}

func setCached(category string, value []RefDataItem) {
	if cache == nil {
		cache = map[string]cacheItem{}
	}

	cache[category] = cacheItem{
		time:  time.Now(),
		value: value,
	}
}
