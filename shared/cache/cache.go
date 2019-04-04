package cache

type Interface interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
}
