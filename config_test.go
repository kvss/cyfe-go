package cyfe

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenFromEnvironmentSetup(t *testing.T) {
	os.Setenv("CYFE_TOKEN_UNITTEST", "12345")
	os.Setenv("CYFE_TOKEN_KVSS_TEST", "KVSS")

	setup()

	assert.NotZero(t, len(config.metricLookups))
	// loop and find them
	foundTest := false
	foundKVSS := false
	for i := range config.metricLookups {
		if config.metricLookups[i].Metric == "UNITTEST" && config.metricLookups[i].Token == "12345" {
			foundTest = true
		}
		if config.metricLookups[i].Metric == "KVSS_TEST" && config.metricLookups[i].Token == "KVSS" {
			foundKVSS = true
		}
	}
	assert.True(t, foundTest)
	assert.True(t, foundKVSS)
}

func TestTokenFromFileSetup(t *testing.T) {
	os.Setenv("CYFE_TOKEN_FILE", "sample")

	setup()

	assert.NotZero(t, len(config.metricLookups))
	// loop and find them
	foundSignup := false
	foundDelete := false
	for i := range config.metricLookups {
		if config.metricLookups[i].Metric == "User Signup" && config.metricLookups[i].Token == "abc123" {
			foundSignup = true
		}
		if config.metricLookups[i].Metric == "User Deleted" && config.metricLookups[i].Token == "123abc" {
			foundDelete = true
		}
	}
	assert.True(t, foundSignup)
	assert.True(t, foundDelete)
}

func TestBadFile(t *testing.T) {
	os.Setenv("CYFE_TOKEN_FILE", "noexist")
	err := setup()
	assert.NotNil(t, err)
	os.Setenv("CYFE_TOKEN_FILE", "sample_bad")
	err = setup()
	assert.NotNil(t, err)

	os.Setenv("CYFE_TOKEN_FILE", "")
}
