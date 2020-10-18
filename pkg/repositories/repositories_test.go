package repositories

import (
	"reflect"
	"testing"
)

func TestAddRepository(t *testing.T) {
	type args struct {
		repo Repository
	}
	tests := []struct {
		name    string
		args    args
		want    RegistryModification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddRepository(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddRepository() got = %v, want %v", got, tt.want)
			}
		})
	}
}
