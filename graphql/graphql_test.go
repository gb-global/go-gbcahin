package graphql

import (
	"testing"
)

func TestBuildSchema(t *testing.T) {
	// Make sure the schema can be parsed and matched up to the object model.
	if _, err := newHandler(nil); err != nil {
		t.Errorf("Could not construct GraphQL handler: %v", err)
	}
}
