package data

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
)

// CacheRSS stores the RSS for a specific namespace in a cache
func CacheRSS(c context.Context, namespace string, value string) error {
	namespace = strings.ToLower(namespace)
	return memcacheStore(c, "RSS"+namespace, []byte(value))
}

// ReadRSSCache checks if there is an entry for the RSS feed of namespace and
// returns it
func ReadRSSCache(c context.Context, namespace string) (string, error) {
	namespace = strings.ToLower(namespace)
	b, err := memcacheRead(c, "RSS"+namespace)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CacheJSON store the JSON for a specific namespace in a cache
func CacheJSON(c context.Context, namespace string, value []byte) error {
	namespace = strings.ToLower(namespace)
	return memcacheStore(c, "JSON"+namespace, value)
}

// ReadJSONCache checks if there is an entry for the JOSN of namespace and
// returns it
func ReadJSONCache(c context.Context, namespace string) ([]byte, error) {
	namespace = strings.ToLower(namespace)
	return memcacheRead(c, "JSON"+namespace)
}

// clearCache removes all data for the given namespace from the cache
func clearCache(c context.Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	return memcache.DeleteMulti(c, []string{"RSS" + namespace, "JSON" + namespace})
}

func memcacheStore(c context.Context, key string, value []byte) error {
	i := &memcache.Item{
		Key:   key,
		Value: value,
	}
	err := memcache.Set(c, i)
	return err
}

func memcacheRead(c context.Context, key string) ([]byte, error) {
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
