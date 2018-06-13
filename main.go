package cyfe

type PushOptions struct {
	// DuplicateOrReplace should be set to replace if you want to replace data on a duplicate post. Otherwise it will
	// add the data to the metric
	DuplicatOrReplace string
	Color             string
	Type              string
}

func Push(metricLabel string, key string, metricValue string) {
	if key == "" || key == "date" {
		// make the key the current date
	}
	// build out the data structure
	send := map[string]interface{}{}
	send["data"] = []map[string]string{
		map[string]string{
			"Date":      key,
			metricLabel: metricValue,
		},
	}
	// now we loop over the options to build the map to send

}
