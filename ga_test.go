package ga

import "testing"

func TestEvent_IsValid(t *testing.T) {
	type fields struct {
		TrackingID string
		ClientID   string
		Category   string
		Action     string
		Label      string
		Value      int
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				TrackingID: tt.fields.TrackingID,
				ClientID:   tt.fields.ClientID,
				Category:   tt.fields.Category,
				Action:     tt.fields.Action,
				Label:      tt.fields.Label,
				Value:      tt.fields.Value,
			}
			got, err := e.IsValid()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}
