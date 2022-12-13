package ga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultHTTPClientTimeout = 5 * time.Second
	DefaultParamsEnvVar      = "GA_EVENT_PARAMS"
	Endpoint                 = "https://www.google-analytics.com/mp/collect?api_secret=%s&measurement_id=%s"
	DebugEndpoint            = "https://www.google-analytics.com/debug/mp/collect?api_secret=%s&measurement_id=%s"
	UserAgent                = "github.com/silinternational/ga-event-tracker"
)

// gaEventBody represents the POST body for reporting events to GA
type gaEventBody struct {
	ClientID string  `json:"client_id"`
	UserID   string  `json:"user_id,omitempty"`
	Events   []Event `json:"events"`
}

type Meta struct {
	// APISecret - [REQUIRED] An API SECRET generated in the Google Analytics UI. To create a new secret, navigate to:
	// Admin > Data Streams > choose your stream > Measurement Protocol > Create
	APISecret string

	// MeasurementID - [REQUIRED] Measurement ID. The identifier for a Data Stream. Found in the Google Analytics UI under:
	// Admin > Data Streams > choose your stream > Measurement ID >
	MeasurementID string

	// ClientID - [REQUIRED] Uniquely identifies a user instance of a web client.
	ClientID string

	// UserID - [OPTIONAL] A unique identifier for a user
	UserID string
}

func (m *Meta) Validate() error {
	if m.APISecret == "" {
		return fmt.Errorf("APISecret cannot be empty")
	}

	if m.MeasurementID == "" {
		return fmt.Errorf("MeasurementID cannot be empty")
	}

	if m.ClientID == "" {
		return fmt.Errorf("ClientID cannot be empty")
	}

	return nil
}

type Params map[string]interface{}

// Event holds the event specific values as well as standard values required for every call to to the measurement protocol
type Event struct {
	// Name - [REQUIRED] Name of event, alphanumeric and underscores only.
	Name string `json:"name"`

	// Params - [OPTIONAL] - Any additional parameters to attach to event
	Params Params `json:"params,omitempty"`
}

func (e *Event) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if isStringInSlice(e.Name, ReservedEventNames) {
		return fmt.Errorf("the event name %s is reserved by Google Analytics", e.Name)
	}

	for key := range e.Params {
		for _, prefix := range ReservedParamPrefixes {
			if strings.HasPrefix(key, prefix) {
				return fmt.Errorf("event %s has param with reserved prefix %s, param: %s", e.Name, prefix, key)
			}
		}
	}

	return nil
}

func callGA(endpoint string, reqBody []byte, meta Meta) (*http.Response, string, error) {
	c := &http.Client{
		Timeout: DefaultHTTPClientTimeout,
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(endpoint, meta.APISecret, meta.MeasurementID),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, "", fmt.Errorf("error creating new http request: %s", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error calling to Google Analytics: %s", err)
	}

	var resBody []byte
	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		return res, "", fmt.Errorf("error reading response body: %s", err)
	}

	return res, fmt.Sprintf("%s", resBody), nil
}

func SendEvent(meta Meta, events []Event) error {
	if err := meta.Validate(); err != nil {
		return err
	}

	for i, ev := range events {
		if err := ev.Validate(); err != nil {
			return fmt.Errorf("validation error for event #%v: %s", i, err)
		}
	}

	gaEv := gaEventBody{
		ClientID: meta.ClientID,
		UserID:   meta.UserID,
		Events:   events,
	}

	// URL encode values for payload
	body, err := json.Marshal(gaEv)
	if err != nil {
		return fmt.Errorf("uanble to marshal event to json: %s", err)
	}

	// Call debug endpoint. The GA4 api fails silently, otherwise.
	res, resBody, err := callGA(DebugEndpoint, body, meta)
	if err != nil {
		return fmt.Errorf("error making call to the debug endpoint: %s", err)
	}
	log.Printf("Results of debug test call. Status: %d. Body: %s", res.StatusCode, resBody)

	// Call non-debug endpoint
	res, resBody, err = callGA(Endpoint, body, meta)
	if err != nil {
		return fmt.Errorf("error making call to the non-debug endpoint: %s", err)
	}

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	return fmt.Errorf("got error calling Google Analytics [status %v]: %s", res.StatusCode, string(resBody))
}

// ReservedEventNames - names reserved by GA, cannot be used for custom events
var ReservedEventNames = []string{
	"ad_activeview",
	"ad_click",
	"ad_exposure",
	"ad_impression",
	"ad_query",
	"adunit_exposure",
	"app_clear_data",
	"app_install",
	"app_update",
	"app_remove",
	"error",
	"first_open",
	"first_visit",
	"in_app_purchase",
	"notification_dismiss",
	"notification_foreground",
	"notification_open",
	"notification_receive",
	"os_update",
	"screen_view",
	"session_start",
	"user_engagement",
}

var ReservedParamPrefixes = []string{
	"google_",
	"ga_",
	"firebase_",
}

// IsStringInSlice iterates over a slice of strings, looking for the given
// string. If found, true is returned. Otherwise, false is returned.
func isStringInSlice(needle string, haystack []string) bool {
	for _, hs := range haystack {
		if needle == hs {
			return true
		}
	}

	return false
}

func GetParamsFromEnv(varName string, required bool) (Params, error) {
	if varName == "" {
		varName = DefaultParamsEnvVar
	}

	value := os.Getenv(varName)
	if value == "" && !required {
		return nil, nil
	}
	if value == "" && required {
		return nil, fmt.Errorf("required params env var %s is empty", varName)
	}

	var params Params
	if err := json.Unmarshal([]byte(value), &params); err != nil {
		return nil, fmt.Errorf("Value of params env var %s does not appear to be JSON, error: %s", varName, err)
	}

	return params, nil
}
