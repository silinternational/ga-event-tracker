package ga

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultHTTPClientTimeout = 5 * time.Second
	Endpoint                 = "https://www.google-analytics.com/collect"
	ProtocolVersion          = "1"
	UserAgent                = "github.com/silinternational/ga-event-tracker"
)

// Event holds the event specific values as well as standard values required for every call to to the measurement protocol
type Event struct {
	// TrackingID - [REQUIRED] The GA property ID, example: UA-123456-1
	TrackingID string

	// ClientID - [REQUIRED] A unique identifier for the user or thing you want this event attached to.
	// For example we use this library to track events for pull request merges so we send the repository name through
	// as the ClientID
	ClientID string

	// Category - [REQUIRED] Specifies the event category. Must not be empty.
	Category string

	// Action - [REQUIRED] Specifies the event action. Must not be empty.
	Action string

	// Label - [OPTIONAL] Specifies the event label.
	Label string

	// Value - [OPTIONAL] Specifies the event value. Values must be non-negative.
	Value int
}

func (e *Event) IsValid() (bool, error) {
	if e.TrackingID == "" || !strings.HasPrefix(e.TrackingID, "UA") {
		return false, fmt.Errorf("TrackingID cannot be empty and must start with UA")
	}

	if e.ClientID == "" {
		return false, fmt.Errorf("ClientID cannot be empty")
	}

	if e.Category == "" {
		return false, fmt.Errorf("category cannot be empty")
	}

	if e.Action == "" {
		return false, fmt.Errorf("action cannot be empty")
	}

	if e.Value < 0 {
		return false, fmt.Errorf("value cannot be negative")
	}

	return true, nil
}

func SendEvent(e Event) error {
	if isValid, err := e.IsValid(); !isValid {
		return err
	}

	values := url.Values{
		"t":   []string{"event"},
		"v":   []string{ProtocolVersion},
		"tid": []string{e.TrackingID},
		"cid": []string{e.ClientID},
		"ec":  []string{e.Category},
		"el":  []string{e.Label},
	}

	if e.Action != "" {
		values.Set("ea", e.Action)
	}

	if e.Value != 0 {
		values.Set("ev", strconv.Itoa(e.Value))
	}

	body := values.Encode()

	c := &http.Client{
		Timeout: DefaultHTTPClientTimeout,
	}

	req, err := http.NewRequest(http.MethodPost, Endpoint, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating new http request: %s", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("error calling to Google Analytics: %s", err)
	}

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	var resBody []byte
	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("did not get OK response, got code %v, unable to read response body: %s", res.StatusCode, err)
	}

	return fmt.Errorf("got error calling Google Analytics [status %v]: %s", res.StatusCode, string(resBody))
}
