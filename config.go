package cyfe

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config = &configStruct{}

type configStruct struct {
	environment   string
	cyfeRoot      string
	metricLookups []metricLookup
}

type metricLookup struct {
	Metric string
	Token  string
}

func init() {

	setup()
}

func isProd() (isProd bool) {
	if config.environment == "production" {
		isProd = true
	}
	return
}

func setup() (err error) {
	// setup the application; this is broken out from the init()
	// so that it may be called by unit tests to change env vars
	config.environment = strings.ToLower(os.Getenv("CYFE_ENV"))
	if config.environment == "" {
		config.environment = "dev"
	}
	config.cyfeRoot = "https://app.cyfe.com/api/push/"
	config.metricLookups = []metricLookup{}

	// each push to the API requires a chart token. The specific metric to specific chart
	// can come from the file (CYFE_TOKEN_FILE) or can come from the environment
	// (looping over all CYFE_TOKEN_* environment variables)
	// this should only be created during initialization

	// file first
	fileLocation := os.Getenv("CYFE_TOKEN_FILE")
	if fileLocation != "" {
		viper.AddConfigPath(".")
		viper.SetConfigName(fileLocation)
		err := viper.ReadInConfig()
		if err != nil {
			//  the user said to find a file, but we couldn't so we error out
			// of course, panicking in an included library is bad, so just output a warning
			// and hope they see it
			// TODO: integrate logrus here
			err = fmt.Errorf("ERROR: CYFE_TOKEN_FILE named but not found")
			fmt.Println(err.Error())
			return err
		}
		parsed := []metricLookup{}
		// metrics := viper.Get("metric")
		err = viper.UnmarshalKey("metric", &parsed)
		if err != nil || len(parsed) == 0 {
			//  the config could not be parsed
			err = fmt.Errorf("ERROR: CYFE_TOKEN_FILE could not be parsed")
			fmt.Println(err.Error())
			return err
		}
		fmt.Printf("\n%+v\n%+v\n", err, config.metricLookups)
		config.metricLookups = parsed
	}
	// TODO: do we want to consider a default toml file?

	// now the env
	env := os.Environ()
	for i := range env {
		if strings.HasPrefix(env[i], "CYFE_TOKEN") {
			split := strings.Split(env[i], "=")
			metric := strings.TrimPrefix(split[0], "CYFE_TOKEN_")
			config.metricLookups = append(config.metricLookups, metricLookup{
				Metric: metric,
				Token:  split[1],
			})
		}
	}
	return nil
}
