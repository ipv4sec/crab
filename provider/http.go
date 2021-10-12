package provider

import (
	"github.com/gojek/heimdall/v7/httpclient"
	"time"
)

var (
	HTTPClient = httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))
)