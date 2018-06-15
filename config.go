package cyfe

import (
	"os"
	"strings"
)

var config = &configStruct{}

type configStruct struct {
	environment string
	cyfeRoot    string
}

func init() {
	config.environment = strings.ToLower(os.Getenv("CYFE_ENV"))
	if config.environment == "" {
		config.environment = "dev"
	}
	config.cyfeRoot = "https://app.cyfe.com/api/push/"
}

func isProd() (isProd bool) {
	if config.environment == "production" {
		isProd = true
	}
	return
}
