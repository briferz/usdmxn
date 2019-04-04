package banxico

import "os"

const (
	envBanxicoAPIKey = "BANXICO_API_KEY"
)

func banxicoAPIKey() (string, bool) {
	return os.LookupEnv(envBanxicoAPIKey)
}
