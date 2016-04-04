package data

import (
	"strings"

	"google.golang.org/appengine/memcache"
)

// CacheRSS stores the RSS for a specific namespace in a cache
func CacheRSS(con *Context, namespace string, value string) error {
	namespace = strings.ToLower(namespace)
	return memcacheStore(con, "RSS"+namespace, []byte(value))
}

// ReadRSSCache checks if there is an entry for the RSS feed of namespace and
// returns it
func ReadRSSCache(con *Context, namespace string) (string, error) {
	namespace = strings.ToLower(namespace)
	b, err := memcacheRead(con, "RSS"+namespace)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// CacheJSON store the JSON for a specific namespace in a cache
func CacheJSON(con *Context, namespace string, value []byte) error {
	namespace = strings.ToLower(namespace)
	return memcacheStore(con, "JSON"+namespace, value)
}

// ReadJSONCache checks if there is an entry for the JOSN of namespace and
// returns it
func ReadJSONCache(con *Context, namespace string) ([]byte, error) {
	namespace = strings.ToLower(namespace)
	return memcacheRead(con, "JSON"+namespace)
}

// clearCache removes all data for the given namespace from the cache
func clearCache(con *Context, namespace string) error {
	namespace = strings.ToLower(namespace)
	return memcache.DeleteMulti(con.C, []string{"RSS" + namespace, "JSON" + namespace})
}

func memcacheStore(con *Context, key string, value []byte) error {
	i := &memcache.Item{
		Key:   key,
		Value: value,
	}
	err := memcache.Set(con.C, i)
	return err
}

func memcacheRead(con *Context, key string) ([]byte, error) {
	i, err := memcache.Get(con.C, key)
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
