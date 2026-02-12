package model

import (
	"testing"
)

func TestCountRequest_Validate(t *testing.T) {
	countVal := 10
	negCount := -1

	tests := []struct {
		name    string
		req     CountRequest
		wantErr bool
		msg     string
	}{
		{
			name: "Valid request",
			req: CountRequest{
				ExternalID: "ext-1",
				Count:      &countVal,
			},
			wantErr: false,
		},
		{
			name: "Missing external_id",
			req: CountRequest{
				ExternalID: "",
				Count:      &countVal,
			},
			wantErr: true,
			msg:     "missing external_id",
		},
		{
			name: "Missing count",
			req: CountRequest{
				ExternalID: "ext-1",
				Count:      nil,
			},
			wantErr: true,
			msg:     "missing count",
		},
		{
			name: "Invalid count value",
			req: CountRequest{
				ExternalID: "ext-1",
				Count:      &negCount,
			},
			wantErr: true,
			msg:     "invalid count value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CountRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.msg {
				t.Errorf("CountRequest.Validate() error message = %v, want %v", err.Error(), tt.msg)
			}
		})
	}
}
