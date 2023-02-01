package main

import (
	"encoding/json"
	"log"
	"net/mail"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/MicahParks/jsontype"
)

const exampleConfig = `{
  "ends": "Wed, 04 Oct 2022 00:00:00 MST",
  "getInterval": "1h30m",
  "notificationMsg": "Your item is on sale!",
  "notify": "EXAMPLE@example.com",
  "targetPage": "https://www.example.com",
  "targetRegExp": "example",
  "targetUUID": "84abbfc2-b7a8-4446-a351-927c0fd26a3a"
}`

type myConfig struct {
	Ends            *jsontype.JSONType[time.Time]      `json:"ends"`
	GetInterval     *jsontype.JSONType[time.Duration]  `json:"getInterval"`
	NotificationMsg string                             `json:"notificationMsg"`
	Notify          *jsontype.JSONType[*mail.Address]  `json:"notify"`
	TargetPage      *jsontype.JSONType[*url.URL]       `json:"targetPage"`
	TargetRegExp    *jsontype.JSONType[*regexp.Regexp] `json:"targetRegExp"`
}

func main() {
	logger := log.New(os.Stdout, "", 0)
	var config myConfig

	// Set non-default unmarshal behavior.
	endOpts := jsontype.Options{
		TimeFormatUnmarshal: time.RFC1123,
	}
	config.Ends = jsontype.NewWithOptions(time.Time{}, endOpts)

	// Unmarshal the configuration.
	err := json.Unmarshal(json.RawMessage(exampleConfig), &config)
	if err != nil {
		logger.Fatalf("failed to unmarshal JSON: %s", err)
	}

	// Access fields on the unmarshalled configuration.
	logger.Printf("Ends: %s", config.Ends.Get().String())
	logger.Printf("Get interval: %s", config.GetInterval.Get().String())

	// Set non-default marshal behavior.
	emailOpts := jsontype.Options{
		MailAddressAddressOnlyMarshal: true,
		MailAddressLowerMarshal:       true,
	}
	config.Notify = jsontype.NewWithOptions(config.Notify.Get(), emailOpts)

	// Marshal the configuration back to JSON.
	remarshaled, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		logger.Fatalf("failed to re-marshal configuration: %s", err)
	}
	logger.Println(string(remarshaled))
}
