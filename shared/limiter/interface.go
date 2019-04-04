package limiter

type Interface interface {
	Allow(string) (bool, error)
}
