package validator

import (
	"testing"
)

func TestValidator(t *testing.T) {
	v := New()

	v.AddError("a", "b")

	if v.Valid() {
		t.Error("Expected valid validator")
	}

}
