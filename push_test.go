package cyfe

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimplePrepare(t *testing.T) {
	send, err := Prepare("Test", "1", "", "", nil)
	assert.Nil(t, err)
	assert.NotNil(t, send)
	assert.Equal(t, send.Data[0]["Date"], time.Now().UTC().Format("20060102"))

	b, err := json.Marshal(send)
	assert.Nil(t, err)
	assert.NotZero(t, len(b))
	assert.NotEmpty(t, string(b))
}

func TestAllOptions(t *testing.T) {
	options := PushOptions{
		ReplaceInstead:      true,
		Color:               "#000000",
		Type:                "Area",
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
	request, err := Prepare("Test", "10", "", "", &options)
	assert.Nil(t, err)
	assert.NotNil(t, request)
	fmt.Printf("\n%+v\n", request)
	// check all the values for sanity
	assert.Equal(t, request.Data[0]["Date"], time.Now().UTC().Format("20060102"))
	assert.Equal(t, (*request.OnDuplicate)["Test"], "replace")
	assert.Equal(t, (*request.Color)["Test"], "#000000")
	assert.Equal(t, (*request.Type)["Test"], "Area")
	assert.Equal(t, (*request.Cumulative)["Test"], "1")
	assert.Equal(t, (*request.Average)["Test"], "1")
	assert.Equal(t, (*request.Total)["Test"], "1")
	assert.Equal(t, (*request.Comparison)["Test"], "1")
	assert.Equal(t, (*request.Reverse)["Test"], "1")
	assert.Equal(t, (*request.ReverseGraph)["Test"], "1")
	assert.Equal(t, (*request.YAxis)["Test"], "1")
	assert.Equal(t, (*request.YAxisMin)["Test"], "-2")
	assert.Equal(t, (*request.YAxisMax)["Test"], "10")
	assert.Equal(t, (*request.YAxisShow)["Test"], "1")
	assert.Equal(t, (*request.LabelShow)["Test"], "1")
}
