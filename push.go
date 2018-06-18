// Package cyfe implements the basic Push API for Cyfe in Go according to https://www.cyfe.com/api
// It is important to remember that no calls will actually be made if CYFE_ENV is not set to production. See Push() for more information
package cyfe

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

// PushOptions are a set of options for a specific Push request. Sane defaults are used (hence why the names differ than in the Cyfe API docs) so that, for
// example, false equates to the default behavior for instantiation
type PushOptions struct {
	// ReplaceInstead toggles whether a new push adds to the value (false, default) or replaces the value (true)
	ReplaceInstead      bool
	Color               string
	Type                string
	IsCumulative        bool
	DisplayAverages     bool
	OverwriteTotal      bool
	OverwriteComparison bool
	// IsBad specifies if a higher value for this metric is "bad" (see reverse in docs)
	IsBad             bool
	IsUpsideDownGraph bool
	UnsyncYAxis       bool
	YAxisMin          string
	YAxisMax          string
	YAxisShow         bool
	ShowLabel         bool
	// Token is available to override a token read from the end point
	Token string
}

// PushSendRequest is the formatted request to be sent to the Cyfe server
type PushSendRequest struct {
	Data         []map[string]string `json:"data,omitempty"`
	OnDuplicate  *map[string]string  `json:"onduplicate,omitempty"`
	Color        *map[string]string  `json:"color,omitempty"`
	Type         *map[string]string  `json:"type,omitempty"`
	Cumulative   *map[string]string  `json:"cumulative,omitempty"`
	Average      *map[string]string  `json:"average,omitempty"`
	Total        *map[string]string  `json:"total,omitempty"`
	Comparison   *map[string]string  `json:"comparison,omitempty"`
	Reverse      *map[string]string  `json:"reverse,omitempty"`
	ReverseGraph *map[string]string  `json:"reversegraph,omitempty"`
	YAxis        *map[string]string  `json:"yaxis,omitempty"`
	YAxisMin     *map[string]string  `json:"yaxismin,omitempty"`
	YAxisMax     *map[string]string  `json:"yaxismax,omitempty"`
	YAxisShow    *map[string]string  `json:"yaxisshow,omitempty"`
	LabelShow    *map[string]string  `json:"labelshow,omitempty"`
	ChartToken   string              `json:"-"`
}

// APIReturn is the success return from a successful push
type APIReturn struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// Prepare prepares a metric to be sent by filling out all of the options and formatting the data. Currently, you can
// only push one metric per call. A future improvement would be to allow multiple metricLabel/metricValue pairs. If keyLabel is empty or keyLabel is date
// AND keyValue is empty, the keyLabel will be set to Date (intentional capitalization as per the docs) and the current UTC timestamp
func Prepare(metricLabel, metricValue, keyLabel, keyValue string, options *PushOptions) (request *PushSendRequest, err error) {
	request = CreateDefaultSendRequest(metricLabel)
	if keyLabel == "" || (keyLabel == "date" && keyValue == "") {
		// make the key the current date
		keyLabel = "Date"
		keyValue = time.Now().UTC().Format("20060102")
	}
	// lookup the chart token; if we can't find it, error
	chartToken := ""
	// first, check if there is an override
	if options != nil && options.Token != "" {
		chartToken = options.Token
	} else {
		for i := range config.metricLookups {
			if config.metricLookups[i].Metric == metricLabel {
				chartToken = config.metricLookups[i].Token
				break
			}
		}
	}
	if chartToken == "" {
		err = errors.New("chart token not found for metric " + metricLabel)
		return
	}
	request.ChartToken = chartToken

	// build out the data structure
	request.Data = []map[string]string{
		map[string]string{
			keyLabel:    keyValue,
			metricLabel: metricValue,
		},
	}
	// now we loop over the options to build the map to send
	if options != nil {
		if options.ReplaceInstead {
			request.OnDuplicate = &map[string]string{
				metricLabel: "replace",
			}
		}
		if options.Color != "" {
			request.Color = &map[string]string{
				metricLabel: options.Color,
			}
		}
		if options.Type != "" {
			request.Type = &map[string]string{
				metricLabel: options.Type,
			}
		}
		if options.IsCumulative {
			request.Cumulative = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.DisplayAverages {
			request.Average = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.OverwriteTotal {
			request.Total = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.OverwriteComparison {
			request.Comparison = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.IsBad {
			request.Reverse = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.IsUpsideDownGraph {
			request.ReverseGraph = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.UnsyncYAxis {
			request.YAxis = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.YAxisMin != "" {
			request.YAxisMin = &map[string]string{
				metricLabel: options.YAxisMin,
			}
		}
		if options.YAxisMax != "" {
			request.YAxisMax = &map[string]string{
				metricLabel: options.YAxisMax,
			}
		}
		if options.YAxisShow {
			request.YAxisShow = &map[string]string{
				metricLabel: "1",
			}
		}
		if options.ShowLabel {
			request.LabelShow = &map[string]string{
				metricLabel: "1",
			}
		}
	}

	return
}

// JustPush is a simpler Push implementation which uses just the defaults. Useful if you don't like typing and just want to get
// a metric to the server
func JustPush(metricLabel, metricValue string) (request *PushSendRequest, ret APIReturn, err error) {
	request, err = Prepare(metricLabel, metricValue, "", "", nil)
	if err != nil {
		return
	}
	ret, err = Push(request)
	return
}

// Push actually makes the push request. NOTE: If the CYFE_ENV environment variable is not set to production, the request is
// NOT actually sent. This is to prevent accidentally sending metrics in test or development environments.
func Push(request *PushSendRequest) (ret APIReturn, err error) {
	if request == nil {
		err = errors.New("request cannot be nil")
	}
	if strings.HasPrefix(request.ChartToken, "/") {
		request.ChartToken = request.ChartToken[1:]
	}
	// is the environment isn't production, we don't make the call
	if !isProd() {
		// return as if it was a code call
		// if the token is set to badtoken, we are likely in a test and want to return an error
		if request.ChartToken == "badtoken" {
			ret = APIReturn{
				Status:     "error",
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid widget key",
			}
			err = errors.New(ret.Message)
		} else {
			ret = APIReturn{
				Status:     "ok",
				StatusCode: http.StatusOK,
				Message:    "Data pushed",
			}
		}
		return
	}
	httpRequest, err := resty.R().
		SetHeader("Accept", "application/json").
		SetBody(request).
		Post(fmt.Sprintf("%s%s", config.cyfeRoot, request.ChartToken))
	if err != nil {
		return
	}

	err = json.Unmarshal(httpRequest.Body(), &ret)
	if err != nil {
		err = errors.New("could not unmarshal the JSON response; check the API")
		return
	}
	ret.StatusCode = httpRequest.StatusCode()

	if httpRequest.StatusCode() != http.StatusOK {
		// the err should be the parsed message
		err = errors.New(ret.Message)
	}

	return
}

// CreateDefaultSendRequest initializes sane defaults for the send request
func CreateDefaultSendRequest(metricLabel string) (send *PushSendRequest) {
	return &PushSendRequest{
		OnDuplicate: &map[string]string{
			metricLabel: "duplicate",
		},
	}
}
