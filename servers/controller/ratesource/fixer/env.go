package fixer

import "os"

const (
	envFixerAPIKey = "FIXER_API_KEY"
)

func fixerAPIKey() (string, bool) {
	return os.LookupEnv(envFixerAPIKey)
}
