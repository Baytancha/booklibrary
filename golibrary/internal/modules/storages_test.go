package modules

import (
	"testing"
)

func TestNewStorages(t *testing.T) {
	storages := NewStorages(nil, nil)
	if storages == nil {
		t.Fatal("storages is nil")
	}
}
