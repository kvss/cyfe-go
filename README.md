# cyfe-go

cyfe-go is a Go SDk for integrating with the Cyfe Push API as described in the [docs](https://www.cyfe.com/api). Configuration is handled through a file or through the environment.

## Installation

Using your favorite dependency manager, make sure to add `github.com/kvss/cyfe-go`. For example, with dep,

`dep ensure -add github.com/kvss/cyfe-go`

## Usage

Usage is fairly straight-forward. See Configuration for more information about setting up and configuring the SDK.

**Important:** No calls will actually be made unless you set `CYFE_ENV` to `production`. This is to prevent populating your widgets with data from things like unit tests or development environments.

Once configured, there are two primary APIs to interact with. `Push` is the primary function, allowing for customizing the metrics. `JustPush` uses only the defaults and exists to be a simpler call if no further customization is needed.

```go
request, ret, err := JustPush("User Signup", "1")
```

```go
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
request, err := Prepare("User Signup", "10", "", "", &options)
result, err := Push(request)
```

## Configuration

Cyfe requires a full API end point for each widget that is pushed to. Obviously, it could be a challenge to hard code all of these "tokens". So, essentially, we allow you to map metrics to tokens using a TOML file or the environment. We recommend TOML, since the environment does not allow things such as spaces in metric names whereas TOML does.

The order of loading is the TOML file first. Then the environment is read and *will replace any duplicates found in the file*. Therefore, a common practice would be to populate the TOML file and then replace any environment-specific points with the environment.

**If a metric is used but no entry is found for the metric in either the environment or the file, the calls will fail**. We don't know where to send the metric. Future versions may allow a default widget to be specified.

## Environment Variables

* `CYFE_ENV` set to production to actually make calls

* `CYFE_TIMEZONE` the timezone used for default dates if no key is pushed; leave blank for UTC. Uses timezone database: `CYFE_TIMEZONE=America/New_York`. Any valid timezone from the [go time docs](https://golang.org/pkg/time/#Location) is valid.

* `CYFE_TOKEN_FILE` the name (no path, no extension) of a toml configuration file of metric/chart token pairs (see sample.toml)

* `CYFE_TOKEN_*` after parsing the file (if provided), any `CYFE_TOKEN_*` environment variables will be parsed and added

## Other Libraries

We use the following additional tools in this library, and thank the maintainers and contributors of those libraries:

* [testify](https://github.com/stretchr/testify) - Makes our unit tests more readable and management

* [viper](https://github.com/spf13/viper) - Interacting with configurations files made better

## Bugs

There is currently a bug in the Cyfe docs (as of 20180702). The docs mention that if onduplicate is set to replace, it would replace the data instead of accumulating. The behavior we are seeing is that passing anything for the field triggers the replacement, even if we send in a blank string.