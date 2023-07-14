package uuidgen_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
)

func TestMust(t *testing.T) {
	result := uuidgen.Must("test1", "test2", "test3")
	_, err := uuid.FromString(result)
	if err != nil {
		t.Errorf("Must did not return a valid UUID: %v", err)
	}

	result1 := uuidgen.Must("test1")
	result2 := uuidgen.Must("test2")

	if result1 == result2 {
		t.Errorf("UUIDs for different names should not be the same: got %v", result1)
	}

	result3 := uuidgen.Must("test1", "test2")
	result4 := uuidgen.Must("test1", "test2")

	if result3 != result4 {
		t.Errorf("UUIDs for the same set of names should be the same: got %v and %v", result3, result4)
	}
}

func TestMustSingle(t *testing.T) {
	result := uuidgen.MustSingle("test")
	_, err := uuid.FromString(result)
	if err != nil {
		t.Errorf("MustSingle did not return a valid UUID: %v", err)
	}

	result1 := uuidgen.MustSingle("test1")
	result2 := uuidgen.MustSingle("test1")

	if result1 != result2 {
		t.Errorf("UUIDs for the same name should be the same: got %v and %v", result1, result2)
	}
}
