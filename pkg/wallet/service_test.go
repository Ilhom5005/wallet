package wallet

import (
	"reflect"
	"testing"

	"github.com/Ilhom5005/wallet/v1/pkg/types"
)

func TestService_FindAccountByID(t *testing.T) {
	type args struct {
		accountID int
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *types.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FindAccountByID(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FindAccountByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
