package consumer

import (
	log "github.com/mashmorsik/L0/pkg/logger"
	"github.com/mashmorsik/L0/test/testdata"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateMsg(t *testing.T) {
	log.BuildLogger()

	type args struct {
		msg nats.Msg
	}

	tests := []struct {
		name     string
		args     args
		expected error
	}{
		{name: "ValidateMsg_true",
			args: args{msg: nats.Msg{
				Data: []byte(testdata.JSONData)}},
			expected: nil},
		{name: "ValidateMsg_false_order_id_field_not_filled",
			args: args{msg: nats.Msg{
				Data: []byte(testdata.JSONDataEmptyId)}},
			expected: errors.New("Field validation for 'OrderUid' failed on the 'required' tag")},
		{name: "ValidateMsg_false_address_details_not_filled",
			args: args{msg: nats.Msg{
				Data: []byte(testdata.JSONDataEmptyAddressDetails)}},
			expected: errors.New("Field validation for 'Address' failed on the 'required' tag")},
		{name: "ValidateMsg_false_track_number_details_not_filled",
			args: args{msg: nats.Msg{
				Data: []byte(testdata.JSONDataEmptyTrackNumber)}},
			expected: errors.New("Field validation for 'TrackNumber' failed on the 'required' tag")},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := ValidateMsg(tt.args.msg)
			if got != nil {
				if !assert.Contains(t, got.Error(), tt.expected.Error()) {
					t.Errorf("ValidateMsg() = %v, expected %v", got, tt.expected)
				}
			}
		})
	}
}
