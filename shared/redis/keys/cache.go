package keys

const (
	cachePrefix = "cache:"
)

func CacheKey(key string) string {
	return cachePrefix + key
}
