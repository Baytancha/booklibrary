package helpers

import (
	"net/url"
	"reflect"
	"test/internal/infrastructure/validator"
	"testing"
)

func TestReadString(t *testing.T) {
	type args struct {
		qs           url.Values
		key          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadString(tt.args.qs, tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	type args struct {
		qs           url.Values
		key          string
		defaultValue []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadCSV(tt.args.qs, tt.args.key, tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadInt(t *testing.T) {
	type args struct {
		qs           url.Values
		key          string
		defaultValue int
		v            *validator.Validator
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadInt(tt.args.qs, tt.args.key, tt.args.defaultValue, tt.args.v); got != tt.want {
				t.Errorf("ReadInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
