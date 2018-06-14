package cyfe

import (
	"time"
)

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
	SyncYAxis         bool
	YAxisMin          string
	YAxisMax          string
	YAxisShow         bool
	ShowLabel         bool
}

type PushSendRequest struct {
	Data        []map[string]string
	OnDuplicate map[string]string
}

func Prepare(metricLabel, metricValue, key string, options *PushOptions) (send *PushSendRequest, err error) {
	send = CreateDefaultSendRequest(metricLabel)
	if key == "" || key == "date" {
		// make the key the current date
		key = time.Now().UTC().Format("20060102")
	}
	// build out the data structure
	send.Data = []map[string]string{
		map[string]string{
			"Date":      key,
			metricLabel: metricValue,
		},
	}
	// now we loop over the options to build the map to send
	if options != nil {
		if options.ReplaceInstead {
			send.OnDuplicate = map[string]string{
				metricLabel: "replace",
			}
		}
	}

	return
}

func CreateDefaultSendRequest(metricLabel string) (send *PushSendRequest) {
	return &PushSendRequest{
		OnDuplicate: map[string]string{
			metricLabel: "duplicate",
		},
	}
}
