package env

import "os"

const (
	envFixerAPIKey = "FIXER_API_KEY"
)

func FixerAPIKey() (string, bool) {
	return os.LookupEnv(envFixerAPIKey)
}
