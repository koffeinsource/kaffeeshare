package data

import (
	"appengine"
	"appengine/memcache"
)

// CacheRSS stores the RSS for a specific namespace in a cache
func CacheRSS(c appengine.Context, namespace string, value string) error {
	return memcacheStore(c, "RSS"+namespace, []byte(value))
}

// ReadRSSCache checks if there is an entry for the RSS feed of namespace and
// returns it
func ReadRSSCache(c appengine.Context, namespace string) (string, error) {
	b, err := memcacheRead(c, "RSS"+namespace)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CacheJSON store the JSON for a specific namespace in a cache
func CacheJSON(c appengine.Context, namespace string, value []byte) error {
	return memcacheStore(c, "JSON"+namespace, value)
}

// ReadJSONCache checks if there is an entry for the JOSN of namespace and
// returns it
func ReadJSONCache(c appengine.Context, namespace string) ([]byte, error) {
	return memcacheRead(c, "JSON"+namespace)
}

// clearCache removes all data for the given namespace from the cache
func clearCache(c appengine.Context, namespace string) error {
	return memcache.DeleteMulti(c, []string{"RSS" + namespace, "JSON" + namespace})
}

func memcacheStore(c appengine.Context, key string, value []byte) error {
	i := &memcache.Item{
		Key:   key,
		Value: value,
	}
	err := memcache.Set(c, i)
	return err
}

func memcacheRead(c appengine.Context, key string) ([]byte, error) {
	i, err := memcache.Get(c, key)
	if err == memcache.ErrCacheMiss {
		// cache miss
		return nil, err
	}
	if err != nil {
		// real error
		return nil, err
	}
	// hit
	return i.Value, nil
}
