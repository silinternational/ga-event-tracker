# Google Analytics - Event Tracker
This package is a lightweight wrapper for the 
[Google Analytics Measurement Protocol v4](https://developers.google.com/analytics/devguides/collection/protocol/ga4) 
specifically for Events.

We like to use GA to track various operational events/metrics like deploying apps, merging PRs, etc. So we just need 
a very simple way to send these events to GA. While the Measurement Protocol can handle a lot more than this package
supports, we've created this as minimal as possible to keep it as simple as possible to use.

If there are additional features of MP that you'd like supported, open an issue to discuss it with us and then possibly 
a pull request to add the feature, although try to keep the API the same.

## Usage

The `SendEvent` method takes two parameters, the first is a `Meta` struct intended to hold parameters that are not 
individual event related but needed for the API call. Initially it is just APISecret and UserID, though it 
could be extended to include other option fields like `timestamp_micros`, `user_properties` and `non_personalized_ads`.
See https://developers.google.com/analytics/devguides/collection/protocol/ga4/reference#payload_post_body for more info 
on the request body structure.

Each call to the MP API can include up to 25 events, so the second parameter is an array of `Event` structs. Each 
`Event` must have a `Name`, and can optionally contain a list of `Params` in form of `map[string]interface{}`. 

```go
package main

import (
	"github.com/silinternational/ga-event-tracker"
	"log"
)

func main() {
	err := ga.SendEvent(ga.Meta{
		APISecret: "asdf1234",
		ClientID: "abc1234",
		MeasurementID: "G-N1235ZM",
	}, []ga.Event{
		{
			Name: "custom_event",
			Params: ga.Params{
				"category": "something",
				"project":  "whatever",
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
```