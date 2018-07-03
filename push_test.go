package cyfe

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimplePrepare(t *testing.T) {
	os.Setenv("CYFE_TOKEN_Test", "test")
	setup()

	send, err := Prepare("Test", "1", "", "", nil)
	require.Nil(t, err)
	assert.NotNil(t, send)
	assert.Equal(t, send.Data[0]["Date"], time.Now().UTC().Format("20060102"))

	b, err := json.Marshal(send)
	require.Nil(t, err)
	assert.NotZero(t, len(b))
	assert.NotEmpty(t, string(b))

}

func TestAllOptions(t *testing.T) {
	os.Setenv("CYFE_TOKEN_TestToken", "thanks_for_all_the_fish")
	setup()
	options := PushOptions{
		ReplaceInstead:      true,
		Color:               "#000000",
		Type:                "Line",
		IsCumulative:        true,
		DisplayAverages:     true,
		OverwriteTotal:      true,
		OverwriteComparison: true,
		IsBad:               true,
		IsUpsideDownGraph:   true,
		UnsyncYAxis:         true,
		YAxisMin:            "-2",
		YAxisMax:            "10",
		YAxisShow:           true,
		ShowLabel:           true,
	}
	request, err := Prepare("TestToken", "10", "", "", &options)
	assert.Nil(t, err)
	assert.NotNil(t, request)

	// check all the values for sanity
	assert.Equal(t, request.Data[0]["Date"], time.Now().UTC().Format("20060102"))
	assert.Equal(t, (*request.OnDuplicate)["TestToken"], "replace")
	assert.Equal(t, (*request.Color)["TestToken"], "#000000")
	assert.Equal(t, (*request.Type)["TestToken"], "Line")
	assert.Equal(t, (*request.Cumulative)["TestToken"], "1")
	assert.Equal(t, (*request.Average)["TestToken"], "1")
	assert.Equal(t, (*request.Total)["TestToken"], "1")
	assert.Equal(t, (*request.Comparison)["TestToken"], "1")
	assert.Equal(t, (*request.Reverse)["TestToken"], "1")
	assert.Equal(t, (*request.ReverseGraph)["TestToken"], "1")
	assert.Equal(t, (*request.YAxis)["TestToken"], "1")
	assert.Equal(t, (*request.YAxisMin)["TestToken"], "-2")
	assert.Equal(t, (*request.YAxisMax)["TestToken"], "10")
	assert.Equal(t, (*request.YAxisShow)["TestToken"], "1")
	assert.Equal(t, (*request.LabelShow)["TestToken"], "1")

	result, err := Push(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestTokenOption(t *testing.T) {
	options := PushOptions{
		ReplaceInstead: true,
		Token:          "override",
	}
	request, err := Prepare("TestToken", "10", "", "", &options)
	assert.Nil(t, err)
	assert.NotNil(t, request)
	assert.Equal(t, request.Data[0]["Date"], time.Now().UTC().Format("20060102"))
}

func TestBadPush(t *testing.T) {
	options := PushOptions{
		ReplaceInstead: true,
		Token:          "badtoken",
	}
	request, err := Prepare("Test", "10", "", "", &options)
	assert.Nil(t, err)
	assert.NotNil(t, request)

	result, err := Push(request)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestJustPushProductionError(t *testing.T) {
	os.Setenv("CYFE_ENV", "production")
	os.Setenv("CYFE_METRIC_UNITTEST", "production")
	setup()
	request, ret, err := JustPush("UNITTEST", "1")
	assert.NotNil(t, err)
	assert.NotNil(t, ret)
	assert.NotNil(t, request)
	os.Setenv("CYFE_METRIC_UNITTEST", "test")
}
