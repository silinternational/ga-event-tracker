package ga

import (
	"os"
	"testing"
)

func TestMeta_Validate(t *testing.T) {
	type fields struct {
		APISecret string
		UserID    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty meta",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "valid",
			fields: fields{
				APISecret: "something",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meta{
				APISecret: tt.fields.APISecret,
				UserID:    tt.fields.UserID,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEvent_Validate(t *testing.T) {
	type fields struct {
		Name   string
		Params map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty events",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "valid, no params",
			fields: fields{
				Name: "testing",
			},
			wantErr: false,
		},
		{
			name: "reserved event name",
			fields: fields{
				Name: "first_open",
			},
			wantErr: true,
		},
		{
			name: "reserved param prefix",
			fields: fields{
				Name: "custom_event",
				Params: Params{
					"google_param_name": "some value",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				Name:   tt.fields.Name,
				Params: tt.fields.Params,
			}
			if err := e.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendEvent(t *testing.T) {
	type args struct {
		meta   Meta
		events []Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid event, no params",
			args: args{
				meta: Meta{
					APISecret: os.Getenv("GA_API_SECRET"),
				},
				events: []Event{
					{
						Name:   "test-event-no-params",
						Params: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid event, with params",
			args: args{
				meta: Meta{
					APISecret: os.Getenv("GA_API_SECRET"),
				},
				events: []Event{
					{
						Name: "test-event-with-params",
						Params: Params{
							"project":     "ga-event-tracking",
							"environment": "dev",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendEvent(tt.args.meta, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("SendEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
