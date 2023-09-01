package round_robin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestWeightBalancer_Pick(t *testing.T) {
	type fields struct {
		connections []*weightConn
	}
	type args struct {
		info balancer.PickInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    balancer.PickResult
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WeightBalancer{
				connections: tt.fields.connections,
			}
			got, err := w.Pick(tt.args.info)
			if !tt.wantErr(t, err, fmt.Sprintf("Pick(%v)", tt.args.info)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Pick(%v)", tt.args.info)
		})
	}
}
