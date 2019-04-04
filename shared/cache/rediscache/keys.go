package rediscache

const (
	cachePrefix = "cache:"
)

func cacheKey(key string) string {
	return cachePrefix + key
}
