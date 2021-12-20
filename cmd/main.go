package main

import (
	"log"
	"os"

	"github.com/silinternational/ga-event-tracker"
)

func main() {
	log.SetOutput(os.Stdout)

	meta := ga.Meta{
		APISecret:     os.Getenv("GA_API_SECRET"),
		MeasurementID: os.Getenv("GA_MEASUREMENT_ID"),
		ClientID:      os.Getenv("GA_CLIENT_ID"),
		UserID:        os.Getenv("GA_USER_ID"),
	}

	name := os.Getenv("GA_EVENT_NAME")
	if name == "" {
		log.Fatal("Env var EVENT_NAME is required")
	}

	params, err := ga.GetParamsFromEnv(ga.DefaultParamsEnvVar, false)
	if err != nil {
		log.Fatal(err)
	}

	err = ga.SendEvent(meta, []ga.Event{
		{
			Name:   name,
			Params: params,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`Event "%s" sent`, name)
}
